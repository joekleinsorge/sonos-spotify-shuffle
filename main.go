package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// Get GitHub repository secrets
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	userID := os.Getenv("SPOTIFY_USER_ID")
	playlistsToShuffle := strings.Split(os.Getenv("SPOTIFY_PLAYLISTS"), ",")

	// Set up Spotify API client
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}
	client := config.Client(context.Background())

	// Set up Spotify Web API client
	api := spotify.NewClient(client)

	// Get playlist IDs
	playlistIDs, err := getPlaylistIDs(api, userID, playlistsToShuffle)
	if err != nil {
		log.Fatal(err)
	}

	// Shuffle and update playlists
	for _, playlistID := range playlistIDs {
		err = shufflePlaylist(api, playlistID)
		if err != nil {
			log.Printf("Failed to shuffle playlist %s: %v", playlistID, err)
		} else {
			fmt.Printf("Playlist %s shuffled successfully\n", playlistID)
		}
	}
}

func getPlaylistIDs(api spotify.Client, userID string, playlistsToShuffle []string) ([]spotify.ID, error) {
	var playlistIDs []spotify.ID

	// Get user's playlists
	playlists, err := api.GetPlaylistsForUser(userID)
	if err != nil {
		return nil, err
	}

	// Loop through user's playlists and find the ones to shuffle
	for _, playlist := range playlists.Playlists {
		for _, playlistName := range playlistsToShuffle {
			if strings.Contains(playlist.Name, playlistName) {
				playlistIDs = append(playlistIDs, playlist.ID)
			}
		}
	}

	return playlistIDs, nil
}

func shufflePlaylist(api spotify.Client, playlistID spotify.ID) error {
	// Get playlist tracks
	tracks, err := api.GetPlaylistTracks(playlistID)
	if err != nil {
		return err
	}

	// Shuffle the tracks
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tracks.Tracks), func(i, j int) {
		tracks.Tracks[i], tracks.Tracks[j] = tracks.Tracks[j], tracks.Tracks[i]
	})

	// Update the playlist with shuffled tracks
	playlistShuffleOptions := spotify.PlaylistReorderOptions{
		RangeStart: 0,
		RangeLength: len(tracks.Tracks),
		InsertBefore: 0,
	}

	_, err = api.ReorderPlaylistTracks(playlistID, playlistShuffleOptions)
	if err != nil {
		return err
	}

	return nil
}

