package balancer

import (
	"fmt"
	"strconv"
)

// Master is a implementation of Balancer.
type Master struct {
	host  string
	port  int
	nodes []NodeConnection
}

// ListenAndServe is used to accept new work for the nodes.
func (balancer *Master) ListenAndServe() error {
	return nil
}

// MaintainNodes should be run alongside ListenAndServe to maintain nodes.
func (balancer *Master) MaintainNodes() error {
	// Do initial start up, and then enter a loop.

	// Is this going to be file notified or socket notified?
	//   Decide later and do that as startup.

	// Initial start up

	//

	// Endless loop.
	for {
		// Determine if any nodes should be taken down.
		for _, node := range balancer.nodes {
			status, err := node.GetStatus()
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
				strconv.Itoa(node.GetPort()) + " - " + status.String())
		}

		// Accept new nodes.
		// If this is socket notified, I need to have started this in
		//   a goroutine before now. Determine that later.
		/* ----- TO DO ----- */

		return nil
	}
}

// Delegate passes along work as a string to a node and returns
//   the reply if any.
func (balancer *Master) Delegate(string) (string, error) {
	return "", nil
}

func (balancer *Master) GetHost() string {
	return balancer.host
}

func (balancer *Master) GetPort() int {
	return balancer.port
}
