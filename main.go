package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopePlaylistModifyPublic))
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// Authenticate with Spotify User
	ctx := context.Background()

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	//  Get list of playlists to shuffle from environment variable
	playlistsToShuffle := strings.Split(os.Getenv("SPOTIFY_PLAYLISTS"), ",")

	// Get the users playlists
	playlists, err := client.GetPlaylistsForUser(ctx, user.ID)
	if err != nil {
		log.Fatalf("error retrieve user playlists: %v", err)
	}

	// Get the IDs of the playlists to shuffle
	var playlistIDs []spotify.ID
	for _, playlist := range playlists.Playlists {
		for _, playlistToShuffle := range playlistsToShuffle {
			if playlist.Name == playlistToShuffle {
				playlistIDs = append(playlistIDs, playlist.ID)
			}
		}
	}
	
	// Shuffle the playlists
	for _, playlistID := range playlistIDs {
		
		// Get the first 100 tracks in the playlist
		tracks, err := client.GetPlaylistTracks(ctx, playlistID)
		if err != nil {
			log.Fatalf("error retrieve playlist tracks: %v", err)
		}

		// // Print the tracks in current order
		// fmt.Println("Original tracks:")
		// for _, track := range tracks.Tracks {
		// 	fmt.Println( "- ", track.Track.ID)
		// }

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

		// // Print the new track IDs
		// fmt.Println("New track IDs:")
		// for _, trackID := range newTracks {
		// 	fmt.Println( "- ", trackID)
		// }

		// Replace the tracks in the playlist with the new order
		err = client.ReplacePlaylistTracks(ctx, playlistID, newTracks...)
		if err != nil {
			log.Fatalf("error replace playlist tracks: %v", err)
		}
		
		fmt.Println("Playlist shuffled!")
	}

	fmt.Println("Done!")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}