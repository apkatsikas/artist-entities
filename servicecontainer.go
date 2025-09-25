package entities

import (
	"os"
	"sync"

	"github.com/apkatsikas/artist-entities/controllers"
	"github.com/apkatsikas/artist-entities/infrastructures"
	"github.com/apkatsikas/artist-entities/infrastructures/fileutil"
	"github.com/apkatsikas/artist-entities/infrastructures/flagutil"
	"github.com/apkatsikas/artist-entities/infrastructures/logutil"
	"github.com/apkatsikas/artist-entities/migrate"
	"github.com/apkatsikas/artist-entities/repositories"
	"github.com/apkatsikas/artist-entities/router"
	"github.com/apkatsikas/artist-entities/services"
	"github.com/apkatsikas/artist-entities/services/rules"
	"github.com/apkatsikas/artist-entities/storageclient"
	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
)

const (
	dbFile = "entities.db"

	// 2am every day
	schedule = "0 2 * * *"
)

type IServiceContainer interface {
	Setup() *chi.Mux
}

type kernel struct {
	sqliteHandler *infrastructures.SQLiteHandler
}

func (k *kernel) Setup() *chi.Mux {
	// Setup logs
	logutil.Setup()
	logutil.Info("Running...")

	// Setup flags
	fu := flagutil.Get()
	fu.Setup()

	// Setup sqlite
	k.sqliteHandler = &infrastructures.SQLiteHandler{}

	// Connect to SQLite
	err := k.sqliteHandler.ConnectSQLite(dbFile)
	if err != nil {
		logutil.Fatal("Failed to connect to SQLite. Error was %v", err)
	}

	// Bring everything online
	storage := storageclient.New()
	artistRules := &rules.ArtistRules{}
	adminRules := &rules.AdminRules{}
	fileUtil := &fileutil.FileUtil{}

	// Inject dependencies

	// Admin
	adminRepository := &repositories.AdminRepository{IDB: k.sqliteHandler}
	adminService := &services.AdminService{
		AdminRepository: adminRepository,
		FileUtil:        fileUtil,
		StorageClient:   storage, Rules: adminRules,
	}
	userRepository := &repositories.UserRepository{IDB: k.sqliteHandler}
	authService := &services.AuthService{UserRepository: userRepository}

	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		panic("JWT_SIGNING_KEY must be set")
	}
	authService.SetJwtSigningKey(signingKey)

	// Web
	artistRepository := &repositories.ArtistRepository{IDB: k.sqliteHandler}
	artistService := &services.ArtistService{
		ArtistRepository: artistRepository,
		Rules:            artistRules,
	}

	artistController := &controllers.ArtistController{ArtistService: artistService,
		AuthService: authService}
	authController := &controllers.AuthController{AuthService: authService}

	if fu.MigrateUser != "" && fu.MigratePassword != "" {
		err := userRepository.Migrate()
		if err != nil {
			logutil.Error("Failed to migrate user table, error was %v", err)
		}
		_, err = authService.CreateUser(fu.MigrateUser, fu.MigratePassword)
		if err != nil {
			logutil.Error("Failed to create user %v, error was %v", err)
		}
		return nil
	}

	// Setup cron
	c := cron.New()

	c.AddFunc(schedule, func() {
		err = adminService.Backup()
		if err != nil {
			logutil.Error("Failed to backup entities DB, error was %v", err)
		}
	})
	c.Start()

	// Migrate
	if fu.MigrateDB {
		migrate.Migrate(artistRepository, artistService, fu.Secret)
	}

	// Setup router
	return router.ChiRouter().InitRouter(artistController, authController)
}

// Setup singleton
var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
