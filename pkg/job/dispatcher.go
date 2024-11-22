package job

import (
	"github.com/ch55secake/dizzy/pkg/client"
	"log"
	"sync"
)

// TODO: Figure out how to batch jobs/worker pool at around 300
type Dispatcher struct {
	WorkerPool chan chan *Job
	JobQueue   chan *Job
	Workers    []*Worker
	wg         *sync.WaitGroup
}

// NewDispatcher create a new dispatcher for a given job queue, with provided number of workers and provided queue size
func NewDispatcher(numWorkers int, queueSize int) *Dispatcher {
	workerPool := make(chan chan *Job, numWorkers)
	jobQueue := make(chan *Job, queueSize)
	workers := make([]*Worker, numWorkers)

	return &Dispatcher{
		WorkerPool: workerPool,
		JobQueue:   jobQueue,
		Workers:    workers,
		wg:         &sync.WaitGroup{},
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
		//log.Printf("Starting worker %d", worker.ID)
		worker.Start()
		d.WorkerPool <- worker.JobChannel
		d.Workers[i] = worker
	}

	go d.dispatch()
}

// dispatch assigns jobs to available workers
func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue {
		workerChannel := <-d.WorkerPool
		workerChannel <- job
		d.WorkerPool <- workerChannel
	}
}

// Submit adds a job to the job queue
func (d *Dispatcher) Submit(job *Job) {
	log.Println("Submitting job", job.ID)
	d.wg.Add(1)
	d.JobQueue <- job
}

// Wait blocks until all jobs are processed
func (d *Dispatcher) Wait() {
	d.wg.Wait()
	close(d.JobQueue)
}
