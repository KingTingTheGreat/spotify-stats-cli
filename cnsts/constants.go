package cnsts

const (
	REDIRECT_URI    = "http://localhost:3000/callback"
	AUTH_URL        = "https://accounts.spotify.com/authorize"
	TOKEN_URL       = "https://accounts.spotify.com/api/token"
	SCOPE           = "user-top-read"
	TOP_TRACKS_URL  = "https://api.spotify.com/v1/me/top/tracks?time_range=short_term"
	TOP_ARTISTS_URL = "https://api.spotify.com/v1/me/top/artists?time_range=short_term"
	// DELIM             = "\x1d"
	DELIM             = "%##$%"
	ANSI_RESET        = "\x1b[0m"
	CHAR              = "█"
	FONT_ASPECT_RATIO = 0.46
	DIM               = 44
	OUTPUT_FILE_NAME  = "output.txt"
)
