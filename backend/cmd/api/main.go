package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	appwod "github.com/rxwod/backend/internal/application/wod"
	deliveryhttp "github.com/rxwod/backend/internal/delivery/http"
	"github.com/rxwod/backend/internal/infrastructure/config"
	"github.com/rxwod/backend/internal/infrastructure/postgres"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.NewDB(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewWODRepository(db)
	service := appwod.NewService(repo, clock.System{}, idgen.UUIDGenerator{})
	router := deliveryhttp.NewRouter(service)

	go func() {
		address := fmt.Sprintf(":%d", cfg.HTTPPort())
		log.Printf("starting api on %s", address)
		if err := router.Start(address); err != nil {
			log.Printf("server stopped: %v", err)
			cancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := router.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
