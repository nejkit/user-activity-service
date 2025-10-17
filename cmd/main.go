package main

import (
	"context"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-activity-service/config"
	"user-activity-service/handlers"
	"user-activity-service/server"
	"user-activity-service/services"
	"user-activity-service/storage"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL

	cfg := config.GetConfig()

	customizeLogger(cfg.LoggerLevel)
	db, err := sqlx.Connect("postgres", cfg.DatabaseConfig.ToConnectionString())

	if err != nil {
		os.Exit(1)
	}

	usersStorage := storage.NewPgUserStorage(db)
	activitiesStorage := storage.NewPgUserActivityStorage(db)
	activityHistoryStorage := storage.NewPgUserActivityHistoryStorage(db)
	activityPeriodsStorage := storage.NewPgActivityPeriodsStorage(db)

	usersService := services.NewUserService(usersStorage)
	activityService := services.NewActivityService(usersStorage, activitiesStorage, activityHistoryStorage)
	backgroundService := services.NewBackgroundActivityService(activityPeriodsStorage, activityHistoryStorage, activitiesStorage, &cfg.BackgroundWorkerConfig)

	usersHandler := handlers.NewUserHandler(usersService)
	activityHandler := handlers.NewActivityHandler(activityService)

	go backgroundService.Run(ctx)

	engine := handlers.InitEngine(usersHandler, activityHandler, cfg.CorsOrigins)

	httpServer := server.NewHttpServer(engine, cfg.ApplicationPort)
	shutdownChan := make(chan struct{})

	go func() {
		if err := httpServer.Run(shutdownChan); err != nil {
			fmt.Println(err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	handleShutdown(ctx, backgroundService.Wg, usersHandler, activityHandler, shutdownChan, httpServer.ServerStoppedChan)
}

func customizeLogger(level string) {
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	parsedLevel, err := logger.ParseLevel(level)

	if err != nil {
		logger.Warningln("Error parsing log level from config. Using default: INFO")
		parsedLevel = logger.InfoLevel
	}

	logger.SetLevel(parsedLevel)
}

func handleShutdown(ctx context.Context, backgroundServiceWg *sync.WaitGroup, usersHandler *handlers.UserHandler, activityHandler *handlers.ActivityHandler, shutdownChan chan<- struct{}, completeShutdownChan <-chan struct{}) {
	logger.Infoln("run handler for stop application")
	<-ctx.Done()

	logger.Infoln("start shutdown application")
	logger.Infoln("mark api as stopped")

	activityHandler.ApiStopped = true
	usersHandler.ApiStopped = true

	logger.Infoln("wait stop scheduler of activity tasks")
	backgroundServiceWg.Wait()
	logger.Infoln("wait stop user controller")
	usersHandler.Wg.Wait()
	logger.Infoln("wait stop activity controller")
	activityHandler.Wg.Wait()

	shutdownChan <- struct{}{}
	close(shutdownChan)

	<-completeShutdownChan

	logger.Infoln("application stopped successfully")
}
