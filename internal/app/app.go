package app

import (
	"context"

	grpcapp "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/app/grpc"
	httpapp "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/app/http"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/config"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/postgres"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/redis"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/auth"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/user"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/service/verification"
	userstore "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/postgres/user"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/redis/refreshtoken"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/redis/verificationcode"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpapp.App
	PG         *postgres.Postgres
	RDB        *redis.Redis
}

func New(ctx context.Context, cfg *config.Config) *App {
	// Init logger
	logger.New(cfg.Log.Level)

	// Postgres connection
	pg, err := postgres.New(ctx, cfg.DB.URL)
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to database: %s", err.Error())
	}

	// Redis connection
	rdb, err := redis.New(ctx, redis.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		logger.Log().Fatal(ctx, "error with connection to redis: %s", err.Error())
	}

	// Auth config
	authConfig := auth.NewConfig(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	// Store
	userStore := userstore.New(pg)
	refreshTokenStore := refreshtoken.New(rdb)
	verificationCodeStore := verificationcode.New(rdb)

	// Service
	authService := auth.New(userStore, refreshTokenStore, authConfig, verificationCodeStore)
	userService := user.New(userStore, verificationCodeStore)

	verificationService := verification.New(verificationCodeStore, userStore)
	verificationService.RegisterEmailService(cfg.SMPT.Host, cfg.SMPT.Port, cfg.SMPT.Username, cfg.SMPT.Password, cfg.SMPT.Sender)
	verificationService.RegisterSMSService(cfg.SMS.Sender)

	// gRPC server
	gRPCApp := grpcapp.New(
		ctx,
		cfg,
		authService,
		userService,
		authConfig,
		verificationService,
	)

	// HTTP server
	httpServer := httpapp.New(ctx, cfg)

	return &App{
		GRPCServer: gRPCApp,
		HTTPServer: httpServer,
		PG:         pg,
		RDB:        rdb,
	}
}
