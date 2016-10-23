package balancer

import (
	"errors"
	"time"
)

// Server s
type Server struct {
	balancer     Balancer
	worker       Worker
	finishedJobs []Work
	jobs         []Work
}

func (server *Server) getBalancer() Balancer {
	return server.balancer
}

func (server *Server) setBalancer(balancer Balancer) {
	server.balancer = balancer
}

// Register registers the node to a load balancer
func (server *Server) Register() error {
	return server.balancer.AddNode(server)
}

// QueueJob adds another job to the server's workload
func (server *Server) QueueJob(job Job) error {
	if work, ok := job.(Work); ok {
		server.jobs = append(server.jobs, work)
		return nil
	}
	return errors.New("Invalid job")
}

// Status returns the status of the node to the load balancer
func (server *Server) Status(requestTime time.Time) Status {
	return Status{
		requestTime:  requestTime,
		responseTime: time.Now(),
		workAmount:   len(server.jobs),
	}
}

// DoWork does an arbitrary task
func (server *Server) DoWork() {

}

// NewServer s
func NewServer(balancer Balancer, worker Worker) Node {
	server := Server{
		balancer,
		worker,
		make([]Work, 0),
		make([]Work, 0),
	}

	return &server
}
