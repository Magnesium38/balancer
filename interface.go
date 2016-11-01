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

// A NodeConnection is how the balancer interfaces with the nodes.
type NodeConnection interface {
	GetHost() string
	GetPort() int

	// GetStatus returns a status struct with information about how the
	//   node is currently running.
	GetStatus() (Status, error)

	// Send allows for a simple way to delegate work to a specific node.
	Send(string) (string, error)
}

type Status interface {
	String() string
	GetIdleTime() time.Time
}
