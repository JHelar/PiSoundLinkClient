package soundlink

import (
	"log"
	"net"
)

type NodeMessage struct {
	nodeID  string
	message string
}

type NodeBag struct {
	nodes    map[string]*Node
	Joins    chan net.Conn
	incoming chan NodeMessage
	outgoing chan string
}

func (nb *NodeBag) send(data string, nodeID string) {
	if node, ok := nb.nodes[nodeID]; ok {
		node.outgoing <- data
	}
}

func (nb *NodeBag) join(connection net.Conn) {
	node := NewClient(connection)
	nb.nodes[node.ID] = node
	go func() {
		for {
			nb.incoming <- NodeMessage{
				nodeID:  node.ID,
				message: <-node.incoming,
			}
		}
	}()
}

func (nb *NodeBag) listen() {
	go func() {
		for {
			select {
			case data := <-nb.incoming:
				log.Printf("GOT: %s, FROM: %d", data.message, data.nodeID)
			case conn := <-nb.Joins:
				log.Printf("New node connecting...")
				nb.join(conn)
			}
		}
	}()
}

func NewNodeBag() *NodeBag {
	nodeBag := &NodeBag{
		nodes:    make(map[string]*Node, 0),
		Joins:    make(chan net.Conn),
		incoming: make(chan NodeMessage),
		outgoing: make(chan string),
	}
	nodeBag.listen()

	return nodeBag
}
