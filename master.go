package balancer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Master is a implementation of Balancer.
type Master struct {
	host          string
	port          int
	nodes         []NodeConnection
	nodeFilePath  string
	pingFrequency time.Duration
	nodeFactory   NodeFactory
}

// ListenAndServe is used to accept new work for the nodes.
func (balancer *Master) ListenAndServe() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(balancer.GetPort()))
	if err != nil {
		return err
	}

	for {
		// Accept new connections.
		connection, err := listener.Accept()
		if err != nil {
			/* ----- TO DO ----- */
			// Determine what causes the accept to fail and handle here.
			//   I have a feeling that just blindly returning would
			//   be a really poor way of handling it.
		}

		// Handle those new connections.
		go func(conn net.Conn) {
			// Read everything from the connection.
			b, err := ioutil.ReadAll(conn)
			if err != nil {
				/* ----- TO DO ----- */
				// Determine what causes the accept to fail and handle here.
				//   I have a feeling that just blindly returning would
				//   be a really poor way of handling it.
			}

			// Delegate the work and get the response.
			response, err := balancer.Delegate(string(b))
			if err != nil {
				/* ----- TO DO ----- */
				// Determine what causes the accept to fail and handle here.
				//   I have a feeling that just blindly returning would
				//   be a really poor way of handling it.
			}

			// If there is a response, send it back.
			if response != "" {
				_, err := conn.Write([]byte(response))
				if err != nil {
					/* ----- TO DO ----- */
					// Determine what causes the accept to fail and handle here.
					//   I have a feeling that just blindly returning would
					//   be a really poor way of handling it.
				}
			}

			conn.Close()
		}(connection)
	}
}

// MaintainNodes should be run alongside ListenAndServe to maintain nodes.
func (balancer *Master) MaintainNodes() error {
	// Initial start up
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Initial read of file
	// Open the file
	file, err := os.Open(balancer.nodeFilePath)
	if err != nil {
		return err
	}

	// Read each line as a node into nodes.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Transform line into a node connection
		node := balancer.nodeFactory.Create(line)

		// Establish a connection to the node.
		err := node.Connect()
		if err != nil {
			return err
		}

		// The connection was successful, keep track of the node.
		balancer.nodes = append(balancer.nodes, node)
	}

	// Check for scanner error.
	if err := scanner.Err(); err != nil {
		return err
	}

	// Nodes are all loaded

	// Endless loop.
	for {
		// Accept new nodes using fsnotify to watch a file.
		select {
		case event := <-watcher.Events:
			fmt.Println(event)
			/* ----- TO DO ----- */
			continue
		case err := <-watcher.Errors:
			fmt.Println(err)
			/* ----- TO DO ----- */
			continue
		default:
			// When there is no event from the watcher, ping nodes if needed.
			//   If they are unresponsive, attempt to reconnect.
			select {
			case <-time.After(balancer.pingFrequency):
				// Handle pinging here.
				for _, node := range balancer.nodes {
					err := node.UpdateStatus()
					status := node.GetStatus()
					if err != nil {
						/* ----- TO DO ----- */
						// Determine what errors GetStatus can return
						//   and how I want to handle them.
					}

					// Handle status.
					// Consider stopping nodes that have been inactive.
					// Need to consider how to bring nodes back up then.
					// For now, just print it nicely.
					fmt.Println(node.GetHost() + ":" +
						strconv.Itoa(node.GetPort()) +
						" - " + status.String())
				}
			default:
			}
		}
	}
}

// Delegate passes along work as a string to a node and returns
//   the reply if any.
func (balancer *Master) Delegate(work string) (string, error) {
	// Find the lowest load and give that node more work.
	var worker NodeConnection
	lowestLoad := -1

	if len(balancer.nodes) == 0 {
		/* ----- TO DO ----- */
		// If there are no NodeConnections, work cannot be delegated.
		//   This needs to be handled somehow.
		return "", nil
	}

FindLowestLoad:
	for _, node := range balancer.nodes {
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

	return worker.Send(work)
}

// GetHost returns the hostname that the balancer is being ran on.
func (balancer *Master) GetHost() string {
	return balancer.host
}

// GetPort returns the port that the balancer is being ran on.
func (balancer *Master) GetPort() int {
	return balancer.port
}
