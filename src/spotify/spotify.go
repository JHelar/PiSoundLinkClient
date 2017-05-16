package spotify

import (
	"soundlink"
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"io/ioutil"
)

type SpotifyClient struct {
	*url.URL
	*http.Client
}

type SpotifyErrorResponse struct {
	Status int `json:"status, omitempty"`
	Message string `json:"message, omitempty"`
}

type SpotifySearchResponse struct {
	Href string `json:"href, omitempty"`
	Limit int `json:"limit, omitempty"`
	Next int `json:"next, omitempty"`
	Offset int `json:"offset, omitempty"`
	Previous int `json:"previous, omitempty"`
	Total int `json:"total, omitempty"`
	SpotifyErrorResponse
}

type SpotifyTrackSearchResponse struct {
	SpotifySearchResponse
	Tracks []SpotifyTrack `json:"items"`
}

type SpotifyTrack struct {
	Album SpotifyAlbum `json:"album"`
	Artists []SpotifyArtist `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	DiscNumber int `json:"disc_number"`
	Duration int `json:"duration"`
	Explicit bool `json:"explicit"`
	ExternalID map[string]string `json:"external_id"`
	ExternalURL map[string]string `json:"external_url"`
	Href string `json:"href"`
	ID string `json:"id"`
	IsPlayable bool `json:"is_playable"`
	LinkedFrom SpotifyLinkedTrack `json:"linked_from"`
	Name string `json:"name"`
	Popularity int `json:"popularity"`
	PreviewUrl string `json:"preview_url"`
	TrackNumber int `json:"track_number"`
	Type string `json:"type"`
	URI string `json:"uri"`
}

type SpotifyAlbum struct {
	AlbumType string `json:"album_type"`
	Artists []SpotifyArtist `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	ExternalURLS map[string]string `json:"external_urls"`
	Href string `json:"href"`
	ID string `json:"id"`
	Images []interface{} `json:"images"`
	Name string `json:"name"`
	URI string `json:"uri"`
}

type SpotifyArtist struct {
	ExternalURLS map[string]string `json:"external_urls"`
	Followers interface{} `json:"followers, omitempty"`
	Genres []string `json:"genres, omitempty"`
	Href string `json:"href"`
	ID string `json:"id"`
	Images []interface{} `json:"images, omitempty"`
	Name string `json:"name"`
	Popularity int `json:"popularity, omitempty"`
	Type string `json:"type"`
	URI string `json:"uri"`
}

type SpotifyLinkedTrack struct {
	ExternalURLS map[string]string `json:"external_urls"`
	Href string `json:"href"`
	ID string `json:"id"`
	Type string `json:"type"`
	URI string `json:"uri"`
}

func (sc *SpotifyClient) SearchTrack(query string) (*SpotifyTrackSearchResponse, error)  {
	url := fmt.Sprintf("/search?type=track&q=%s", url.QueryEscape(query))
	var tracks *SpotifyTrackSearchResponse
	var r interface{}

	if byte, err := sc.request(http.MethodGet, url); err == nil {
		err = json.Unmarshal(byte, &tracks)
		json.Unmarshal(byte, &r)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}else{
			log.Printf("%v", r)
			return tracks, nil
		}
	}else{
		log.Fatal(err)
		return nil, err
	}
}

func (sc *SpotifyClient) request(method, url string) ([]byte, error) {
	method = strings.ToUpper(method)
	url = strings.Join([]string{sc.URL.String(), url}, "")

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := sc.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

const (
	BASE_ADDRESS = "https://api.spotify.com/v1"
	SOURCE_NAME = "Spotify"
	CLIENT_ID = ""
	CLIENT_SECRET = ""
)



var source *soundlink.SoundLinkSource
var client *SpotifyClient

func search(query string, artist, track, album bool)(*soundlink.SongResult, error){
	result := &soundlink.SongResult{
			SongCount: 0,
			Songs: make([]*soundlink.Song, 0),
		}
	if(track){
		if response, err := client.SearchTrack(query); err == nil {
			result.SongCount += response.Total

			for _, track := range response.Tracks {
				result.Songs = append(result.Songs, &soundlink.Song{Name:track.Name})
			}
		}else{
			return nil, err
		}
	}
	return result, nil
}

func Register(master *soundlink.SoundLinkMaster) error  {
	source = master.RegisterSource(SOURCE_NAME)
	source.Search = search

	url, err := url.Parse(BASE_ADDRESS)
	if err != nil {
		return err
	}
	client = &SpotifyClient{url, &http.Client{}}
	return  nil
}

