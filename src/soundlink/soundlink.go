package soundlink

import (
	"log"
	"errors"
)

const TAG = "SoundLinkMaster: "

type SoundLinkMaster struct {
	sources map[string]*SoundLinkSource
}

type SoundLinkSource struct {
	Source string
	Search func(query string, artist, track, album bool) (*SongResult, error)
}

type SongResult struct {
	Songs []*Song
	SongCount int
}

type Song struct {
	Name string
}

func (slm *SoundLinkMaster) RegisterSource(sourceName string) *SoundLinkSource {
	log.Printf("%sRegistering %s", TAG, sourceName)
	source := &SoundLinkSource{Source:sourceName}
	slm.sources[sourceName] = source
	return source
}

func (slm *SoundLinkMaster) SearchSpecificSource(sourceName, query string, artist, track, album bool) (*SongResult, error){
	if val, ok := slm.sources[sourceName]; ok {
		return val.Search(query, artist, track, album)
	}
	return nil, errors.New(TAG + "No registered source with that name.")
}

func New() (*SoundLinkMaster)  {
	log.Printf("%sCreating new master.", TAG)
	return &SoundLinkMaster{
		sources: make(map[string]*SoundLinkSource),
	}
}