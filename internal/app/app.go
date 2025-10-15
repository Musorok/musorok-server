package app

import (
	"context"
	"net/http"
	"time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct { Server *http.Server }

func (a *App) Start() error {
	log.Info().Str("addr", a.Server.Addr).Msg("server starting")
	return a.Server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("server shutting down")
	return a.Server.Shutdown(ctx)
}

func SetupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.With().Str("app", "musorok").Logger()
}
