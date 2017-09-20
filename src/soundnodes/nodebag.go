package soundnodes

import (
	"fmt"
	"log"
	"net"
)

type NodeMessage struct {
	nodeID  string
	Message string
}

type NodeBag struct {
	nodes          map[string]*Node
	Joins          chan net.Conn
	incoming       chan NodeMessage
	outgoing       chan string
	MessageHandler func(message NodeMessage) (NodeMessage, error)
}

// Send sends message to given client ID
func (nb *NodeBag) Send(message NodeMessage) {
	if node, ok := nb.nodes[message.nodeID]; ok {
		node.outgoing <- message.Message
	}
}

// SendRaw sends message to given client ID
func (nb *NodeBag) SendRaw(nodeID, message string) {
	if node, ok := nb.nodes[nodeID]; ok {
		node.outgoing <- message
	}
}

// Broadcast sends given message to all clients
func (nb *NodeBag) Broadcast(message NodeMessage) {
	for _, val := range nb.nodes {
		val.outgoing <- message.Message
	}
}

func (nb *NodeBag) join(connection net.Conn) {
	node := NewClient(connection)
	nb.nodes[node.ID] = node
	go func() {
		for {
			nb.incoming <- NodeMessage{
				nodeID:  node.ID,
				Message: <-node.incoming,
			}
		}
	}()
}

func (nb *NodeBag) listen() {
	go func() {
		for {
			select {
			case data := <-nb.incoming:
				go func() {
					if response, err := nb.MessageHandler(data); err == nil {
						nb.Send(response)
					} else {
						log.Print(err)
					}
				}()
			case conn := <-nb.Joins:
				log.Printf("New node connecting...")
				nb.join(conn)
			}
		}
	}()
}

// NewNodeBag: Create a new NodeBag
func NewNodeBag() *NodeBag {
	nodeBag := &NodeBag{
		nodes:    make(map[string]*Node, 0),
		Joins:    make(chan net.Conn),
		incoming: make(chan NodeMessage),
		outgoing: make(chan string),
		MessageHandler: func(message NodeMessage) (NodeMessage, error) {
			message.Message = fmt.Sprintf("RESPONSE %s", message.Message)
			return message, nil
		},
	}
	nodeBag.listen()

	return nodeBag
}
