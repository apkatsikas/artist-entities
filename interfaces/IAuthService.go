package interfaces

type IAuthService interface {
	IsAuthorized(token string) bool
	GenerateJWT(name string, password string) (string, error)
}
