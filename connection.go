package balancer

// A ConnectionFactory creates the connection to the nodes for the balancer.
type ConnectionFactory struct {
}

// Create takes the connection info and creates the connection struct.
func (factory *ConnectionFactory) Create(connInfo string) Connection {
	return Connection{}
}

// A Connection is a way for a balancer to communicate with a node.
type Connection struct {
}

// GetHost returns the hostname that the node is listening on.
func (conn *Connection) GetHost() string {
	return ""
}

// GetPort returns the port that the node is listening on.
func (conn *Connection) GetPort() string {
	return ""
}

// AddJob increments the internal counter of jobs by 1.
func (conn *Connection) AddJob() {

}

// FinishJob decrements the internal counter of jobs by 1.
func (conn *Connection) FinishJob() {

}

// GetWorkLoad returns the current estimated work load that
//   the node has. It is made possible through Add/FinishJob
func (conn *Connection) GetWorkLoad() int {
	return 0
}

// Connect initiates the connection between the balancer
//   and the node. In this implementation, this only exists
//   to verify that it can be connected to.
func (conn *Connection) Connect() error {
	return nil
}

// GetStatus returns the current Status of the node, or an error
//   if it cannot retrieve the status.
func (conn *Connection) GetStatus() (Status, error) {
	return nil, nil
}

// Send is how a balancer can send work to the nodes. This
//   implementation is using RPC.
func (conn *Connection) Send(work string) (string, error) {
	return "", nil
}
