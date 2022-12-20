package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	authConfig := &clientcredentials.Config{
		ClientID:     "<your_client_id>",
		ClientSecret: "<your_client_secret",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)
	playlistID := spotify.ID(os.Getenv("SPOTIFY_PLAYLIST_ID"))


	// Get the first 100 tracks in the playlist
	tracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		log.Fatalf("error retrieve playlist tracks: %v", err)
	}

	// Print the tracks in current order
	fmt.Println("Original tracks:")
	for _, track := range tracks.Tracks {
		fmt.Println( "- ", track.Track.ID)
	}

	// Create a new slice of tracks ID to reorder
	var newTracks []spotify.ID
	for _, track := range tracks.Tracks {
		newTracks = append(newTracks, track.Track.ID)
	}

	// Shuffle the tracks
	for i := range newTracks {
		j := rand.Intn(i + 1)
		newTracks[i], newTracks[j] = newTracks[j], newTracks[i]
	}

	// Print the new track IDs
	fmt.Println("New track IDs:")
	for _, trackID := range newTracks {
		fmt.Println( "- ", trackID)
	}

	// Replace the tracks in the playlist with the new order
	err = client.ReplacePlaylistTracks(playlistID, newTracks...)
	if err != nil {
		log.Fatalf("error replace playlist tracks: %v", err)
	}
	
	fmt.Println("Playlist shuffled!")
}