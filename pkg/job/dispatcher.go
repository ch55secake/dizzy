package job

import (
	"github.com/ch55secake/dizzy/pkg/client"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Dispatcher struct {
	WorkerPool chan chan *Job
	JobQueue   chan *Job
	Workers    []*Worker
	wg         *sync.WaitGroup
	batchSize  int
}

// NewDispatcher creates a new dispatcher for a given job queue, with provided number of workers and provided queue size
func NewDispatcher(numWorkers, queueSize int) *Dispatcher {
	workerPool := make(chan chan *Job, numWorkers)
	jobQueue := make(chan *Job, queueSize)
	workers := make([]*Worker, numWorkers)

	return &Dispatcher{
		WorkerPool: workerPool,
		JobQueue:   jobQueue,
		Workers:    workers,
		wg:         &sync.WaitGroup{},
		batchSize:  300, // hardcode limit of 300 batch so that it doesn't fail overload the execution
	}
}

// Run starts the dispatcher and workers, will also dispatch with the given requester
func (d *Dispatcher) Run(r *client.Requester) {
	for i := range d.Workers {
		worker := &Worker{
			ID:         i,
			JobChannel: make(chan *Job),
			Requester:  r,
			wg:         d.wg,
		}
		worker.Start()
		d.WorkerPool <- worker.JobChannel
		d.Workers[i] = worker
	}

	go d.dispatch()
}

// dispatch assigns jobs to workers in batches
func (d *Dispatcher) dispatch() {
	var batch []*Job
	timer := time.NewTicker(100 * time.Millisecond) // Optional batching timeout
	defer timer.Stop()

	for {
		select {
		case job, ok := <-d.JobQueue:
			if !ok {
				d.dispatchBatch(batch)
				return
			}

			batch = append(batch, job)

			if len(batch) >= d.batchSize {
				d.dispatchBatch(batch)
				batch = []*Job{}
			}
		case <-timer.C:
			if len(batch) > 0 {
				d.dispatchBatch(batch)
				batch = []*Job{}
			}
		}
	}
}

// dispatchBatch sends a batch of jobs to workers
func (d *Dispatcher) dispatchBatch(batch []*Job) {
	for _, job := range batch {
		workerChannel := <-d.WorkerPool
		workerChannel <- job
		d.WorkerPool <- workerChannel
	}
}

// Submit adds a job to the job queue
func (d *Dispatcher) Submit(job *Job) {
	log.Debugf("Submitting job: %v", job.ID)
	d.wg.Add(1)
	d.JobQueue <- job
}

// Wait blocks until all jobs are processed
func (d *Dispatcher) Wait() {
	d.wg.Wait()
	close(d.JobQueue)
}
