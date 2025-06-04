package spotify

type SearchResponse struct {
	Limit   int                  `json:"limit"`
	Offsite int                  `json:"offset"`
	Items   []SpotifyTrackObject `json:"items"` // Assuming spotifyTrack is a struct that holds individual track details
	Total   int                  `json:"total"`
}

type SpotifyTrackObject struct {
	//album releated fields
	AlbumType        string   `json:"albumType"`      // Assuming AlbumType is a string representing the type of album
	AlbumTotalTracks int      `json:"totalTracks"`    // Assuming TotalTracks is an integer representing the total number of tracks in the album
	AlbumImagesURL   []string `json:"AlbumImagesURL"` // Assuming Images is a struct that holds image details
	AlbumName        string   `json:"albumName"`      // Name of the album

	//artist related fields
	ArtistsName []string `json:"artistsName"` // Assuming spotifyArtist is a struct that holds artist name

	//track related fields
	Explicit bool   `json:"explicit"` // Assuming Explicit is a boolean indicating if the track is explicit
	ID       string `json:"id"`       // ID of the track
	Name     string `json:"name"`     // Name of the track
}
