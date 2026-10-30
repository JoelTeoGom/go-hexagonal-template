package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JoelTeoGom/login-service/internal/adapters/in/httpapi"
	"github.com/JoelTeoGom/login-service/internal/adapters/out/mysql"
	"github.com/JoelTeoGom/login-service/internal/adapters/out/password"
	"github.com/JoelTeoGom/login-service/internal/adapters/out/token"
	"github.com/JoelTeoGom/login-service/internal/config"
	"github.com/JoelTeoGom/login-service/internal/service"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	appConfig, err := config.Load()
	if err != nil {
		log.Error("load config failed", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := mysql.Open(ctx, appConfig.MySQLDSN)
	if err != nil {
		log.Error("open database failed", "err", err)
		os.Exit(1)
	}

	// Outbound adapters
	users := mysql.NewUserRepository(db)
	hasher := password.NewBcrypt(0)
	tokens := token.NewJWT(appConfig.JWTSecret, appConfig.JWTTTL, "login-service")

	// Use cases
	authService := service.NewAuthService(users, hasher, tokens)
	userService := service.NewUserService(users)

	// Inbound adapter
	handler := httpapi.NewHandler(authService, userService, log)
	server := &http.Server{
		Addr:              appConfig.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", appConfig.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// From here on the database is open and os.Exit skips defers, so every
	// exit path closes it explicitly first.
	select {
	case <-ctx.Done():
		log.Info("shutting down")
	case err := <-serverErrors:
		log.Error("server failed", "err", err)
		if err := db.Close(); err != nil {
			log.Error("close database failed", "err", err)
		}
		os.Exit(1)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown failed", "err", err)
		if err := db.Close(); err != nil {
			log.Error("close database failed", "err", err)
		}
		os.Exit(1)
	}

	if err := db.Close(); err != nil {
		log.Error("close database failed", "err", err)
		os.Exit(1)
	}
}
