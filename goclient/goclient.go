package goclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlLib "net/url"
	"path"
	"time"

	"github.com/apkatsikas/artist-entities/viewmodels"
)

const (
	timeout   = 30 * time.Second
	artistStr = "artist"
)

// ResponseMetadata represents an http response minus the body
type ResponseMetadata struct {
	StatusCode int
	Header     http.Header
}

// RawResponse represents a response from an endpoint that does not return a struct/JSON
type RawResponse struct {
	*ResponseMetadata
	Body []byte
}

// ArtistResponse represents a response from the /artist endpoint
type ArtistResponse struct {
	*ResponseMetadata
	Artist *viewmodels.ArtistVM
}

// BackendClient represents an API client for an http service
type BackendClient struct {
	baseURL    string
	httpClient *http.Client
	JwtToken   string
}

// New returns a BackendClient for a given baseURL
func New(baseURL string) *BackendClient {
	if baseURL == "" {
		panic("baseURL for BackendClient cannot be blank")
	}

	return &BackendClient{baseURL: baseURL,
		httpClient: &http.Client{Timeout: timeout}}
}

func (bc *BackendClient) buildURL(path string, qs urlLib.Values) (*urlLib.URL, error) {
	url, err := urlLib.Parse(bc.baseURL)

	if err != nil {
		return nil, err
	}

	url.Path = path
	url.RawQuery = qs.Encode()

	return url, nil
}

func (bc *BackendClient) sendRequest(url *urlLib.URL, httpMethod string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url.String(), body)
	if err != nil {
		return nil, err
	}

	res, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetArtist calls the /artist endpoint for a given id and returns the Artist
func (bc *BackendClient) GetArtist(id string) (*ArtistResponse, error) {
	// Setup our artist
	artist := viewmodels.ArtistVM{}

	// Build URL
	url, err := bc.buildURL(path.Join(artistStr, id), nil)
	if err != nil {
		return nil, err
	}

	// Send request
	res, err := bc.sendRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode and return response
	err = json.NewDecoder(res.Body).Decode(&artist)
	if err != nil {
		return nil, err
	}
	return &ArtistResponse{
		ResponseMetadata: &ResponseMetadata{StatusCode: res.StatusCode, Header: res.Header},
		Artist:           &artist}, nil
}

// GetArtistRandom calls the /artist/random endpoint and returns the Artist
func (bc *BackendClient) GetArtistRandom() (*ArtistResponse, error) {
	// Setup our artist
	artist := viewmodels.ArtistVM{}

	// Build URL
	url, err := bc.buildURL(path.Join(artistStr, "random"), nil)
	if err != nil {
		return nil, err
	}

	// Send request
	res, err := bc.sendRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode and return response
	err = json.NewDecoder(res.Body).Decode(&artist)
	if err != nil {
		return nil, err
	}
	return &ArtistResponse{
		ResponseMetadata: &ResponseMetadata{StatusCode: res.StatusCode, Header: res.Header},
		Artist:           &artist}, nil
}

func (bc *BackendClient) Login(userName string, password string) (string, error) {
	url, err := bc.buildURL("/login", nil)
	if err != nil {
		return "", err
	}

	loginJSON, err := json.Marshal(&viewmodels.UserVM{
		UserName: userName,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	res, err := bc.sendRequest(url, http.MethodPost, bytes.NewBuffer(loginJSON))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var jwt string
	err = json.NewDecoder(res.Body).Decode(&jwt)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

// CreateArtist sends data to the /artist endpoint and returns the Artist
func (bc *BackendClient) CreateArtist(name string) (*ArtistResponse, error) {
	// Setup our artist
	artist := viewmodels.ArtistVM{Name: name}

	// Build URL
	url, err := bc.buildURL(artistStr, nil)
	if err != nil {
		return nil, err
	}

	// Marshal the data
	artistJSON, err := json.Marshal(artist)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(artistJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bc.JwtToken))

	res, err := bc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Decode and return the response
	err = json.NewDecoder(res.Body).Decode(&artist)
	if err != nil {
		return nil, err
	}
	return &ArtistResponse{
		ResponseMetadata: &ResponseMetadata{StatusCode: res.StatusCode, Header: res.Header},
		Artist: &viewmodels.ArtistVM{
			Name: artist.Name,
			ID:   artist.ID,
		}}, nil
}
