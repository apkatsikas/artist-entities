package services

import (
	"fmt"
	"testing"

	"github.com/apkatsikas/artist-entities/interfaces/mocks"
	"github.com/apkatsikas/artist-entities/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	userName       = "user"
	password       = "password"
	hashedPassword = "$2a$10$FB0lrtyiqn5mCbfCFuZoPuW1vcU8QWgyuz95hMlQjUIEyubxic2h2"
)

func TestIsAuthorized(t *testing.T) {
	userRepository := mocks.NewIUserRepository(t)
	userRepository.EXPECT().Get(userName).Return(&models.User{
		Name:     userName,
		Password: hashedPassword,
	}, nil)

	service := AuthService{UserRepository: userRepository}
	service.SetJwtSigningKey(password)

	token, err := service.GenerateJWT(userName, password)
	require.NoError(t, err)

	isAuthorized := service.IsAuthorized(token)

	require.True(t, isAuthorized)
}

func TestGenerateJWTNoUser(t *testing.T) {
	userRepository := mocks.NewIUserRepository(t)
	userRepository.EXPECT().Get(userName).Return(nil, fmt.Errorf("user does not exist"))

	service := AuthService{UserRepository: userRepository}
	service.SetJwtSigningKey(password)

	token, err := service.GenerateJWT(userName, password)
	require.Error(t, err)
	require.Empty(t, token)
}

func TestGenerateJWTWrongPassword(t *testing.T) {
	userRepository := mocks.NewIUserRepository(t)
	userRepository.EXPECT().Get(userName).Return(&models.User{
		Name:     userName,
		Password: hashedPassword,
	}, nil)

	service := AuthService{UserRepository: userRepository}
	service.SetJwtSigningKey(password)

	token, err := service.GenerateJWT(userName, "bloop")
	require.Error(t, err)
	require.Empty(t, token)
}

func TestIsAuthorizedFail(t *testing.T) {
	service := AuthService{}
	service.SetJwtSigningKey(password)

	isAuthorized := service.IsAuthorized("token")
	require.False(t, isAuthorized)
}

func TestEmptyKeyIsAuthorized(t *testing.T) {
	service := AuthService{}

	require.Panics(t, func() {
		service.IsAuthorized("token")
	})
}

func TestEmptyKeyGenerateJwt(t *testing.T) {
	service := AuthService{}

	require.Panics(t, func() {
		service.GenerateJWT(userName, password)
	})
}

func TestCreateUser(t *testing.T) {
	user := &models.User{}
	userRepository := mocks.NewIUserRepository(t)
	userRepository.EXPECT().Create(userName, mock.AnythingOfType("string")).Return(user, nil)

	service := AuthService{UserRepository: userRepository}
	createdUser, err := service.CreateUser(userName, password)
	require.Nil(t, err)
	require.Equal(t, user, createdUser)
}
