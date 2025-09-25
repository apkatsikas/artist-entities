package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/viewmodels"
)

type AuthController struct {
	AuthService interfaces.IAuthService
}

func (ac *AuthController) Login(res http.ResponseWriter, req *http.Request) {
	var user viewmodels.UserVM
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	decodeError := decoder.Decode(&user)

	if decodeError != nil {
		handleRes(
			res,
			ResponseError{Message: BAD_REQUEST},
			http.StatusBadRequest,
		)
	} else if user.Password == "" || user.UserName == "" {
		handleRes(
			res,
			ResponseError{Message: BAD_REQUEST},
			http.StatusBadRequest,
		)
	} else {
		jwt, err := ac.AuthService.GenerateJWT(user.UserName, user.Password)

		if err != nil {
			handleRes(
				res,
				ResponseError{Message: UNAUTHORZIED},
				http.StatusUnauthorized,
			)
		} else {
			handleRes(res, jwt, http.StatusOK)
		}
	}
}
