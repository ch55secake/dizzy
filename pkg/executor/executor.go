package executor

// Executor interface, will kick of dispatchers and complete queued jobs
type Executor interface {
	Execute()
}
