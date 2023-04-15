package relayer

import (
	"context"
	"fmt"
	"log"
	"time"

	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

type Worker struct {
	ID               int
	MaxJobsRetries   int
	MaxTaskRetry     int
	waitOnTaskFailed time.Duration
	waitOnJobFailed  time.Duration
	inputChan        chan *rtypes.Job
	quitChan         chan bool
}

func NewWorker() *Worker {
	return &Worker{
		inputChan:        make(chan *rtypes.Job),
		quitChan:         make(chan bool),
		MaxJobsRetries:   3,
		MaxTaskRetry:     3,
		waitOnTaskFailed: 10 * time.Second,
		waitOnJobFailed:  10 * time.Second,
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("Worker %d started\n", w.ID)
	for {
		select {
		case job := <-w.inputChan:
			log.Printf("Worker %d received job %s \n", w.ID, job.ID)
			if err := w.ProcessJob(job); err != nil {
				log.Printf("Worker %d has reached the maximum number of retries for job %s. Stopping worker\n", w.ID, job.ID)
				return

			}
			if !job.Tasks.Result().Finished {
				log.Printf("Worker %d has not been able to complete the job %s correctly . Stopping worker\n", w.ID, job.ID)
				return
			}
			log.Printf("Worker %d completed job %s successfully\n", w.ID, job.ID)

		case <-ctx.Done():
			fmt.Printf("Worker %d stopped\n", w.ID)
			return
		}
	}
}

func (w *Worker) ProcessJob(job *rtypes.Job) error {
	var err error
	for i := 0; job.Retry < w.MaxJobsRetries; i++ {
		err = w.doJob(job)
		if err == nil {
			return nil
		}
		job.Retry++
		wait := w.waitOnJobFailed * time.Duration(job.Retry)

		//log.Printf("Worker %d will retry job %d after %s\n", w.ID, job.ID, wait)
		log.Printf("Worker %d encountered an error processing job %s (attempt %d of %d): %s\n", w.ID, job.ID, job.Retry, w.MaxJobsRetries, err.Error())

		time.Sleep(wait)

	}
	return err
}

func (w *Worker) doJob(job *rtypes.Job) error {
	var currentTask *rtypes.Task
	for {

		nextTask := job.Tasks.GetNextTask(currentTask)
		if nextTask == nil {
			break
		}
		err := w.processTask(nextTask)
		if err != nil {
			return err
		}
		currentTask = nextTask
	}
	return nil
}

func (w *Worker) processTask(t *rtypes.Task) error {
	for i := 0; t.Retry < w.MaxTaskRetry; i++ {
		log.Println("Processing task", t.ID, "attempt", t.Retry)
		result := t.Do()
		if result.Err.Err == nil || result.Err.Skip {
			return nil
		}

		t.Retry++
		wait := w.waitOnJobFailed * time.Duration(t.Retry)
		time.Sleep(wait)

	}
	return t.TResult.Err.Err

}
