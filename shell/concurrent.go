/////////////////////////////////////////////////////////////////////////
// Package for running jobs that make REST requests
//
// Note: With concurrecy the RestClient will open new TCP/IP connections for
// each worker in a concurrecy set. This leads to many requests having
// inflated time under concurrency.
//
// Warming jobs are used to ensure a TCP/IP connection is opened by each
// thread in a concurrency set. The warming job makes a connection that
// is not counted within the requested iterations or time.started with extra.
//
// (TBD) Are there better mechanisms as warming jobs can still be
// oustanding when jobs are started (affects benchmark wall time, but iteration
// time should not be affected as long as MaxIdleConsPerHost isn't
// exceeded or some other anomallys of Go or the OS)
//

package shell

import (
	"errors"
	"sync"
	"time"
)

// Mockable interfaces for concurrent testing and benchmarks
var mockableTimeSince = time.Since
var mockableTimeNow = time.Now

// JobMonitor -- interface to support benchmark capabilities
type JobMonitor interface {
	Start()
	StartIteration(int) JobContext
	FinalizeIteration(JobContext)
	End()
}

// JobContext is returned by the start of an iteration to collect
// iteration information for the completed iteration
type JobContext interface {
	EndIteration(error)
	UpdateError(error)
}

// JobFunction -- Function prototype for a function that will perform an instance of the job
type JobFunction func() (*RestResponse, error)

// JobCompletion -- Function prototype for a function that can parse the response
type JobCompletion func(job int, jc JobContext, resp *RestResponse)

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
	Logger            Logger
}

type JobProcessor struct {
	maker             JobMaker
	monitor           JobMonitor
	throttle          int
	logger            Logger
	completionHandler JobCompletion
}

func NewJobProcessor(logger Logger, maker JobMaker, monitor JobMonitor, completion JobCompletion, throttleMs int) JobProcessor {
	if logger == nil {
		logger = NewLogger(false, false)
	}
	return JobProcessor{
		maker:             maker,
		monitor:           monitor,
		throttle:          throttleMs,
		logger:            logger,
		completionHandler: completion,
	}
}

func (jp JobProcessor) RunProcessor(iterations int, concurrency int, duration time.Duration, cancelPtr *bool) {
	var waitGroup sync.WaitGroup
	var endTime time.Time
	logger := jp.logger
	jobs := make(chan int, concurrency*2)
	closeJobs := true
	defer func() {
		if closeJobs {
			close(jobs)
			closeJobs = false
		}
	}()

	// Setup workers for consuming jobs
	for t := 0; t < concurrency; t++ {
		waitGroup.Add(1) // Add a waiter
		go func() {
			defer waitGroup.Done() // Subtract a waiter
			worker := NewWorker(jp.logger, jp.maker, jp.monitor, jp.completionHandler, jp.throttle, cancelPtr)
			worker.Process(jobs, &endTime, cancelPtr)
		}()
	}

	// Run the jobs ; stop after all iterations or end of time whichever comes first
	logger.LogDebug("Initializing monitor")
	jp.startMonitor()
	if duration == 0 {
		endTime = time.Now().Add(time.Hour * 72)
	} else {
		endTime = time.Now().Add(duration)
	}
	logger.LogDebug("Starting jobs")
	for i := 0; (iterations == 0 || i < iterations) && (duration == 0 || time.Now().Before(endTime)); i++ {
		if cancelPtr != nil && *cancelPtr {
			break
		}
		jobs <- i
	}
	logger.LogDebug("Finished job loop")

	// Close job channel so workers exit when empty
	if closeJobs {
		close(jobs)
		closeJobs = false
	}

	// Wait for all job threads to exit
	waitGroup.Wait()
	jp.endMonitor()
	logger.LogDebug("Finished waiting for jobs to complete")
}

func (jp JobProcessor) startMonitor() {
	if jp.monitor != nil {
		jp.monitor.Start()
	}
}

func (jp JobProcessor) endMonitor() {
	if jp.monitor != nil {
		jp.monitor.End()
	}
}

// GetJobOptionsFromParams -- initializes options from common command line options
func GetJobOptionsFromParams() JobOptions {
	return JobOptions{
		Concurrency:   GetCmdConcurrencyValue(),
		Iterations:    GetCmdIterationValue(),
		Duration:      GetCmdDurationValueWithFallback(0),
		ThrottleInMs:  GetCmdIterationThrottleMs(),
		EnableWarming: IsCmdWarmingEnabled(),
		Logger:        NewLogger(IsCmdVerboseEnabled(), IsCmdDebugEnabled()),
	}
}

// ProcessJob -- Run the jobs based on provided options
func ProcessJob(options JobOptions, jm JobMonitor) {
	logger := options.Logger
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

	if options.EnableWarming {
		warmer := NewJobProcessor(logger, options.JobMaker, nil, nil, 0)
		warmer.RunProcessor(concurrency, concurrency, 0, nil)
	}

	processor := NewJobProcessor(logger, options.JobMaker, jm, options.CompletionHandler, options.ThrottleInMs)
	processor.RunProcessor(options.Iterations, concurrency, options.Duration, options.CancelPtr)
}

// MakeJobCompletionForExpectedStatus -- Create a completion handler
// to accept a different status than StatusOK
func MakeJobCompletionForExpectedStatus(status int) JobCompletion {
	return func(job int, jc JobContext, resp *RestResponse) {
		if resp.GetStatus() != status {
			msg := resp.GetStatusString()
			jc.UpdateError(errors.New(msg))
		}
	}
}
