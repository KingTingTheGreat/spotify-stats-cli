package env

import (
	"os"
	"spotify-stats-cli/types"
	"spotify-stats-cli/util"

	"github.com/joho/godotenv"
)

var env_file = ""

func WriteToEnvFile(clientId, clientSecret, accessToken, refreshToken string) {
	envVars := map[string]string{
		"SPOTIFY_CLIENT_ID":     clientId,
		"SPOTIFY_CLIENT_SECRET": clientSecret,
		"SPOTIFY_ACCESS_TOKEN":  accessToken,
		"SPOTIFY_REFRESH_TOKEN": refreshToken,
	}

	godotenv.Write(envVars, env_file)
}

func LoadEnvVars() types.SpotifyVars {
	basePath := util.BasePath()
	// fmt.Println("Base Path:", basePath)

	err := godotenv.Load(basePath + "\\.env")
	if err != nil {
		err = godotenv.Load(basePath + "/.env")
		if err != nil {
			util.EndWithErr("cannot load .env file")
		} else {
			env_file = basePath + "/.env"
		}
	} else {
		env_file = basePath + "\\.env"
	}

	return types.SpotifyVars{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		AccessToken:  os.Getenv("SPOTIFY_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		State:        randomString(),
	}
}
