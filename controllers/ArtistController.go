package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	ce "github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/infrastructures/logutil"
	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/viewmodels"
	"github.com/go-chi/chi/v5"
)

type ArtistController struct {
	ArtistService interfaces.IArtistService
	AuthService   interfaces.IAuthService
}

func (ac *ArtistController) Get(res http.ResponseWriter, req *http.Request) {
	artistID := chi.URLParam(req, "artistID")

	// Try to parse int from string param
	u64, err := strconv.ParseUint(artistID, 10, 32)
	uintID := uint(u64)
	if err != nil {
		handleRes(
			res,
			ResponseError{Message: BAD_REQUEST},
			http.StatusBadRequest,
		)
	} else {
		// Get the artist from the service
		artist, err := ac.ArtistService.Get(uintID)

		if err != nil {
			// Record not found
			if errors.Is(err, ce.ErrRecordNotFound) {
				handleRes(
					res,
					ResponseError{Message: err.Error()},
					http.StatusNotFound,
				)
			} else {
				logutil.Error("Failed to get artist with ID of %v. Error was: %v", artistID, err)
				handleRes(
					res,
					ResponseError{Message: UNEXPECTED_ERROR},
					http.StatusInternalServerError,
				)
			}
		} else {
			// Encode the artist to the response
			encodeRes(res,
				viewmodels.ArtistVM{Name: artist.Name, ID: artist.ID})
		}
	}
}
func getBearerToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	splitBySpace := strings.Split(authHeader, " ")
	if len(splitBySpace) != 2 || splitBySpace[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	token := splitBySpace[1]
	return token, nil
}

func (ac *ArtistController) Create(res http.ResponseWriter, req *http.Request) {
	token, err := getBearerToken(req)
	if err != nil {
		handleRes(
			res,
			ResponseError{Message: err.Error()},
			http.StatusUnauthorized,
		)
	} else if !ac.AuthService.IsAuthorized(token) {
		handleRes(
			res,
			ResponseError{Message: UNAUTHORZIED},
			http.StatusUnauthorized,
		)
	} else {
		var artist viewmodels.ArtistVM
		decodeError := json.NewDecoder(req.Body).Decode(&artist)

		if decodeError != nil {
			handleRes(
				res,
				ResponseError{Message: BAD_REQUEST},
				http.StatusBadRequest,
			)
		} else {
			createdArtist, err := ac.ArtistService.Create(artist.Name)

			if err != nil {
				invalid := errors.Is(err, ce.ErrDataInvalid)
				dataTooLong := errors.Is(err, ce.ErrDataTooLong)
				recordExists := errors.Is(err, ce.ErrRecordExists)

				if invalid || dataTooLong || recordExists {
					handleRes(
						res,
						ResponseError{Message: err.Error()},
						http.StatusBadRequest,
					)
				} else {
					logutil.Error("Failed to create artist %v. Error was: %v", artist.Name, err)
					handleRes(
						res,
						ResponseError{Message: UNEXPECTED_ERROR},
						http.StatusInternalServerError,
					)
				}
			} else {
				// Encode the artist to the response
				handleRes(res,
					viewmodels.ArtistVM{Name: createdArtist.Name, ID: createdArtist.ID}, http.StatusCreated)
			}
		}
	}
}

func (ac *ArtistController) GetRandom(res http.ResponseWriter, req *http.Request) {
	// Get the artist from the service
	artist, err := ac.ArtistService.GetRandom()

	if err != nil {
		logutil.Error("Failed to get random artist. Error was: %v", err)
		handleRes(
			res,
			ResponseError{Message: UNEXPECTED_ERROR},
			http.StatusInternalServerError,
		)
	} else {
		// Encode the artist to the response
		encodeRes(res,
			viewmodels.ArtistVM{Name: artist.Name, ID: artist.ID})
	}
}
