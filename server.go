package balancer

import (
	"bufio"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

// Server is an implementation of Node
type Server struct {
	host       string
	port       int
	worker     Worker
	configPath string
}

// Register allows the node to register itself to the load balancer.
func (node *Server) Register() error {
	// Open up the file
	file, err := os.OpenFile(node.configPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if it is already registered.
	addr := node.GetHost() + ":" + strconv.Itoa(node.GetPort())
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == addr {
			/* ----- TO DO ----- */
			// Create an error to indicate already registered.
			return nil
		}
	}

	// Register itself.
	writer := bufio.NewWriter(file)
	writer.WriteString(addr + "\n")

	return nil
}

// GetHost returns the hostname that the node is being ran on.
func (node *Server) GetHost() string {
	return node.host
}

// GetPort returns the port that the node is being ran on.
func (node *Server) GetPort() int {
	return node.port
}

// ListenAndServe is how the balancer assigns work.
func (node *Server) ListenAndServe() error {
	// Handle using RPC.
	rpc.RegisterName("Status", node.Status)
	rpc.RegisterName("Do", node.Do)
	rpc.HandleHTTP()

	// Setup the listener.
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(node.GetPort()))
	if err != nil {
		return err
	}

	// Serve.
	return http.Serve(listener, nil)
}

// Status is the implementation to get the status of the node
//   in an RPC-accessible way.
func (node *Server) Status(requestTime time.Time, response *string) error {
	status := node.worker.Status(requestTime).String()
	response = &status
	return nil
}

// Do (work) is the implementation to instruct the node to do work
//   in an RPC-accessible way.
func (node *Server) Do(work string, response *string) error {
	finishedWork, err := node.worker.Do(work)
	response = &finishedWork
	return err
}

// NewNode returns an implementation of Node.
func NewNode(host string, port int, worker Worker, configPath string) Node {
	return &Server{host, port, worker, configPath}
}
