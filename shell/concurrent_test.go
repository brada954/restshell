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

func (jc *JobCounter) Start() {
}

func (jc *JobCounter) End() {
}

func (jc *JobCounter) StartIteration(i int) {
	jc.mutex.Lock()
	jc.Initiations++
	jc.mutex.Unlock()
}

func (jc *JobCounter) EndIteration(i int) {
	jc.mutex.Lock()
	jc.Iterations++
	jc.mutex.Unlock()
}

func (jc *JobCounter) EndIterationWithError(i int, err error) {
	jc.mutex.Lock()
	jc.Iterations++
	jc.mutex.Unlock()
	jc.SetIterationStatus(i, err)
}

func (jc *JobCounter) SetIterationStatus(i int, err error) {
	if err != nil {
		jc.mutex.Lock()
		jc.Errors++
		jc.mutex.Unlock()
	}
}

func (jc *JobCounter) UpdateIterationError(i int, err error) {
	jc.SetIterationStatus(i, err)
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

	jc := &JobCounter{}
	ProcessJob(options, jc)

	if jc.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", jc.Iterations)
	}

	if jc.Errors != 0 {
		t.Errorf("Unexpected error count: 1<>%d", jc.Errors)
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

	jc := &JobCounter{}
	ProcessJob(options, jc)

	if jc.Iterations != 3 {
		t.Errorf("Unexpected iteration count: 3<>%d", jc.Iterations)
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

	jc := &JobCounter{}
	ProcessJob(options, jc)

	if jc.Iterations != 6 {
		t.Errorf("Unexpected iteration count: 6<>%d", jc.Iterations)
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

	jc := &JobCounter{}
	ProcessJob(options, jc)

	if jc.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", jc.Iterations)
	}

	if jc.Errors != 1 {
		t.Errorf("Unexpected error count: 1<>%d", jc.Errors)
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

	jc := &JobCounter{}
	ProcessJob(options, jc)

	if jc.Iterations != 5 {
		t.Errorf("Unexpected iteration count: 5<>%d", jc.Iterations)
	}

	if jc.Errors != 4 {
		t.Errorf("Unexpected error count: 4<>%d", jc.Errors)
	}
}
