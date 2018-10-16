package shell

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type JobProcessor func() (*RestResponse, error)
type JobCompletion func(job int, bm *Benchmark, resp *RestResponse)
type JobMaker func() JobProcessor

// Process jobs concurrently
// Note: RestClient will open new connections for each worker because
// most likely all connections will be in use until a worker completes
// a request and frees a connection.
//
// Warming jobs was introduced to help get all workers started with a
// extra iterations upfront that effectively start all workers.
// (TBD) Are there better mechanisms as warming jobs can still be
// oustanding when jobs are started (affects BM wall time, but iteration
// time should not be affected as long as MaxIdleConsPerHost isn't
// exceeded or some other anomallys of Go or the OS)
//
func ProcessJob(makeJob JobMaker, completion JobCompletion, cancel *bool) *Benchmark {
	iterations := GetCmdIterationValue()
	concurrency := GetCmdConcurrencyValue()
	warming := 0
	if IsCmdWarmingEnabled() {
		warming = concurrency
	}

	var waitGroup sync.WaitGroup

	bm := NewBenchmark(iterations)

	jobs := make(chan int, iterations)
	closeJobs := true
	defer func() {
		if closeJobs {
			close(jobs)
			closeJobs = false
		}
	}()

	// Setup workers for consuming jobs
	for t := 0; t < concurrency; t++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			for job := range jobs {
				if cancel != nil && *cancel {
					break
				}

				if job >= 0 {
					processor := makeJobWithThrottle(makeJob, GetCmdIterationThrottleMs())

					bm.StartIteration(job)
					resp, err := processor()
					bm.EndIteration(job)
					if err != nil {
						bm.SetIterationStatus(job, err)
					} else {
						if completion != nil {
							// Handle command specific completion
							completion(job, &bm, resp)
						} else if resp.GetStatus() != http.StatusOK {
							// Handle default completion
							// TODO: This should not be an error
							msg := resp.GetStatusString()
							bm.SetIterationStatus(job, errors.New(msg))
						}
					}
				} else {
					// Performing warming; do not care ab out throttling or results
					processor := makeJob()
					_, _ = processor()
				}
			}
		}()
	}

	// Run warming jobs
	if warming > 0 {
		for i := 0; i < warming; i++ {
			jobs <- -1
		}

		// Give a little time for connections; longer the better but we don't want to wait too long
		Delay(350 * time.Millisecond)
	}

	// Run the jobs
	bm.Start()
	for i := range bm.Iterations {
		if *cancel {
			break
		}
		jobs <- i
	}

	// Close job channel so workers exit when empty
	if closeJobs {
		close(jobs)
		closeJobs = false
	}

	// Wait for all job threads to exit
	waitGroup.Wait()
	bm.End()
	return &bm
}

// Function overlaps the creation of a job with a delay as job
// creation could potentially make network calls
func makeJobWithThrottle(makejob JobMaker, throttleMs int) JobProcessor {
	var processor JobProcessor

	if throttleMs > 0 {
		// When throttling overlap make job and delay
		messages := make(chan bool)
		go func() {
			Delay(time.Duration(throttleMs) * time.Millisecond)
			messages <- true
		}()
		processor = makejob()
		_ = <-messages
	} else {
		processor = makejob()
	}
	return processor
}
