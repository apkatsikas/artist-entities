package router

import (
	"sync"

	"github.com/apkatsikas/artist-entities/controllers"
	"github.com/apkatsikas/artist-entities/infrastructures/logutil"
	"github.com/go-chi/chi/v5"
)

type IChiRouter interface {
	InitRouter(ac *controllers.ArtistController, authController *controllers.AuthController) *chi.Mux
}

type router struct{}

func (router *router) InitRouter(ac *controllers.ArtistController,
	authController *controllers.AuthController) *chi.Mux {
	// Create router
	r := chi.NewRouter()
	r.HandleFunc(controllers.ARTIST_RP, ac.Get)
	r.HandleFunc(controllers.POST_ARTIST_RP, ac.Create)
	r.HandleFunc(controllers.RANDOM_ARTIST_RP, ac.GetRandom)

	r.HandleFunc(controllers.LOGIN, authController.Login)

	logutil.Info("Router initialized")

	return r
}

// Setup singleton
var (
	m          *router
	routerOnce sync.Once
)

func ChiRouter() IChiRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}
