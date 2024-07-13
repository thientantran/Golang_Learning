package asyncjob

import (
	"context"
	"log"
	"time"
)

//Job requirements:
// 1. Job can do something(handler)
// 2. Job can retry
//   2.1 Config retry times and duration
// 3. Should be stateful
// 4. We should have job manager to manage all jobs

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDurations(times []time.Duration)
}

const (
	defautMaxTimeout     = 10 * time.Second
	defaultMaxRetryCount = 3
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 5, time.Second * 10}
)

type JobHandler func(ctx context.Context) error

func (js JobState) String() string {
	return []string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type jobConfig struct {
	maxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler) *job {
	j := job{
		config: jobConfig{
			maxTimeout: defautMaxTimeout,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		retryIndex: -1,
		state:      StateInit,
		stopChan:   make(chan bool),
	}

	return &j
}

func (j *job) Execute(ctx context.Context) error {
	j.state = StateRunning

	var err error

	err = j.handler(ctx)

	if err != nil {
		j.state = StateFailed
		return err
	}

	j.state = StateCompleted
	return nil

	//ch := make(chan error)
	//ctxJob, doneFunc := context.WithCancel(ctx)
	//
	//go func() {
	//	j.state = StateRunning
	//	var err error
	//
	//	err = j.handler(ctxJob)
	//
	//	if err != nil {
	//		j.state = StateFailed
	//		ch <- err
	//		return
	//	}
	//
	//	j.state = StateCompleted
	//	ch <- err
	//}()
	//
	//for {
	//	select {
	//	case <-j.stopChan:
	//		break
	//	default:
	//		fmt.Println("Job is running")
	//
	//	}
	//}
	//
	////go func() {
	////	for {
	////	}
	////}()
	//select {
	//case err := <-ch:
	//	doneFunc()
	//	return err
	//case <-j.stopChan:
	//	doneFunc()
	//	return nil
	//}
	//
	//return <-ch
}

func (j *job) Retry(ctx context.Context) error {
	//if j.retryIndex == len(j.config.Retries) - 1 {
	//	return nil
	//}
	j.retryIndex += 1

	time.Sleep(j.config.Retries[j.retryIndex])
	log.Println("Retry index:", j.retryIndex)
	//j.state = StateRunning
	err := j.Execute(ctx)

	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == len(j.config.Retries)-1 {
		j.state = StateRetryFailed
		return err
	}

	j.state = StateFailed

	return err
}

func (j *job) State() JobState {
	return j.state
}

func (j *job) RetryIndex() int {
	return j.retryIndex
}

func (j *job) SetRetryDurations(times []time.Duration) {
	if len(times) == 0 {
		return
	}

	j.config.Retries = times
}
