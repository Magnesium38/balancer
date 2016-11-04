package balancer

import "time"

// Balancer is the interface that defines what methods a load balancer has.
type Balancer interface {
	// ListenAndServe is the main method to be run for the balancer
	//   to delegate work
	ListenAndServe() error

	// MaintainNodes should be run alongside ListenAndServe to handle
	//   nodes that have gone down and accept new nodes
	MaintainNodes() error

	// Delegate accepts a string as work to do and passes it on to a
	//   a node and then waits for the response before returning.
	Delegate(string) (string, error)

	// GetHost and GetPort is what the balancer is listening on.
	GetHost() string
	GetPort() int
}

// A NodeConnection is how the balancer interfaces with the nodes.
type NodeConnection interface {
	// GetHost and GetPort make up the address that the node is located at.
	GetHost() string
	GetPort() int

	// AddJob, FinishJob and AmountOfJobs gives the balancer a general
	//   idea of how much work the balancer currently has.
	AddJob()
	FinishJob()
	GetLoad() int

	// Connect establishes the connection from the balancer to the node.
	Connect() error

	// GetStatus returns a status struct with information about how the
	//   node is currently running.
	GetStatus() (Status, error)

	// Send allows for a simple way to delegate work to a specific node.
	Send(string) (string, error)
}

// A NodeFactory is how a balancer can create new NodeConnections for itself.
type NodeFactory interface {
	Create(string) NodeConnection
}

// Node is the interface that defines the nodes that a load balancer can use.
type Node interface {
	// Register allows the node to self-register itself with the load balancer.
	Register() error

	// GetHost and GetPort is what the node is listening on.
	GetHost() string
	GetPort() int

	// ListenAndServe should be run by the node to continuously accept
	//   work as strings and then reply with finished work as strings.
	ListenAndServe() error
}

// A WorkPool is how a node can complete work and report on its work.
type WorkPool interface {
	Work(string) (string, error)

	// Status accepts a time which is interpreted as the time
	//   the request was made for the status.
	Status(time.Time) Status
}

// A Status is what a balancer can easily ask about its nodes.
type Status interface {
	String() string
	Parse(string)
	GetIdleTime() time.Time
}
