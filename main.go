package main

import (
	"os"
	"spotify-stats-cli/env"
	"spotify-stats-cli/spotify"
	"spotify-stats-cli/util"
)

func main() {
	env.LoadEnvVars()

	if spotify.SpotifyVars.ClientID == "" || spotify.SpotifyVars.ClientSecret == "" {
		util.EndWithErr("Missing client ID or client secret")
	}

	if spotify.SpotifyVars.AccessToken == "" || spotify.SpotifyVars.RefreshToken == "" {
		spotify.InitializeSpotifyAccess()
	}

	topTrack := spotify.TopSong()

	trackText := topTrack.Name + " - " + topTrack.Artists[0].Name

	// get ansi image
	ansiImage := util.AnsiImage(topTrack.Album.Images[0].Url)

	outputFile := util.WriteOutputToFile(ansiImage, trackText)

	os.Stdout.WriteString(outputFile)
}
