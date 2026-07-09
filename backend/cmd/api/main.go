package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	appathletedefaultsession "github.com/rxwod/backend/internal/application/athletedefaultsession"
	appauth "github.com/rxwod/backend/internal/application/auth"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appclassbooking "github.com/rxwod/backend/internal/application/classbooking"
	appclasssession "github.com/rxwod/backend/internal/application/classsession"
	appgym "github.com/rxwod/backend/internal/application/gym"
	appwod "github.com/rxwod/backend/internal/application/wod"
	appwodresult "github.com/rxwod/backend/internal/application/wodresult"
	deliveryhttp "github.com/rxwod/backend/internal/delivery/http"
	"github.com/rxwod/backend/internal/infrastructure/config"
	"github.com/rxwod/backend/internal/infrastructure/email"
	infrajwt "github.com/rxwod/backend/internal/infrastructure/jwt"
	"github.com/rxwod/backend/internal/infrastructure/password"
	"github.com/rxwod/backend/internal/infrastructure/postgres"
	"github.com/rxwod/backend/internal/platform/clock"
	"github.com/rxwod/backend/internal/platform/idgen"
	"github.com/rxwod/backend/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	appLogger, err := logger.Setup(cfg.LogLevel())
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}
	slog.SetDefault(appLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.NewDB(ctx, cfg.DatabaseURL())
	if err != nil {
		slog.Error("connect database failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	systemClock := clock.System{}
	uuidGenerator := idgen.UUIDGenerator{}

	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	resetTokenRepo := postgres.NewPasswordResetTokenRepository(db)
	verificationTokenRepo := postgres.NewEmailVerificationTokenRepository(db)
	gymRepo := postgres.NewGymRepository(db)
	wodRepo := postgres.NewWODRepository(db)
	wodResultRepo := postgres.NewWODResultRepository(db)
	classSessionRepo := postgres.NewClassSessionRepository(db)
	classBookingRepo := postgres.NewClassBookingRepository(db)
	athleteDefaultSessionRepo := postgres.NewAthleteDefaultSessionRepository(db)

	gymService := appgym.NewService(gymRepo, userRepo, systemClock, uuidGenerator, 7*24*time.Hour)
	accessTokenIssuer := infrajwt.NewAccessTokenIssuer(cfg.JWTSecret(), cfg.AccessTokenTTL())
	emailSender := email.NewLogSender()
	authService := appauth.NewService(
		userRepo,
		refreshTokenRepo,
		resetTokenRepo,
		verificationTokenRepo,
		password.NewBcryptHasher(),
		accessTokenIssuer,
		gymService,
		emailSender,
		systemClock,
		uuidGenerator,
		cfg.RefreshTokenTTL(),
		cfg.FrontendURL(),
		cfg.PasswordResetTTL(),
		cfg.EmailVerificationTTL(),
	)
	authorizer := appauthz.NewAuthorizer(gymRepo)
	wodService := appwod.NewService(wodRepo, systemClock, uuidGenerator)
	wodResultService := appwodresult.NewService(wodResultRepo, wodRepo, gymRepo, systemClock, uuidGenerator)
	classSessionService := appclasssession.NewService(classSessionRepo, classBookingRepo, athleteDefaultSessionRepo, gymRepo, systemClock, uuidGenerator)
	classBookingService := appclassbooking.NewService(classBookingRepo, classSessionRepo, gymRepo, systemClock, uuidGenerator)
	athleteDefaultSessionService := appathletedefaultsession.NewService(athleteDefaultSessionRepo, gymRepo, systemClock, uuidGenerator)
	router := deliveryhttp.NewRouter(
		authService,
		gymService,
		wodService,
		wodResultService,
		classSessionService,
		classBookingService,
		athleteDefaultSessionService,
		authorizer,
		db,
		cfg.AllowedOrigins(),
		cfg,
	)

	go func() {
		address := fmt.Sprintf(":%d", cfg.HTTPPort())
		slog.Info("starting api", "event", "app.starting", "address", address, "env", cfg.AppEnv())
		if err := router.Start(address); err != nil {
			slog.Error("server stopped unexpectedly", "event", "app.server_stopped", "error", err)
			cancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := router.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "event", "app.shutdown_error", "error", err)
	}
}
