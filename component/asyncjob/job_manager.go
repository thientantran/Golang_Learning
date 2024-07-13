package asyncjob

import (
	"context"
	"log"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}
			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}

func (g *group) Run(ctx context.Context) error {
	//g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		//if g.isConcurrent {
		//	go func(aj Job) {
		//		defer common.AppRecover()
		//		errChan <- g.runJob(ctx, aj)
		//		g.wg.Done()
		//	}(g.jobs[i])
		//	continue
		//}

		job := g.jobs[i]

		err := g.runJob(ctx, job)

		// dừng lại nếu có job bị lỗi
		if err != nil {
			return err
		}

		errChan <- err
		//g.wg.Done()
	}

	//close(errChan)

	var err error

	//for v := range errChan {
	//	log.Println(v)
	//	if v != nil {
	//		err = v
	//	}
	//}
	for i := 1; i <= len(g.jobs); i++ {
		v := <-errChan
		if v != nil {
			err = v
		}
	}
	//g.wg.Wait()
	return err
}
