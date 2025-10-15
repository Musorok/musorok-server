package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/joho/godotenv"

	"github.com/musorok/server/config"
	"github.com/musorok/server/internal/app"
	httpapi "github.com/musorok/server/internal/http"
	"github.com/musorok/server/internal/repo/postgres"
	redisrepo "github.com/musorok/server/internal/repo/redis"
	"github.com/musorok/server/internal/core/payments/paynetworks"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil { panic(err) }
	app.SetupLogger()

	db, err := postgres.Open(cfg.DBDSN)
	if err != nil { log.Fatal().Err(err).Msg("open db") }

	r := redisrepo.Open(cfg.RedisAddr)
	if err := redisrepo.Ping(context.Background(), r); err != nil { log.Fatal().Err(err).Msg("redis ping") }

	pay := paynetworks.New(cfg.PayAPIKey, cfg.PayReturnURL)
	router := httpapi.NewRouter(db, cfg.JWTSecret, cfg.JWTRefreshSecret, int64(cfg.JWTAccessTTL.Seconds()), int64(cfg.JWTRefreshTTL.Seconds()), pay)
	srv := &http.Server{ Addr: ":"+cfg.AppPort, Handler: router }

	application := &app.App{ Server: srv }
	go func(){ if err := application.Start(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("server") } }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = application.Shutdown(ctx)
}
