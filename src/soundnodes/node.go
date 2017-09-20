package soundnodes

import (
	"bufio"
	"log"
	"net"

	"github.com/nu7hatch/gouuid"
)

type Node struct {
	ID       string
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (node *Node) read() {
	for {
		line, _ := node.reader.ReadString('\n')
		node.incoming <- line
	}
}

func (node *Node) write() {
	for data := range node.outgoing {
		node.writer.WriteString(data)
		node.writer.Flush()
	}
}

func (node *Node) listen() {
	go node.read()
	go node.write()
}

func NewClient(connection net.Conn) *Node {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	node := &Node{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	u, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
	}

	node.ID = u.String()

	node.listen()

	log.Printf("New client: %+v", node)

	return node
}
