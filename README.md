# Sonos Spotify Shuffle

[![Go Report Card](https://goreportcard.com/badge/joekleinsorge/sonos-spotify-shuffle)](https://goreportcard.com/report/joekleinsorge/sonos-spotify-shuffle)

A simple script to force shuffle Spotify playlists because Google Assistant on Sonos can't do it.

## How to use

1. Have a Spotify Premium account
2. Have `Go` installed locally
3. [Create a Spotify app](https://developer.spotify.com/my-applications/) and get the `ID` and `Secret`
   1. Set the redirect URI to `http://localhost:8080/callback`
4. Clone this repo `gh repo clone joekleinsorge/sonos-spotify-shuffle`
5. Create the following environment variables:
    - `SPOTIFY_ID="<Your Spotify ID>"`
    - `SPOTIFY_SECRET="<Your Spotify Secret>"`
    - `SPOTIFY_PLAYLISTS="<Your Playlist names>,<Comma separated with no spaces>"`
6. Run `go run main.go`
7. When prompted, go to the URL and login to Spotify
8. The script will then run and shuffle your playlists

## Acknowledgements

- [Spotify Go Wrapper](github.com/zmb3/spotify)
