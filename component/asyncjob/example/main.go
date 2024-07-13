package main

import (
	"Food-delivery/component/asyncjob"
	"context"
	"errors"
	"log"
	"time"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second)
		log.Println("I'm job 1")
		return errors.New("Something went wrong at job 1")
	})

	//set time
	job1.SetRetryDurations([]time.Duration{time.Second})
	if err := job1.Execute(context.Background()); err != nil {
		log.Println(job1.State(), err)

		for {
			err := job1.Retry(context.Background())
			if err != nil {
				log.Println(job1.State(), err)
			}
			if job1.State() == asyncjob.StateRetryFailed || job1.State() == asyncjob.StateCompleted {
				break
			}
		}
	}
}
