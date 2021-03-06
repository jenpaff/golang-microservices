package app

import (
	"context"
	"fmt"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/persistence"
	"github.com/jenpaff/golang-microservices/users"
	"github.com/jenpaff/golang-microservices/validation"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"time"
)

type App struct {
	server     *http.Server
	port       string
	controller *api.Controller
}

func NewApp(port, configPath, secretsPath, secretsEnv string) (*App, error) {
	cfg, err := config.BuildConfig(configPath, secretsPath, secretsEnv)
	if err != nil {
		return nil, err
	}
	db, err := persistence.ConnectPostgres(cfg.Persistence)
	if err != nil {
		log.Errorf("could not establish database connection to %s:%d: %s", cfg.Persistence.DbHost, cfg.Persistence.DbPort, err.Error())
		return nil, err
	}

	userPersistence := users.NewStorage(db)
	userService := users.NewService(userPersistence)

	validator, err := validation.NewValidate()
	if err != nil {
		return nil, err
	}

	controller := api.NewController(cfg, userService, validator)
	router := api.NewRouter(controller)
	server := &http.Server{Addr: ":" + port, Handler: router}
	return &App{server: server, port: port, controller: controller}, nil
}

func (a App) Start() error {
	ctx := context.Background()

	log.Info("Starting...")

	err := ensureDatabaseConnectivity(ctx, a.controller.Cfg.Persistence)
	if err != nil {
		return err
	}

	// everything below enables graceful shutdown of our service without dropping any requests

	go func() {
		log.Infof("Listening on port %s", a.port)
		err = a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Errorf("Could not listen on port %s", a.port)
			os.Exit(1)
		}
	}()

	return nil
}

func (a App) Stop() {
	log.Info("Shutting down server")

	if err := a.server.Shutdown(context.Background()); err != nil {
		log.WithError(fmt.Errorf("couldn't shutdown server cleanly: %s", err.Error()))
	}

	log.Info("Shutting down done")
}

func ensureDatabaseConnectivity(ctx context.Context, cfg config.PersistenceConfig) error {
	deadline := 20 * time.Second
	pollingDelay := 500 * time.Millisecond

	ctx, cancel := context.WithTimeout(ctx, deadline)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	pgOptions := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", cfg.DbHost, cfg.DbPort, cfg.DbUsername, cfg.DbName, cfg.DbPassword)
	if !cfg.SslEnabled {
		pgOptions = pgOptions + " sslmode=disable"
	}

	g.Go(func() error {
		log.Infof("Checking for database connectivity on host: %s port: %d with user: %s", cfg.DbHost, cfg.DbPort, cfg.DbUsername)
		err := persistence.EnsureConnected(ctx, pgOptions, pollingDelay)
		if err != nil {
			return fmt.Errorf("could not initialise database: %s", err.Error())
		}
		log.Infof("Database connection successful")
		return nil
	})

	return g.Wait()
}
