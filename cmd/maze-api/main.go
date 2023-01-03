package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/egurnov/maze-api/maze-api/app"
	"github.com/egurnov/maze-api/maze-api/jwtservice"
	"github.com/egurnov/maze-api/maze-api/service"
	storepkg "github.com/egurnov/maze-api/maze-api/store"
)

type Config struct {
	// GIN_MODE=(release|debug)
	Port   int    `envconfig:"PORT" default:"8080"`
	JWTKey []byte `envconfig:"JWT_SIGNING_KEY" required:"true"`
	DBURL  string `envconfig:"DB_URL" required:"true"` // default:"root@(localhost:3306)/dreamteam"
}

func main() {
	log.Println("Server starting")

	// Configuration
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.WithError(err).Fatal("cannot read configs")
	}
	flag.Parse()

	// Create Store
	store, err := storepkg.NewMySQLStore(cfg.DBURL)
	if err != nil {
		log.WithError(err).Fatal("cannot start the DB")
	}

	// Create Logger
	logger := log.New()
	logger.SetOutput(os.Stdout)
	if os.Getenv("GIN_MODE") != "release" {
		logger.SetLevel(log.DebugLevel)
	}

	// Create JWTService
	if len(cfg.JWTKey) == 0 {
		log.WithError(err).Fatal("JWT signing key is empty")
	}
	jwtService := jwtservice.New(cfg.JWTKey, 24*time.Hour)

	// Create app
	mazeAPI := &app.App{
		Log:        logger,
		JWTService: jwtService,

		UserService: &service.UserService{Store: &storepkg.UserStore{Store: store}},
		MazeService: &service.MazeService{Store: &storepkg.MazeStore{Store: store}},
	}

	// Initialize DB if requested
	// if len(*adminEmail) > 0 || len(*adminPass) > 0 {
	// 	id, err := mazeAPI.CreateAdmin(*adminEmail, *adminPass)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("could not create admin user")
	// 	}
	// 	log.WithField("id", id).Info("admin user created")
	// }

	// Create Gin engine
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	mazeAPI.SetRoutes(engine)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: engine,
	}

	// Start server
	go func() {
		log.Infof("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("serving http")
		}
	}()

	// Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("cannot shutdown, crashing instead")
	}

	log.Info("Server exiting")
}
