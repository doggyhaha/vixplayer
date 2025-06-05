package vixplayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type VixPlayer struct {
	baseURL   string
	httpClient *http.Client
}

// Creates a new VixPlayer instance with the provided base URL and HTTP client.
func NewVixPlayer(baseURL string, client *http.Client) *VixPlayer {
	if client == nil {
		client = &http.Client{}
	}
	return &VixPlayer{
		baseURL:   baseURL,
		httpClient: client,
	}
}

// Converts the js object for the master playlist into a valid JSON object.
func fixMasterObject(jsObject []byte) []byte {
	// ugly but we can manually fix the js object to be valid JSON since the object is always the same.
	jsObject = bytes.Replace(jsObject, []byte("params:"), []byte("\"params\":"), 1)
	jsObject = bytes.Replace(jsObject, []byte("url:"), []byte("\"url\":"), 1)
	jsObject = bytes.ReplaceAll(jsObject, []byte("'"), []byte("\""))

	// Now we must remove the trailing comma only if it's the last kv pair in the object.
	var buffer bytes.Buffer
	for i := 0; i < len(jsObject); i++ {
		if jsObject[i] == ',' { // Check if the next non-whitespace character is a closing brace
			idx := i + 1
			for idx < len(jsObject) && (jsObject[idx] == ' ' || jsObject[idx] == '\t' || jsObject[idx] == '\n' || jsObject[idx] == '\r') {
				idx++
			}
			// If the next non-whitespace char is '}', skip this comma
			if idx < len(jsObject) && jsObject[idx] == '}' {
				continue
			}
		}
		buffer.WriteByte(jsObject[i])
	}

	return buffer.Bytes()
}

// Private method to retrieve HLS stream from the player URL.
func (vp *VixPlayer) getHLSFromPlayer(playerURL string, lang string) (*PlayerData, error) {
	resp, err := vp.httpClient.Get(playerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, &ContentNotFoundError{}
		}
		return nil, fmt.Errorf("failed to get player page: %s", resp.Status)
	}
	reVideo := regexp.MustCompile(`(?s)window\.video\s*=\s*(\{.*?\});`)
	reStreams := regexp.MustCompile(`(?s)window\.streams\s*=\s*(\[\{.*?\}\]);`)
	reMaster := regexp.MustCompile(`(?s)window\.masterPlaylist\s*=\s*(\{.*\})\s`)
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	videoMatch := reVideo.FindSubmatch(bodyBytes)
	if videoMatch == nil {
		return nil, fmt.Errorf("video data not found in response")
	}
	streamsMatch := reStreams.FindSubmatch(bodyBytes)
	if streamsMatch == nil {
		return nil, fmt.Errorf("streams data not found in response")
	}
	masterMatch := reMaster.FindSubmatch(bodyBytes)
	if masterMatch == nil {
		return nil, fmt.Errorf("master playlist data not found in response")
	}
	// Unmarshal the matched JSON data into the PlayerData struct.
	playerData := &PlayerData{}
	if err := json.Unmarshal(videoMatch[1], &playerData.Video); err != nil {
		return nil, fmt.Errorf("failed to unmarshal video data: %w", err)
	}
	if err := json.Unmarshal(streamsMatch[1], &playerData.Streams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal streams data: %w", err)
	}
	master := fixMasterObject(masterMatch[1])
	if err := json.Unmarshal(master, &playerData.MasterPlaylist); err != nil {
		fmt.Println("Failed to unmarshal master playlist data:", string(master))
		return nil, fmt.Errorf("failed to unmarshal master playlist data: %w", err)
	}

	// Set the MasterURL field based on the master playlist URL.
	if lang == "" {
		lang = "en"
	}
	playerData.MasterURL = playerData.MasterPlaylist.URL + "?token=" + playerData.MasterPlaylist.Params.Token + "&expires=" + playerData.MasterPlaylist.Params.Expires + "&h=1" + "&lang=" + lang
	// Append the ASN parameter if it exists.
	if playerData.MasterPlaylist.Params.ASN != "" {
		playerData.MasterURL += "&asn=" + playerData.MasterPlaylist.Params.ASN
	}
	return playerData, nil
}

// Public method to get the HLS stream for a movie by its ID.
// The lang parameter can be any language code supported by the api, it can also be empty.
func (vp *VixPlayer) GetMovieHLS(tmdbID string, lang string) (*PlayerData, error) {
	playerURL := fmt.Sprintf("%s/movie/%s", vp.baseURL, tmdbID)
	return vp.getHLSFromPlayer(playerURL, lang)
}

// Public method to get the HLS stream for a show by its ID, season, and episode.
// The lang parameter can be any language code supported by the api, it can also be empty.
func (vp *VixPlayer) GetShowHLS(tmdbID string, season int, episode int, lang string) (*PlayerData, error) {
	playerURL := fmt.Sprintf("%s/tv/%s/%d/%d", vp.baseURL, tmdbID, season, episode)
	fmt.Println("Player URL:", playerURL)
	return vp.getHLSFromPlayer(playerURL, lang)
}