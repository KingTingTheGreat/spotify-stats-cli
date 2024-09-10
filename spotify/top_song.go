package spotify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"spotify-stats-cli/cnsts"
	"spotify-stats-cli/env"
	"spotify-stats-cli/types"
	"spotify-stats-cli/util"
	"strings"
)

func refreshAccessToken() string {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", SpotifyVars.RefreshToken)

	req, err := http.NewRequest("POST", cnsts.TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		util.EndWithErr("cannot create Spotify refresh token request")
	}

	req.SetBasicAuth(SpotifyVars.ClientID, SpotifyVars.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.EndWithErr("cannot send request to Spotify refresh token URL")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.EndWithErr("cannot read response body from Spotify refresh token URL")
	}

	var tokenResponse types.TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		util.EndWithErr("cannot unmarshal response body from Spotify refresh token URL")
	}

	SpotifyVars.AccessToken = tokenResponse.AccessToken
	env.WriteToEnvFile(SpotifyVars.ClientID, SpotifyVars.ClientSecret, tokenResponse.AccessToken, SpotifyVars.RefreshToken)

	return tokenResponse.AccessToken
}

func topItemsBytes(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.EndWithErr("cannot create Spotify top tracks request")
	}

	req.Header.Set("Authorization", "Bearer "+SpotifyVars.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.EndWithErr("cannot send request to Spotify top tracks URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		newAccessToken := refreshAccessToken()
		newReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			util.EndWithErr("cannot create Spotify top tracks request")
		}
		newReq.Header.Set("Authorization", "Bearer "+newAccessToken)

		resp, err = client.Do(newReq)
		if err != nil {
			util.EndWithErr("cannot send request to Spotify top tracks URL")
		}
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.EndWithErr("cannot read response body from Spotify top tracks URL")
	}

	return body
}

func TopSong() types.Track {
	topTracksBytes := topItemsBytes(cnsts.TOP_TRACKS_URL)

	topTracksResponse := types.TopTracksResponse{}
	if err := json.Unmarshal(topTracksBytes, &topTracksResponse); err != nil {
		util.EndWithErr("cannot unmarshal top tracks")
	}

	tracks := topTracksResponse.Items
	topTrack := tracks[0]

	return topTrack
}
