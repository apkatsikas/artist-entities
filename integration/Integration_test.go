package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/apkatsikas/artist-entities/goclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func client() *goclient.BackendClient {
	return goclient.New(os.Getenv("BASE_URL"))
}

func TestArtist(t *testing.T) {
	// Setup data
	const expectedStatusCode = http.StatusOK
	const artistID = "1"

	// Setup client and make request
	res, err := client().GetArtist(artistID)

	// Check for no errors
	require.NoErrorf(t, err, "Got an error when calling /artist: %q", err)
	// Check status code
	assert.Equal(t, expectedStatusCode, res.StatusCode)
	// Check artist
	assert.NotEmpty(t, res.Artist.Name)
	assert.NotZero(t, res.Artist.ID)
}

func TestCreateArtist(t *testing.T) {
	// Setup data
	const expectedStatusCode = http.StatusCreated
	now := time.Now().Unix()
	name := fmt.Sprintf("testart%v", now)

	client := client()

	jwtToken, err := client.Login("admin", "password")
	require.NoError(t, err)
	client.JwtToken = jwtToken

	res, err := client.CreateArtist(name)

	// Check for no errors
	require.NoErrorf(t, err, "Got an error when sending data to /artist: %q", err)
	// Check status code
	assert.Equal(t, expectedStatusCode, res.StatusCode)
	// Check artist
	assert.Equal(t, name, res.Artist.Name)
	assert.NotZero(t, res.Artist.ID)
}

func TestArtistRandom(t *testing.T) {
	// Setup data
	const expectedStatusCode = http.StatusOK

	// Setup client and make request
	res, err := client().GetArtistRandom()

	// Check for no errors
	require.NoErrorf(t, err, "Got an error when calling /artist/random: %q", err)
	// Check status code
	assert.Equal(t, expectedStatusCode, res.StatusCode)
	// Check artist
	assert.NotEmpty(t, res.Artist.Name)
	assert.NotZero(t, res.Artist.ID)
}
