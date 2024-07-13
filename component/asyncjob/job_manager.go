package asyncjob

import (
	"Food-delivery/common"
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
	g.wg.Add(len(g.jobs))
	// waitgroup , khi có n goroutine, thì truyền n vào waitgroup để biết khi nào n goroutine xong, mỗi goroutine xong thì waitgroup sẽ done, đúng số n done thì sẽ xong
	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				defer common.AppRecover()
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])
			continue
		}

		job := g.jobs[i]

		err := g.runJob(ctx, job)

		// dừng lại nếu có job bị lỗi
		if err != nil {
			return err
		}

		errChan <- err
		g.wg.Done()
	}
	// tại sao chuyển lên đây, để đợi tất cả routine chạy xong, vì khi để ở dưới, mà code đổi return err trước thì sẽ không chạy hết các goroutine
	g.wg.Wait()

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
			return v
		}
	}

	return err
}
