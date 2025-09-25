package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	ce "github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/interfaces/mocks"
	"github.com/apkatsikas/artist-entities/models"
	"github.com/apkatsikas/artist-entities/viewmodels"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

const (
	weirdError  = "weird error"
	artistRoute = "/artist"
	randomRoute = "/random"
	token       = "asdijsdfu23r329"
)

var authHeader = fmt.Sprintf("Bearer %v", token)

func getArtist(id string) *http.Request {
	return httptest.NewRequest(http.MethodGet,
		fmt.Sprintf(
			"%v/%v", artistRoute, id), nil)
}

func RandomGET() *http.Request {
	return httptest.NewRequest(http.MethodGet,
		fmt.Sprintf(
			"%v%v", artistRoute, randomRoute), nil)
}

func postArtist(body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPost, artistRoute, body)
}

func TestGetArtist(t *testing.T) {
	// Artist data
	artistName := "Lou Reed"
	artistID := uint(1)
	serviceRecord := models.Artist{Name: artistName}
	serviceRecord.ID = artistID

	// Expectations
	expectedArtist := viewmodels.ArtistVM{}
	expectedArtist.Name = artistName
	expectedArtist.ID = artistID
	expectedStatus := http.StatusOK

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().Get(artistID).Return(&serviceRecord, nil)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Convert data to string
	sID := strconv.FormatUint(uint64(artistID), 10)

	// Make the request
	req := getArtist(sID)
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(ARTIST_RP, artistController.Get)
	r.ServeHTTP(w, req)

	// Decode result
	artistResult := viewmodels.ArtistVM{}
	json.NewDecoder(w.Body).Decode(&artistResult)

	// Check the artist
	assert.Equal(t, expectedArtist, artistResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetArtistNoRecord(t *testing.T) {
	// Artist data
	artistID := uint(2)

	// Expectations
	expectedResponseError := ResponseError{}
	expectedResponseError.Message = ce.ErrRecordNotFound.Error()
	expectedStatus := http.StatusNotFound

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().Get(artistID).Return(nil, ce.ErrRecordNotFound)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Convert data to string
	sID := strconv.FormatUint(uint64(artistID), 10)

	// Make the request
	req := getArtist(sID)
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(ARTIST_RP, artistController.Get)
	r.ServeHTTP(w, req)

	// Decode result
	responseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&responseErrorResult)

	// Check the response error
	assert.Equal(t, expectedResponseError, responseErrorResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetArtistUnexpectedError(t *testing.T) {
	// Artist data
	artistID := uint(33)

	// Expectations
	expectedResponseError := ResponseError{}
	expectedResponseError.Message = UNEXPECTED_ERROR
	expectedStatus := http.StatusInternalServerError

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().Get(artistID).Return(nil, errors.New(weirdError))

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Convert data to string
	sID := strconv.FormatUint(uint64(artistID), 10)

	// Make the request
	req := getArtist(sID)
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(ARTIST_RP, artistController.Get)
	r.ServeHTTP(w, req)

	// Decode result
	responseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&responseErrorResult)

	// Check the response error
	assert.Equal(t, expectedResponseError, responseErrorResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetArtistIDBadData(t *testing.T) {
	var testData = []struct {
		name  string
		value string
	}{
		{name: "too big", value: "92233720368599999999"},
		{name: "negative", value: "-12"},
		{name: "float", value: "1.25"},
		{name: "NaN", value: "whatisthis"},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			// Expectations
			expectedResponseError := ResponseError{}
			expectedResponseError.Message = BAD_REQUEST
			expectedStatus := http.StatusBadRequest

			// Setup mock service
			artistService := mocks.NewIArtistService(t)

			// Inject controller with service
			artistController := ArtistController{ArtistService: artistService}

			// Make the request
			req := getArtist(tt.value)
			w := httptest.NewRecorder()
			r := chi.NewRouter()
			r.HandleFunc(ARTIST_RP, artistController.Get)
			r.ServeHTTP(w, req)

			// Decode result
			responseErrorResult := ResponseError{}
			json.NewDecoder(w.Body).Decode(&responseErrorResult)

			// Check the response error
			assert.Equal(t, expectedResponseError, responseErrorResult)
			// Check the status code
			assert.Equal(t, expectedStatus, w.Result().StatusCode)
		})
	}
}

func TestCreateArtistAuthError(t *testing.T) {
	expectedStatus := http.StatusUnauthorized

	// Artist data
	artistName := "James Brown"
	artist := viewmodels.ArtistVM{Name: artistName}

	// Setup request
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(artist)
	req := postArtist(&buf)

	// Inject controller with service
	artistController := ArtistController{ArtistService: mocks.NewIArtistService(t)}

	// Make the request
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetArtistNoInput(t *testing.T) {
	// Expectations
	expectedStatus := http.StatusNotFound

	// Data
	noArtist := ""

	// Setup mock service
	artistService := mocks.NewIArtistService(t)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Make the request
	req := getArtist(noArtist)
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(ARTIST_RP, artistController.Get)
	r.ServeHTTP(w, req)

	// Decode result
	responseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&responseErrorResult)

	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtist(t *testing.T) {
	// Artist data
	artistName := "James Brown"
	artistID := uint(1)
	vmArtist := viewmodels.ArtistVM{Name: artistName}
	serviceRecord := models.Artist{Name: artistName}
	serviceRecord.ID = artistID

	// Setup request
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(vmArtist)
	req := postArtist(&buf)
	req.Header.Add("Authorization", authHeader)

	// Expectations
	expectedArtist := viewmodels.ArtistVM{}
	expectedArtist.Name = artistName
	expectedArtist.ID = artistID
	expectedStatus := http.StatusCreated

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().Create(vmArtist.Name).Return(&serviceRecord, nil)

	authService := mocks.NewIAuthService(t)
	authService.EXPECT().IsAuthorized(token).Return(true)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService, AuthService: authService}

	// Make the request
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	// Decode result
	artistResult := viewmodels.ArtistVM{}
	json.NewDecoder(w.Body).Decode(&artistResult)

	// Check the artist
	assert.Equal(t, expectedArtist, artistResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtistRejected(t *testing.T) {
	var testData = []struct {
		name  string
		value string
		err   error
	}{
		{name: "too big",
			value: "insanelylongnamecanubelievethistrulyincrediblewhatkindoflamebandwouldhavethis",
			err:   ce.ErrDataTooLong,
		},
		{name: "invalid", value: "wow,wow", err: ce.ErrDataInvalid},
		{name: "already exists", value: "Lou Reed", err: ce.ErrRecordExists},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			// Artist
			artist := viewmodels.ArtistVM{Name: tt.name}
			// Expectations
			expectedResponseError := ResponseError{}
			expectedResponseError.Message = tt.err.Error()
			expectedStatus := http.StatusBadRequest

			// Setup request
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(artist)
			req := postArtist(&buf)
			req.Header.Add("Authorization", authHeader)

			// Setup mock service
			artistService := mocks.NewIArtistService(t)
			artistService.EXPECT().Create(artist.Name).Return(nil, tt.err)
			authService := mocks.NewIAuthService(t)
			authService.EXPECT().IsAuthorized(token).Return(true)

			// Inject controller with service
			artistController := ArtistController{ArtistService: artistService, AuthService: authService}

			// Make the request
			w := httptest.NewRecorder()
			r := chi.NewRouter()
			r.HandleFunc(POST_ARTIST_RP, artistController.Create)
			r.ServeHTTP(w, req)

			// Decode result
			responseErrorResult := ResponseError{}
			json.NewDecoder(w.Body).Decode(&responseErrorResult)

			// Check the response error
			assert.Equal(t, expectedResponseError, responseErrorResult)
			// Check the status code
			assert.Equal(t, expectedStatus, w.Result().StatusCode)
		})
	}
}

func TestCreateArtistUnexpectedError(t *testing.T) {
	// Artist data
	artistName := "James Brown"
	artist := viewmodels.ArtistVM{Name: artistName}

	// Expectations
	expectedReponseError := ResponseError{}
	expectedReponseError.Message = UNEXPECTED_ERROR
	expectedStatus := http.StatusInternalServerError

	// Setup request
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(artist)
	req := postArtist(&buf)
	req.Header.Add("Authorization", authHeader)

	// Setup mock service
	returnError := errors.New(weirdError)
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().Create(artist.Name).Return(nil, returnError)
	authService := mocks.NewIAuthService(t)
	authService.EXPECT().IsAuthorized(token).Return(true)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService, AuthService: authService}

	// Make the request
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	// Decode result
	reponseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&reponseErrorResult)

	// Check the response error
	assert.Equal(t, expectedReponseError, reponseErrorResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtistUnauthorized(t *testing.T) {
	artistName := "James Brown"
	artist := viewmodels.ArtistVM{Name: artistName}

	expectedReponseError := ResponseError{}
	expectedReponseError.Message = UNAUTHORZIED
	expectedStatus := http.StatusUnauthorized

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(artist)
	req := postArtist(&buf)
	req.Header.Add("Authorization", authHeader)

	authService := mocks.NewIAuthService(t)
	authService.EXPECT().IsAuthorized(token).Return(false)

	artistController := ArtistController{AuthService: authService}

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	reponseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&reponseErrorResult)

	assert.Equal(t, expectedReponseError, reponseErrorResult)
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtistMissingAuth(t *testing.T) {
	artistName := "James Brown"
	artist := viewmodels.ArtistVM{Name: artistName}

	expectedStatus := http.StatusUnauthorized

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(artist)
	req := postArtist(&buf)

	artistController := ArtistController{}

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	reponseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&reponseErrorResult)

	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtistMalformedAuth(t *testing.T) {
	artistName := "James Brown"
	artist := viewmodels.ArtistVM{Name: artistName}

	expectedStatus := http.StatusUnauthorized

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(artist)
	req := postArtist(&buf)
	req.Header.Add("Authorization", "blooper")

	artistController := ArtistController{}

	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	reponseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&reponseErrorResult)

	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestCreateArtistBadRequest(t *testing.T) {
	// Bad data
	badData := true

	// Expectations
	expectedReponseError := ResponseError{}
	expectedReponseError.Message = BAD_REQUEST
	expectedStatus := http.StatusBadRequest

	// Setup request
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(badData)
	req := postArtist(&buf)
	req.Header.Add("Authorization", authHeader)

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	authService := mocks.NewIAuthService(t)
	authService.EXPECT().IsAuthorized(token).Return(true)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService, AuthService: authService}

	// Make the request
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(POST_ARTIST_RP, artistController.Create)
	r.ServeHTTP(w, req)

	// Decode result
	reponseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&reponseErrorResult)

	// Check the response error
	assert.Equal(t, expectedReponseError, reponseErrorResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetRandomArtist(t *testing.T) {
	// Artist data
	artistName := "Lou Reed"
	artistID := uint(1)
	serviceRecord := models.Artist{Name: artistName}
	serviceRecord.ID = artistID

	// Expectations
	expectedArtist := viewmodels.ArtistVM{}
	expectedArtist.Name = artistName
	expectedArtist.ID = artistID
	expectedStatus := http.StatusOK

	// Setup mock service
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().GetRandom().Return(&serviceRecord, nil)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Make the request
	req := RandomGET()
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(RANDOM_ARTIST_RP, artistController.GetRandom)
	r.ServeHTTP(w, req)

	// Decode result
	artistResult := viewmodels.ArtistVM{}
	json.NewDecoder(w.Body).Decode(&artistResult)

	// Check the artist
	assert.Equal(t, expectedArtist, artistResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}

func TestGetRandomArtistUnexpectedError(t *testing.T) {
	// Expectations
	expectedResponseError := ResponseError{}
	expectedResponseError.Message = UNEXPECTED_ERROR
	expectedStatus := http.StatusInternalServerError

	// Setup mock service
	returnError := errors.New(weirdError)
	artistService := mocks.NewIArtistService(t)
	artistService.EXPECT().GetRandom().Return(nil, returnError)

	// Inject controller with service
	artistController := ArtistController{ArtistService: artistService}

	// Make the request
	req := RandomGET()
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(RANDOM_ARTIST_RP, artistController.GetRandom)
	r.ServeHTTP(w, req)

	// Decode result
	responseErrorResult := ResponseError{}
	json.NewDecoder(w.Body).Decode(&responseErrorResult)

	// Check the value
	assert.Equal(t, expectedResponseError, responseErrorResult)
	// Check the status code
	assert.Equal(t, expectedStatus, w.Result().StatusCode)
}
