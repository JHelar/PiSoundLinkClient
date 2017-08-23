package spotify

import (
	"context"
	"soundlink"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	SourceName   = "Spotify"
	ClientID     = "62ae08f2454740d7b8d3f5e2a03b92e9"
	ClientSecret = "4221fec8c51b4704b1cbd436c0024d37"
)

var source *soundlink.SoundLinkSource
var client spotify.Client

func search(query string, searchtype soundlink.SearchType) (*soundlink.SongResult, error) {
	result := &soundlink.SongResult{
		SongCount: 0,
		Songs:     make([]soundlink.Song, 0),
	}
	if searchtype == soundlink.SearchTrack {
		if response, err := client.Search(query, spotify.SearchTypeTrack); err == nil {
			result.SongCount += response.Tracks.Total

			for _, track := range response.Tracks.Tracks {
				artists := make([]soundlink.Artist, 0)
				for _, artist := range track.Artists {
					artists = append(artists, soundlink.Artist{Name: artist.Name})
				}
				result.Songs = append(result.Songs, soundlink.Song{Name: track.Name, Artist: artists})
			}
		} else {
			return nil, err
		}
	}
	return result, nil
}

func Register(master *soundlink.SoundLinkMaster) error {
	source = master.RegisterSource(SourceName)
	source.Search = search

	config := &clientcredentials.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		return err
	}
	client = spotify.Authenticator{}.NewClient(token)

	return nil
}
