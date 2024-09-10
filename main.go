package main

import (
	"os"
	"spotify-stats-cli/env"
	"spotify-stats-cli/hndle_err"
	"spotify-stats-cli/spotify"
	"spotify-stats-cli/util"
)

func main() {
	env.LoadEnvVars()

	if spotify.SpotifyVars.ClientID == "" || spotify.SpotifyVars.ClientSecret == "" {
		hndle_err.EndWithErr("Missing client ID or client secret")
	}

	if spotify.SpotifyVars.AccessToken == "" || spotify.SpotifyVars.RefreshToken == "" {
		spotify.InitializeSpotifyAccess()
	}

	topTrack := spotify.TopSong()

	trackText := topTrack.Name + " - " + topTrack.Artists[0].Name

	// get ansi image
	ansiImage, err := util.AnsiImage(topTrack.Album.Images[0].Url)
	if err != nil {
		hndle_err.EndWithErr("Cannot get ansi image")
	}

	outputFile, err := util.WriteOutputToFile(ansiImage, trackText)
	if err != nil {
		hndle_err.EndWithErr("Cannot write to file")
	}

	os.Stdout.WriteString(outputFile)
}
