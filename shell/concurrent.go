/////////////////////////////////////////////////////////
//  Package for running jobs that make REST requests
//
// Note: RestClient will open new connections for each worker because
// most likely all connections will be in use until a worker completes
// a request and frees a connection.
//
// Warming jobs was introduced to help get all workers started with extra
// iterations upfront to effectively handle handshakes for first connection.
// (TBD) Are there better mechanisms as warming jobs can still be
// oustanding when jobs are started (affects benchmark wall time, but iteration
// time should not be affected as long as MaxIdleConsPerHost isn't
// exceeded or some other anomallys of Go or the OS)
//

package shell

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// JobMonitor -- interface to support benchmark capabilities
type JobMonitor interface {
	Start()
	End()
	StartIteration(int)
	EndIteration(int)
	SetIterationStatus(int, error)
}

// JobFunction -- Function prototype for a function that will perform an instance of the job
type JobFunction func() (*RestResponse, error)

// JobCompletion -- Function prototype for a function that can parse the response
type JobCompletion func(job int, jm JobMonitor, resp *RestResponse)

// JobMaker -- Function prototype for a function that can create an instance of the job function
type JobMaker func() JobFunction

// JobOptions -- options available to the job processing engine
type JobOptions struct {
	Concurrency       int
	Iterations        int
	Duration          time.Duration
	ThrottleInMs      int
	EnableWarming     bool
	JobMaker          JobMaker
	CompletionHandler JobCompletion
	CancelPtr         *bool
}

// GetJobOptionsFromParams -- initializes options from common command line options
func GetJobOptionsFromParams() JobOptions {
	return JobOptions{
		Concurrency:   GetCmdConcurrencyValue(),
		Iterations:    GetCmdIterationValue(),
		Duration:      time.Duration(0),
		ThrottleInMs:  GetCmdIterationThrottleMs(),
		EnableWarming: IsCmdWarmingEnabled(),
	}
}

// ProcessJob -- Run the jobs based on provided options
func ProcessJob(options JobOptions, jm JobMonitor) {
	if options.Iterations <= 0 && options.Duration <= 0 {
		return
	}

	if options.JobMaker == nil {
		return
	}

	concurrency := options.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	warming := 0
	if options.EnableWarming {
		warming = concurrency
	}

	// Wait group to single the completion of a job thread
	var waitGroup sync.WaitGroup

	// Create a channel to control jobs
	jobs := make(chan int, concurrency*2)
	closeJobs := true
	defer func() {
		if closeJobs {
			close(jobs)
			closeJobs = false
		}
	}()

	// Setup workers for consuming jobs
	endTime := time.Now().Add(options.Duration)
	for t := 0; t < concurrency; t++ {
		waitGroup.Add(1) // Add a waiter
		go func() {
			defer waitGroup.Done() // Subtract a waiter

			for job := range jobs {
				if options.Duration > 0 && time.Now().After(endTime) {
					continue // Keep pulling jobs in case producer is blocked
				}

				if options.CancelPtr != nil && *options.CancelPtr {
					continue
				}

				if job >= 0 {
					processor, err := makeJobWithThrottle(options.JobMaker, options.ThrottleInMs)
					if err != nil {
						jm.StartIteration(job)
						jm.EndIteration(job)
						jm.SetIterationStatus(job, err)
						continue
					}
					jm.StartIteration(job)
					resp, err := callJobWithPanicHandler(processor)
					jm.EndIteration(job)
					if err != nil {
						jm.SetIterationStatus(job, err)
					} else {
						if options.CompletionHandler != nil {
							// Handle command specific completion
							options.CompletionHandler(job, jm, resp)
						} else if resp.GetStatus() != http.StatusOK {
							// Handle default completion; expects 200 status
							// Use a custom completion handler for different statuses
							msg := resp.GetStatusString()
							jm.SetIterationStatus(job, errors.New(msg))
						}
					}
				} else {
					// Performing warming; do not care about throttling or results
					processor, err := makeJobWithThrottle(options.JobMaker, 0)
					if err == nil {
						_, _ = callJobWithPanicHandler(processor)
					}
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

	// Run the jobs ; stop after all iterations or end of time whichever comes first
	jm.Start()
	for i := 0; (options.Iterations == 0 || i < options.Iterations) && (options.Duration == 0 || time.Now().Before(endTime)); i++ {
		if options.CancelPtr != nil && *options.CancelPtr {
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
	jm.End()
}

// MakeJobCompletionForExpectedStatus -- Create a completion handler
// to accept a different status than StatusOK
func MakeJobCompletionForExpectedStatus(status int) JobCompletion {
	return func(job int, jm JobMonitor, resp *RestResponse) {
		if resp.GetStatus() != status {
			msg := resp.GetStatusString()
			jm.SetIterationStatus(job, errors.New(msg))
		}
	}
}

// Function overlaps the creation of a job with a delay as job
// creation could potentially make network calls
func makeJobWithThrottle(makejob JobMaker, throttleMs int) (processor JobFunction, err error) {
	defer func() {
		// Absorb panics from make job
		if r := recover(); r != nil {
			err = fmt.Errorf("Make Job Panic'ed: %v", r)
			processor = nil
		}
	}()

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
	return processor, nil
}

func callJobWithPanicHandler(processor JobFunction) (resp *RestResponse, err error) {
	defer func() {
		// Absorb panics from make job
		if r := recover(); r != nil {
			err = fmt.Errorf("Job Panic'ed: %v", r)
			resp = nil
		}
	}()

	return processor()
}
