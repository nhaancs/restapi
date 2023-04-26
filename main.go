package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"restapi/module/user/userbusiness"
	"restapi/module/user/userstore"
	"restapi/module/user/usertransport"
	"restapi/pkg/config"
	"restapi/pkg/database"
	"restapi/pkg/logger"
	"restapi/pkg/md5hasher"
	"restapi/pkg/tokenprovider/jwtprovider"
	"runtime"
	"syscall"
	"time"
)

type appConfig struct {
	Auth struct {
		Secret string
	}
	DB struct {
		Addr string
		Name string
		User string
		Pass string
	}
	Redis struct {
		Addr string
		Pass string
	}
}

func main() {
	log, err := logger.New("RESTAPI")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Fatal("failed to start service: ", err)
	}
}

func run(log *zap.SugaredLogger) error {
	// Set the correct number of threads for the service based on what is available either by the machine or quotas.
	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Info("GOMAXPROCS: ", runtime.GOMAXPROCS(0))

	var cfg appConfig
	if err := config.Load("config.yaml", &cfg); err != nil {
		fmt.Errorf("error loading config: %+v", err)
	}
	log.Infof("loaded config: %+v", cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1 := r.Group("v1")

	db, err := database.Connect(database.Config{
		Addr: cfg.DB.Addr,
		Name: cfg.DB.Name,
		User: cfg.DB.User,
		Pass: cfg.DB.Pass,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %+v", err)
	}
	defer db.Close()

	userBusiness := userbusiness.New(
		log,
		userstore.New(db),
		jwtprovider.New(cfg.Auth.Secret),
		md5hasher.New(),
	)
	userTransport := usertransport.New(
		log,
		userBusiness,
	)
	userTransport.SetupRoutes(v1)

	// Construct a server to service the requests.
	app := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infof("service started at %s", app.Addr)
		serverErrors <- app.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %+v", err)
	case sig := <-shutdown:
		log.Info("shutdown started with signal:", sig)
		defer log.Info("shutdown completed")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := app.Shutdown(ctx); err != nil {
			err := app.Close()
			if err != nil {
				return fmt.Errorf("app.Close error: %+v", err)
			}
			return fmt.Errorf("app.Shutdown error: %+v", err)
		}
	}

	return nil
}
