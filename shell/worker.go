package shell

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type worker struct {
	jobMaker   JobMaker
	monitor    JobMonitor
	throttle   int
	cancelPtr  *bool
	logger     Logger
	completion JobCompletion
}

func NewWorker(logger Logger, jm JobMaker, monitor JobMonitor, completion JobCompletion, throttleMs int, cancel *bool) worker {
	if logger == nil {
		logger = NewLogger(false, false)
	}
	worker := worker{
		jobMaker:   jm,
		monitor:    monitor,
		throttle:   throttleMs,
		cancelPtr:  cancel,
		logger:     logger,
		completion: completion,
	}
	return worker
}

func (w worker) Process(jobs <-chan int, endTime *time.Time, cancelPtr *bool) {
	logger := w.logger
	for job := range jobs {
		logger.LogDebug("Pulling new job: %d", job)
		if time.Now().After(*endTime) {
			// Keep pulling jobs in case producer is blocked
			// Overall, we are done starting new requests
			continue
		}

		if cancelPtr != nil && *cancelPtr {
			continue
		}

		logger.LogDebug("Making job: %d interation", job)
		processor, err := w.makeJob()
		if err != nil {
			context := w.StartIteration(job)
			context.EndIteration(err)
			w.FinalizeIteration(context)
			continue
		}
		logger.LogDebug("Starting job: %d interation", job)
		context := w.StartIteration(job)

		resp, err := w.invoke(processor)
		context.EndIteration(err)
		if err == nil {
			if w.completion != nil {
				w.completion(job, context, resp)
			} else if resp.GetStatus() != http.StatusOK {
				context.UpdateError(errors.New(resp.GetStatusString()))
			}
		}
		w.FinalizeIteration(context)
	}
}

func (w worker) makeJob() (worker JobFunction, err error) {
	defer func() {
		// Absorb panics from make job
		if r := recover(); r != nil {
			err = fmt.Errorf("make worker failed: %v", r)
			worker = nil
		}
	}()

	if w.throttle > 0 {
		// When throttling overlap make job and delay
		messages := make(chan bool)
		go func() {
			Delay(time.Duration(w.throttle) * time.Millisecond)
			messages <- true
		}()
		worker = w.jobMaker()
		<-messages
	} else {
		worker = w.jobMaker()
	}
	return worker, nil
}

func (w worker) invoke(worker JobFunction) (resp *RestResponse, err error) {
	defer func() {
		// Absorb panics from the worker
		if r := recover(); r != nil {
			err = fmt.Errorf("job panic: %v", r)
			resp = nil
		}
	}()
	return worker()
}

func (w worker) StartIteration(job int) JobContext {
	if w.monitor != nil {
		return w.monitor.StartIteration(job)
	}
	return newNullJobContext()
}

func (w worker) FinalizeIteration(context JobContext) {
	if w.monitor != nil {
		w.monitor.FinalizeIteration(context)
	}
}

type nullJobContext struct {
}

func newNullJobContext() nullJobContext {
	return nullJobContext{}
}
func (n nullJobContext) EndIteration(error) {

}

func (n nullJobContext) UpdateError(error) {

}
