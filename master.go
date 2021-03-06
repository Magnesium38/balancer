package balancer

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Master is a implementation of Balancer.
type Master struct {
	host         string
	port         int
	nodes        []NodeConnection
	nodeRegistry string
	frequency    time.Duration
	factory      NodeFactory
}

// Work s
func (master *Master) Work(argument string, reply *string) error {
	response, err := master.Delegate(argument)
	*reply = response
	return err
}

type MasterWorker struct {
	master *Master
}

func (s *MasterWorker) Work(argument string, reply *string) error {
	return s.master.Work(argument, reply)
}

// ListenAndServe is used to accept new work for the nodes.
func (master *Master) ListenAndServe() error {
	rpc.RegisterName("Master", &MasterWorker{master})
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(master.GetPort()))
	if err != nil {
		return err
	}

	return http.Serve(listener, nil)
}

// updateNodes does a full read of the node repository
func (master *Master) updateNodes() error {
	// Open the file
	file, err := os.Open(master.nodeRegistry)
	if err != nil {
		return err
	}

	// Read each line as a node into nodes.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Transform line into a node connection
		node, err := master.factory.Create(line)
		if err != nil {
			return err
		}

		// Check if this node already has this connection.
		keep := true
		for _, existingNode := range master.nodes {
			if node.GetHost() != existingNode.GetHost() {
				continue
			}
			if node.GetPort() != existingNode.GetPort() {
				continue
			}
			keep = false
			break
		}

		if !keep {
			continue
		}

		// Establish a connection to the node.
		err = node.Connect()
		if err != nil {
			return err
		}

		// The connection was successful, keep track of the node.
		master.nodes = append(master.nodes, node)
	}

	// Check for scanner error.
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// MaintainNodes should be run alongside ListenAndServe to maintain nodes.
func (master *Master) MaintainNodes() error {
	// Initial start up
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Add the file to be watched.
	watcher.Add(master.nodeRegistry)

	// Initial read of file
	// Open the file
	file, err := os.Open(master.nodeRegistry)
	if err != nil {
		return err
	}

	// Read each line as a node into nodes.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Transform line into a node connection
		node, err := master.factory.Create(line)
		if err != nil {
			return err
		}

		// Establish a connection to the node.
		err = node.Connect()
		if err != nil {
			return err
		}

		// The connection was successful, keep track of the node.
		master.nodes = append(master.nodes, node)
	}

	// Check for scanner error.
	if err := scanner.Err(); err != nil {
		return err
	}

	// Create a timer as a notification to check node status.
	timer := time.NewTimer(master.frequency)

	// Endless loop.
	for {
		// Accept new nodes using fsnotify to watch a file.

		select {
		case <-watcher.Events:
			// Typically means a new node should be watched.
			err := master.updateNodes()
			if err != nil {
				fmt.Println(err)
			}
			/* ----- TO DO ----- */
			continue
		case err := <-watcher.Errors:
			fmt.Println(err)
			/* ----- TO DO ----- */
			continue
		case <-timer.C:
			// The timer has finished. Check nodes and then reset it.
			// Handle pinging here.
			for _, node := range master.nodes {
				err := node.UpdateStatus()
				status := node.GetStatus()
				if err != nil {
					/* ----- TO DO ----- */
					// Determine what errors GetStatus can return
					//   and how I want to handle them.
				}

				// Handle status.
				// TO DO:
				// Consider stopping nodes that have been inactive.
				// Need to consider how to bring nodes back up then.
				// For now, just print it nicely.
				fmt.Println(node.GetHost() + ":" +
					strconv.Itoa(node.GetPort()) +
					" - " + status.String())
			}

			timer.Stop()
			timer.Reset(master.frequency)
		}
	}
}

// Delegate passes along work as a string to a node and returns
//   the reply if any.
func (master *Master) Delegate(work string) (string, error) {
	// Find the lowest load and give that node more work.
	var worker NodeConnection
	lowestLoad := -1

	if len(master.nodes) == 0 {
		/* ----- TO DO ----- */
		// If there are no NodeConnections, work cannot be delegated.
		//   This needs to be handled somehow.
		fmt.Println("ERROR No nodes.")
		return "", nil
	}

FindLowestLoad:
	for _, node := range master.nodes {
		currentLoad := node.GetWorkLoad()

		switch {
		case lowestLoad == -1:
			// The first node sets the initial bar for work load.
			lowestLoad = currentLoad
			worker = node
		case currentLoad < lowestLoad:
			lowestLoad = currentLoad
			worker = node
		case currentLoad == 0:
			// If a worker is doing nothing, it's definitely the lowest.
			//   Just set it as the worker and break to send it work.
			worker = node
			break FindLowestLoad
		}
	}

	worker.AddJob()
	defer worker.FinishJob()
	return worker.Send(work)
}

// GetHost returns the hostname that the balancer is being ran on.
func (master *Master) GetHost() string {
	return master.host
}

// GetPort returns the port that the balancer is being ran on.
func (master *Master) GetPort() int {
	return master.port
}

// NewLoadBalancer returns an implementation of Balancer.
func NewLoadBalancer(host string, port int, nodeRegistry string, frequency time.Duration, factory NodeFactory) Balancer {
	master := Master{}

	master.host = host
	master.port = port
	master.nodeRegistry = nodeRegistry
	master.frequency = frequency
	master.factory = factory

	return &master
}
