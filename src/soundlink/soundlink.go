package soundlink

import (
	"errors"
	"log"
)

// SearchTypes
type SearchType int

// The different types
const (
	SearchAlbum  SearchType = 1 << iota
	SearchArtist SearchType = 1 << iota
	SearchTrack  SearchType = 1 << iota
)

const (
	TAG = "SoundLinkMaster: "
)

type SoundLinkMaster struct {
	sources map[string]*SoundLinkSource
	Nodebag *NodeBag
}

type SoundLinkSource struct {
	Source string
	Search func(query string, searchtype SearchType) (*SongResult, error)
}

type SongResult struct {
	Source    string
	Songs     []Song
	SongCount int
}

type Song struct {
	Name   string
	Artist []Artist
}

type Artist struct {
	Name string
}

func (slm *SoundLinkMaster) RegisterSource(sourceName string) *SoundLinkSource {
	log.Printf("%sRegistering %s", TAG, sourceName)
	source := &SoundLinkSource{Source: sourceName}
	slm.sources[sourceName] = source
	return source
}

func (slm *SoundLinkMaster) Search(query string, st SearchType) ([]*SongResult, error) {
	result := make([]*SongResult, 0)
	for _, source := range slm.sources {
		if res, err := source.Search(query, st); err == nil {
			result = append(result, res)
		} else {
			return result, nil
		}
	}
	return result, nil
}

func (slm *SoundLinkMaster) SearchSpecific(sourceName, query string, searchtype SearchType) (*SongResult, error) {
	if val, ok := slm.sources[sourceName]; ok {
		return val.Search(query, searchtype)
	}
	return nil, errors.New(TAG + "No registered source with that name.")
}

func New() *SoundLinkMaster {
	log.Printf("%sCreating new master.", TAG)
	return &SoundLinkMaster{
		sources: make(map[string]*SoundLinkSource),
		Nodebag: NewNodeBag(),
	}
}
