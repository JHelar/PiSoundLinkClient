package main

import (
	"log"
	"os"
	"soundlink"
	"spotify"
)

func main() {
	SLM := soundlink.New()
	err := spotify.Register(SLM)
	if err != nil {
		log.Panic(err)
	}
	query := os.Args[1]
	// Search track
	log.Printf("Searching: %s", query)
	songs, err := SLM.Search(query, soundlink.SearchTrack)
	if err != nil {
		log.Panic(err)
	} else {
		for _, song := range songs {
			log.Printf("Song: %+v", song)
		}
	}
}
