package types

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type Artist struct {
	Name string `json:"name"`
	Href string `json:"href"`
	Uri  string `json:"uri"`
}

type Album struct {
	Name        string   `json:"name"`
	Href        string   `json:"href"`
	Uri         string   `json:"uri"`
	Artists     []Artist `json:"artists"`
	ReleaseDate string   `json:"release_date"`
	TotalTracks int      `json:"total_tracks"`
	Images      []Image  `json:"images"`
}

type Track struct {
	Name        string   `json:"name"`
	Album       Album    `json:"album"`
	Artists     []Artist `json:"artists"`
	Uri         string   `json:"uri"`
	Href        string   `json:"href"`
	TrackNumber int      `json:"track_number"`
	Popularity  int      `json:"popularity"`
}

type TopTracksResponse struct {
	Items []Track `json:"items"`
	Total int     `json:"total"`
	Limit int     `json:"limit"`
}

type TopArtistsResponse struct {
	Items []Artist `json:"items"`
	Total int      `json:"total"`
	Limit int      `json:"limit"`
}

type SpotifyVars struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	State        string
}
