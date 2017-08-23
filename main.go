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

	listener, _ := net.Listen("tcp", ":6666")
	for {
		conn, _ := listener.Accept()
		SLM.Nodebag.Joins <- conn
	}
	// query := os.Args[1]
	// // Search track
	// log.Printf("Searching: %s", query)
	// songs, err := SLM.Search(query, soundlink.SearchTrack)
	// if err != nil {
	// 	log.Panic(err)
	// } else {
	// 	for _, song := range songs {
	// 		log.Printf("Song: %+v", song)
	// 	}
	// }
}
