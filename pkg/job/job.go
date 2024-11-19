package job

import (
	"github.com/ch55secake/dizzy/pkg/client"
	"log"
)

// Task represents the function type for job logic and also what will be done
type Task func(client *client.Requester)

// Job represents a unit of work with custom logic
type Job struct {
	ID      int
	Execute Task
}

// NewJob will return a job with a random id and a given request, this job will then be added to the queue
func NewJob(id int, request client.Request) *Job {
	return &Job{
		ID: id,
		Execute: func(client *client.Requester) {
			err, _ := client.MakeRequest(request)
			if err != nil {
				log.Fatalf("Error creating new job with id: %d and err: %v", id, err)
			}
		},
	}
}
