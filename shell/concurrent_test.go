package shell

import (
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"
)

// JobCounter -- A test monitor for counting iterations during tests
type JobCounter struct {
	Iterations  int
	Initiations int
	Errors      int
	mutex       sync.Mutex
}

type JobIteration struct {
	Error error
}

func (job *JobCounter) Start() {
}

func (job *JobCounter) End() {
}

func (job *JobCounter) StartIteration(i int) JobContext {
	job.mutex.Lock()
	job.Initiations++
	job.mutex.Unlock()
	return &JobIteration{}
}

func (job *JobCounter) FinalizeIteration(jc JobContext) {
	job.mutex.Lock()
	job.Iterations++
	if ji, ok := jc.(*JobIteration); ok {
		if ji.Error != nil {
			job.Errors++
		}
	}
	job.mutex.Unlock()
}

func (jc *JobIteration) EndIteration(err error) {
	jc.Error = err
}

func (jc *JobIteration) UpdateError(err error) {
	jc.Error = err
}

var SuccessResponse = RestResponse{httpResp: &http.Response{Status: "StatusOk", StatusCode: 200}}
var NotFoundResponse = RestResponse{httpResp: &http.Response{Status: "NotFound", StatusCode: 404}}

func TestBasicIteration(t *testing.T) {

	var jm JobMaker = func() JobFunction {
		return func() (*RestResponse, error) { return &SuccessResponse, nil }
	}

	options := JobOptions{
		Iterations: 5,
		JobMaker:   jm,
	}

	jobData := JobCounter{}
	job := &jobData
	ProcessJob(options, job)

	if job.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", job.Iterations)
	}

	if job.Errors != 0 {
		t.Errorf("Unexpected error count: 0<>%d", job.Errors)
	}
}

func TestTimeIteration(t *testing.T) {

	var jm JobMaker = func() JobFunction {
		return func() (*RestResponse, error) {
			//fmt.Println("Going to sleep...")
			time.Sleep(time.Second)
			return &RestResponse{}, nil
		}
	}

	options := JobOptions{
		Iterations: 5,
		Duration:   (time.Second * 3) - (time.Millisecond * 100),
		JobMaker:   jm,
	}

	job := &JobCounter{}
	ProcessJob(options, job)

	if job.Iterations != 3 {
		t.Errorf("Unexpected iteration count: 3<>%d", job.Iterations)
	}
}

func TestConcurrentIterationsWithTime(t *testing.T) {
	var jm JobMaker = func() JobFunction {
		return func() (*RestResponse, error) {
			//fmt.Println("Going to sleep...")
			time.Sleep(time.Second * 3)
			return &RestResponse{}, nil
		}
	}

	options := JobOptions{
		Iterations:  6,
		Duration:    time.Second * 7,
		JobMaker:    jm,
		Concurrency: 3,
	}

	job := &JobCounter{}
	ProcessJob(options, job)

	if job.Iterations != 6 {
		t.Errorf("Unexpected iteration count: 6<>%d", job.Iterations)
	}
}

func TestBasicIterationWithError(t *testing.T) {
	var jobID = 0

	var jm JobMaker = func() JobFunction {
		return func() (*RestResponse, error) {
			jobID++
			if jobID == 3 {
				return &RestResponse{}, errors.New("Test Error")
			} else {
				return &SuccessResponse, nil
			}
		}
	}

	options := JobOptions{
		Iterations: 5,
		JobMaker:   jm,
	}

	job := &JobCounter{}
	ProcessJob(options, job)

	if job.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", job.Iterations)
	}

	if job.Errors != 1 {
		t.Errorf("Unexpected error count: 1<>%d", job.Errors)
	}
}

func TestBasicIterationWithAlternateResponse(t *testing.T) {
	var jobID = 0

	var jm JobMaker = func() JobFunction {
		return func() (*RestResponse, error) {
			jobID++
			if jobID == 3 {
				return &NotFoundResponse, nil
			} else {
				return &SuccessResponse, nil
			}
		}
	}

	options := JobOptions{
		Iterations:        5,
		JobMaker:          jm,
		CompletionHandler: MakeJobCompletionForExpectedStatus(404),
	}

	job := &JobCounter{}
	ProcessJob(options, job)

	if job.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", job.Iterations)
	}

	if job.Errors != 4 {
		t.Errorf("Unexpected error count: 4<>%d", job.Errors)
	}
}
