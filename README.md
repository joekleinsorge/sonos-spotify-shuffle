# Sonos Spotify Shuffle

[![Go Report Card](https://goreportcard.com/badge/joekleinsorge/sonos-spotify-shuffle)](https://goreportcard.com/report/joekleinsorge/sonos-spotify-shuffle)

A simple script to force shuffle Spotify playlists because Google Assistant running through Sonos can't do it.

## Why this exists

I have a few Sonos speakers and use Google Assistant to control my smart home. I also have a few themed Spotify playlists that we like to listen to during dinner. However, it dinner only lasts so long and we end up listening to the same few song again and again.

Now the easy solution would be to open either the Spotify or Sonos apps and just click the shuffle button. However, I don't want to have my phone out during dinner.

My preferred solution would be to just ask Google Assistant to shuffle my playlists. Unfortunately, this is currently not possible. I can ask it to play a playlist, but it will always play the playlist in order. So I decided to write this script to permanently shuffle the playlists every time we get bored.

## What this does

This script will authenticate with Spotify and then it will find the `ID`s of the playlists you added to `SPOTIFY_PLAYLIST`. It will then randomly shuffle the songs in each playlist and replace the existing playlist with the shuffled version.

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
