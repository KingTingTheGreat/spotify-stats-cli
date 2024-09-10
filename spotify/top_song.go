package spotify

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"spotify-stats-cli/cnsts"
	"spotify-stats-cli/env"
	"spotify-stats-cli/types"
	"strings"
)

func refreshAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", SpotifyVars.RefreshToken)

	req, err := http.NewRequest("POST", cnsts.TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(SpotifyVars.ClientID, SpotifyVars.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse types.TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	SpotifyVars.AccessToken = tokenResponse.AccessToken
	env.WriteToEnvFile(SpotifyVars.ClientID, SpotifyVars.ClientSecret, tokenResponse.AccessToken, SpotifyVars.RefreshToken)

	return tokenResponse.AccessToken, nil
}

func fetchTopItems(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+SpotifyVars.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		newAccessToken, err := refreshAccessToken()
		if err != nil {
			return nil, err
		}
		newReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		newReq.Header.Set("Authorization", "Bearer "+newAccessToken)

		resp, err = client.Do(newReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	return io.ReadAll(resp.Body)
}

func TopSong() types.Track {
	topTracksBytes, err := fetchTopItems(cnsts.TOP_TRACKS_URL)
	if err != nil {
		log.Fatal("Failed to fetch top tracks")
	}

	topTracksResponse := types.TopTracksResponse{}
	if err := json.Unmarshal(topTracksBytes, &topTracksResponse); err != nil {
		log.Fatal("Failed to unmarshal top tracks")
	}

	tracks := topTracksResponse.Items
	topTrack := tracks[0]

	return topTrack
}
