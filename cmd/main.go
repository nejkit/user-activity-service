package main

import (
	"context"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	cfg := config.Config{}

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

	go handleApiStoppedSignal(ctx, usersHandler, activityHandler)
	go backgroundService.Run(ctx)

	engine := handlers.InitEngine(usersHandler, activityHandler)

	httpServer := server.NewHttpServer(engine, cfg.ApplicationPort)
	shutdownChan := make(chan struct{})

	go func() {
		if err := httpServer.Run(shutdownChan); err != nil {
			cancel()
		}
	}()

	handleShutdown(ctx, backgroundService.Wg, usersHandler.Wg, activityHandler.Wg, shutdownChan)
}

func handleApiStoppedSignal(ctx context.Context, usersHandler *handlers.UserHandler, activityHandler *handlers.ActivityHandler) {
	<-ctx.Done()

	usersHandler.ApiStopped = true
	activityHandler.ApiStopped = true
}

func handleShutdown(ctx context.Context, backgroundServiceWg, usersHandlerWg *sync.WaitGroup, activityHandlerWg *sync.WaitGroup, shutdownChan chan<- struct{}) {
	<-ctx.Done()
	backgroundServiceWg.Wait()
	usersHandlerWg.Wait()
	activityHandlerWg.Wait()

	shutdownChan <- struct{}{}
	close(shutdownChan)
}
