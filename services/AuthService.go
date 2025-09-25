package services

import (
	"time"

	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	UserRepository  interfaces.IUserRepository
	jwtSignatureKey []byte
}

func fiveMinuteExpiration() time.Time {
	return time.Now().Add(5 * time.Minute)
}

func (as *AuthService) SetJwtSigningKey(signatureKey string) {
	as.jwtSignatureKey = []byte(signatureKey)
}

func (as *AuthService) IsAuthorized(token string) bool {
	as.panicIfEmptyKey()

	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return as.jwtSignatureKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// unauthorized
			return false
		}
		// bad request
		return false
	}
	if !tkn.Valid {
		// unauthorized
		return false
	}
	return true
}

func (as *AuthService) CreateUser(name string, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := as.UserRepository.Create(name, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (as *AuthService) panicIfEmptyKey() {
	if len(as.jwtSignatureKey) == 0 {
		panic("JWT signature key cannot be blank!")
	}
}

func (as *AuthService) GenerateJWT(name string, password string) (string, error) {
	as.panicIfEmptyKey()

	user, err := as.UserRepository.Get(name)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(fiveMinuteExpiration()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSignatureKey)
}
