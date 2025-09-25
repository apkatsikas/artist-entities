package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apkatsikas/artist-entities/interfaces/mocks"
	"github.com/apkatsikas/artist-entities/viewmodels"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	authController *AuthController
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (suite *AuthControllerTestSuite) TestLogin() {
	expectedStatus := http.StatusOK
	expectedToken := "foo"

	userName := "username"
	password := "password"

	authService := mocks.NewIAuthService(suite.T())
	authService.EXPECT().GenerateJWT(userName, password).Return(expectedToken, nil)

	suite.authController = &AuthController{AuthService: authService}

	loginJSON, err := json.Marshal(&viewmodels.UserVM{
		UserName: userName,
		Password: password,
	})
	require.NoError(suite.T(), err)

	w := suite.doRequest(loginJSON)

	var jwt string
	json.NewDecoder(w.Body).Decode(&jwt)

	assert.Equal(suite.T(), expectedToken, jwt)
	assert.Equal(suite.T(), expectedStatus, w.Result().StatusCode)
}

func (suite *AuthControllerTestSuite) TestLoginBadPayload() {
	expectedStatus := http.StatusBadRequest

	userName := "username"
	password := "password"

	suite.authController = &AuthController{}

	loginJSON, err := json.Marshal(map[string]string{
		"user": userName,
		"pass": password,
	})

	require.NoError(suite.T(), err)

	w := suite.doRequest(loginJSON)

	var jwt string
	json.NewDecoder(w.Body).Decode(&jwt)

	assert.Equal(suite.T(), expectedStatus, w.Result().StatusCode)
}

func (suite *AuthControllerTestSuite) TestLoginIncompletePayload() {
	expectedStatus := http.StatusBadRequest

	userName := ""
	password := ""

	suite.authController = &AuthController{}

	loginJSON, err := json.Marshal(&viewmodels.UserVM{
		UserName: userName,
		Password: password,
	})

	require.NoError(suite.T(), err)

	w := suite.doRequest(loginJSON)

	var jwt string
	json.NewDecoder(w.Body).Decode(&jwt)

	assert.Equal(suite.T(), expectedStatus, w.Result().StatusCode)
}

func (suite *AuthControllerTestSuite) TestLoginUnauthorized() {
	expectedStatus := http.StatusUnauthorized

	userName := "user"
	password := "pass"

	authService := mocks.NewIAuthService(suite.T())
	authService.EXPECT().GenerateJWT(userName, password).Return("", fmt.Errorf("no"))

	suite.authController = &AuthController{AuthService: authService}

	loginJSON, err := json.Marshal(&viewmodels.UserVM{
		UserName: userName,
		Password: password,
	})
	require.NoError(suite.T(), err)

	w := suite.doRequest(loginJSON)

	var jwt string
	json.NewDecoder(w.Body).Decode(&jwt)

	assert.Equal(suite.T(), "", jwt)
	assert.Equal(suite.T(), expectedStatus, w.Result().StatusCode)
}

func (suite *AuthControllerTestSuite) doRequest(payload []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, LOGIN, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.HandleFunc(LOGIN, suite.authController.Login)
	r.ServeHTTP(w, req)
	return w
}
