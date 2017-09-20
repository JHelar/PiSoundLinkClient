package soundlink

import (
	"encoding/json"
	"errors"
	"log"
	"soundnodes"
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
	sources       map[string]*SoundLinkSource
	ClientNodeBag *soundnodes.NodeBag
	NodesNodeBag  *soundnodes.NodeBag
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

func (slm *SoundLinkMaster) search(query string, st SearchType) ([]*SongResult, error) {
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

func (slm *SoundLinkMaster) nodeMessageHandler(message soundnodes.NodeMessage) (soundnodes.NodeMessage, error) {
	return message, nil
}

func (slm *SoundLinkMaster) clientMessageHandler(message soundnodes.NodeMessage) (soundnodes.NodeMessage, error) {
	var requestBase soundnodes.Base
	var err error
	if err = json.Unmarshal([]byte(message.Message), &requestBase); err == nil {
		switch requestBase.Type {
		case "Search":
			if result, err := slm.search(requestBase.Data["query"], SearchTrack); err == nil {
				if bytes, err := json.Marshal(result); err == nil {
					message.Message = string(bytes[:])
				}
			}
			break
		case "Play":
			if val, ok := requestBase.Data["nodeid"]; ok {
				slm.NodesNodeBag.SendRaw(val, message.Message)
			}
		default:
			break
		}
	}
	if err != nil {
		message.Message = err.Error()
	}
	return message, nil
}

func (slm *SoundLinkMaster) SearchSpecific(sourceName, query string, searchtype SearchType) (*SongResult, error) {
	if val, ok := slm.sources[sourceName]; ok {
		return val.Search(query, searchtype)
	}
	return nil, errors.New(TAG + "No registered source with that name.")
}

func New() *SoundLinkMaster {
	log.Printf("%sCreating new master.", TAG)
	slm := &SoundLinkMaster{
		sources:       make(map[string]*SoundLinkSource),
		ClientNodeBag: soundnodes.NewNodeBag(),
		NodesNodeBag:  soundnodes.NewNodeBag(),
	}
	slm.ClientNodeBag.MessageHandler = slm.clientMessageHandler
	return slm
}
