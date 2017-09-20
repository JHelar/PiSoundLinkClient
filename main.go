package main

import (
	"log"
	"net"
	"soundlink"
	"spotify"
)

func main() {
	SLM := soundlink.New()
	err := spotify.Register(SLM)
	if err != nil {
		log.Panic(err)
	}

	nodeListener, _ := net.Listen("tcp", ":6666")
	log.Print("NodeListener Active")
	go func() {
		for {
			conn, _ := nodeListener.Accept()
			SLM.NodesNodeBag.Joins <- conn
		}
	}()

	clientListener, _ := net.Listen("tcp", ":8080")
	log.Print("ClientListener Active")
	for {
		conn, _ := clientListener.Accept()
		SLM.ClientNodeBag.Joins <- conn
	}
}
