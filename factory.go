package balancer

import (
	"strconv"
	"strings"
)

// A ConnectionFactory creates the connection to the nodes for the balancer.
type ConnectionFactory struct {
	status StatusFactory
}

// Create takes the connection info and creates the connection struct.
func (factory *ConnectionFactory) Create(connInfo string) (NodeConnection, error) {
	parts := strings.Split(connInfo, ":")

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	conn := Connection{}
	conn.host = parts[0]
	conn.port = port
	conn.jobCount = 0
	conn.status = factory.status.Create()
	conn.client = nil

	return &conn, nil
}

// NewConnectionFactory returns an implementation of NodeFactory
func NewConnectionFactory(factory StatusFactory) NodeFactory {
	return &ConnectionFactory{factory}
}
