package balancer

import "time"

// Balancer is the interface that defines what methods a load balancer has.
type Balancer interface {
	// ListenAndServe is the main method to be run for the balancer
	// to delegate work
	ListenAndServe() error

	// MaintainNodes should be run alongside ListenAndServe to handle
	// nodes that have gone down and accept new nodes
	MaintainNodes() error

	AddNode(Node) error
	RemoveNode(Node) error

	// Delegate accepts a Job and passes it on to a node
	// and then waits for the response before returning.
	Delegate(Job) (interface{}, error)
}

// A Connection allows the Nodes to self-register
// with the associated load balancer.
type Connection interface {
	// Register accepts a Node and notifies the load balancer in some way
	// that a new server would like to be a part of the cluster.
	Register(Node) error
}

// Node is the interface that defines the nodes that a load balancer can use.
type Node interface {
	// Register allows the node to self-register itself with the load balancer.
	Register() error

	// GetConnection returns the load balancer connection.
	GetConnection() Connection

	GetHost() string
	GetPort() int

	// GetStatus returns a simple overview of the Node's status.
	GetStatus() Status

	//// QueueJob accepts a new job for the server to complete.
	//QueueJob(Job) error

	// DoWork is the main method that actually completes the given work.
	DoWork() error

	// AcceptNewJobs should be ran in a goroutine to accept new jobs
	// so that DoWork() can complete them when finished.
}

/*
// Node s
type Node interface {
	getBalancer() Balancer
	setBalancer(Balancer)
	QueueJob(Job) error
	Register() error
	Status(time.Time) Status
	DoWork()
}
*/

// Job s
type Job interface {
}

// Worker s
type Worker interface {
	Do(interface{})
}

// Status s
type Status struct {
	requestTime  time.Time
	responseTime time.Time
	workAmount   int
}

// Work s
type Work struct {
}
