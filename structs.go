package vixplayer

// Video holds metadata about the video file.
type Video struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Size          int    `json:"size"`
	Quality       int    `json:"quality"`
	Duration      int    `json:"duration"`
	Views         int    `json:"views"`
	IsViewable    int    `json:"is_viewable"`
	Status        string `json:"status"`
	FPS           int    `json:"fps"`
	Legacy        int    `json:"legacy"`
	FolderID      string `json:"folder_id"`
	CreatedAtDiff string `json:"created_at_diff"`
}

// Stream represents a single streaming server entry.
type Stream struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
	URL    string `json:"url"`
}

// MasterPlaylistParams holds the query parameters for the master playlist URL.
type MasterPlaylistParams struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
	ASN     string `json:"asn"`
}

// MasterPlaylist wraps the params object plus its URL.
type MasterPlaylist struct {
	Params MasterPlaylistParams `json:"params"`
	URL    string               `json:"url"`
}

// PlayerData combines all pieces (video, streams, playlist) and includes a
// top-level MasterURL string field.
type PlayerData struct {
	Video          Video          `json:"video"`
	Streams        []Stream       `json:"streams"`
	MasterPlaylist MasterPlaylist `json:"masterPlaylist"`
	MasterURL      string         `json:"masterUrl"`
}

type ContentNotFoundError struct {}

func (e *ContentNotFoundError) Error() string {
	return "content not found"
}