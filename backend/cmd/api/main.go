package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	deliveryhttp "github.com/rxwod/backend/internal/delivery/http"
	"github.com/rxwod/backend/internal/infrastructure/config"
	infrajwt "github.com/rxwod/backend/internal/infrastructure/jwt"
	"github.com/rxwod/backend/internal/infrastructure/email"
	"github.com/rxwod/backend/internal/infrastructure/password"
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

	systemClock := clock.System{}
	uuidGenerator := idgen.UUIDGenerator{}

	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	resetTokenRepo := postgres.NewPasswordResetTokenRepository(db)
	gymRepo := postgres.NewGymRepository(db)
	wodRepo := postgres.NewWODRepository(db)

	gymService := appgym.NewService(gymRepo, systemClock, uuidGenerator, 7*24*time.Hour)
	accessTokenIssuer := infrajwt.NewAccessTokenIssuer(cfg.JWTSecret(), cfg.AccessTokenTTL())
	emailSender := email.NewLogSender()
	authService := appauth.NewService(
		userRepo,
		refreshTokenRepo,
		resetTokenRepo,
		password.NewBcryptHasher(),
		accessTokenIssuer,
		gymService,
		emailSender,
		systemClock,
		uuidGenerator,
		cfg.RefreshTokenTTL(),
		cfg.FrontendURL(),
		cfg.PasswordResetTTL(),
	)
	authorizer := appauthz.NewAuthorizer(gymRepo)
	wodService := appwod.NewService(wodRepo, systemClock, uuidGenerator)
	router := deliveryhttp.NewRouter(authService, gymService, wodService, authorizer)

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
