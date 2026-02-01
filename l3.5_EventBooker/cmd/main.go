package main

import (
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/config"
	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/handler"
	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/service"
	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/storage"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("Starting application...")

	var cfg config.Config
	cfgLoader := config.New()
	if err := cfgLoader.Load("config/config.yaml"); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("Failed to load config file")
	}

	if err := cfgLoader.Unmarshal(&cfg); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("Failed to unmarshal config")
	}

	db := storage.InitDB(cfg.DB)
	defer storage.CloseDB(db)

	store := storage.New(db)
	srv := service.New(store)

	go startExpiredBookingWorker(srv)

	h := handler.New(srv)
	zlog.Logger.Info().Str("addr", cfg.Server.Address).Msg("starting server")
	if err := h.Router().Run(cfg.Server.Address); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("server failed to start")
	}
}

func startExpiredBookingWorker(srv *service.Service) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := cancelExpiredBookings(srv); err != nil {
			zlog.Logger.Error().Err(err).Msg("failed to cancel expired bookings")
		}
	}
}

func cancelExpiredBookings(srv *service.Service) error {
	return srv.CancelExpiredBookings()
}
