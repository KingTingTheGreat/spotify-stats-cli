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

func exchangeCodeForToken(code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", cnsts.REDIRECT_URI)

	req, err := http.NewRequest("POST", cnsts.TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth(SpotifyVars.ClientID, SpotifyVars.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenResponse types.TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	SpotifyVars.AccessToken = tokenResponse.AccessToken
	SpotifyVars.RefreshToken = tokenResponse.RefreshToken
	env.WriteToEnvFile(SpotifyVars.ClientID, SpotifyVars.ClientSecret, tokenResponse.AccessToken, tokenResponse.RefreshToken)

	return nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("callbackHandler")
	defer func() { handler_chan <- true }()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
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
