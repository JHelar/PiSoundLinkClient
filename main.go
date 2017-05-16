package main

import (
	"soundlink"
	"spotify"
	"log"
)

func main(){
	SLM := soundlink.New()
	err := spotify.Register(SLM)
	if err != nil {
		log.Panic(err)
	}
	// Search track
	songs, err := SLM.SearchSpecificSource(spotify.SOURCE_NAME, "polis", false, true, false)
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("Found: %d songs.", songs.SongCount)
		for _,song := range songs.Songs {
			log.Printf("Song: %s", song.Name)
		}
	}
}
