package executor

import "github.com/ch55secake/dizzy/pkg/job"

type DefaultExecutor struct {
	job.Dispatcher
	WorkerCount int
	QueueSize   int
}

// Execute will use the dispatcher provided on the struct and then kick off all jobs that are ready to be dispatched
// Will also rely on the queue size and worker count provided when using the dispatcher to execute any jobs
func (e *DefaultExecutor) Execute() {

}
