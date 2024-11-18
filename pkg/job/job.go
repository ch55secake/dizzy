package job

import "github.com/ch55secake/dizzy/pkg/client"

// Task represents the function type for job logic and also what will be done
type Task func(client *client.Requester)

// Job represents a unit of work with custom logic
type Job struct {
	ID      int
	Execute Task
	Request client.Request
}

// NewJob will return a job with a random id and a given request, this job will then be added to the queue
func NewJob(id int, request client.Request) *Job {
	return &Job{}
}
