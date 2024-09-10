package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"spotify-stats-cli/cnsts"
	"spotify-stats-cli/env"
	"spotify-stats-cli/types"
	"spotify-stats-cli/util"
	"strings"
)

var handler_chan = make(chan bool, 1)

func getAuthURL() string {
	params := url.Values{}
	params.Set("client_id", SpotifyVars.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", cnsts.REDIRECT_URI)
	params.Set("scope", cnsts.SCOPE)
	params.Set("state", SpotifyVars.State)

	return fmt.Sprintf("%s?%s", cnsts.AUTH_URL, params.Encode())
}

func exchangeCodeForToken(code string) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", cnsts.REDIRECT_URI)

	req, err := http.NewRequest("POST", cnsts.TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		util.EndWithErr("cannot create Spotify token request")
	}

	req.SetBasicAuth(SpotifyVars.ClientID, SpotifyVars.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.EndWithErr("cannot send request to Spotify token URL")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.EndWithErr("cannot read response body from Spotify token URL")
	}

	var tokenResponse types.TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		util.EndWithErr("cannot unmarshal response body from Spotify token URL")
	}

	SpotifyVars.AccessToken = tokenResponse.AccessToken
	SpotifyVars.RefreshToken = tokenResponse.RefreshToken
	env.WriteToEnvFile(SpotifyVars.ClientID, SpotifyVars.ClientSecret, tokenResponse.AccessToken, tokenResponse.RefreshToken)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("callbackHandler")
	defer func() { handler_chan <- true }()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	exchangeCodeForToken(code)
}

func InitializeSpotifyAccess() {
	fmt.Println("Visit the following URL to authorize the application:")
	fmt.Println(getAuthURL())

	http.HandleFunc("/callback", callbackHandler)
	for {
		http.ListenAndServe(":8080", nil)
		<-handler_chan
		break
	}
}
