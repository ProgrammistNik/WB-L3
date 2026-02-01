package storage

import (
	"fmt"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/config"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

func InitDB(cfg config.DB) *dbpg.DB {
	masterDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name,
	)

	opts := &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
	}

	db, err := dbpg.New(masterDSN, nil, opts)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("Failed to initialize DB connection")
	}

	if err := db.Master.Ping(); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("Failed to ping master DB")
	}

	zlog.Logger.Info().Msg("Database connected successfully (master)")

	return db
}

func CloseDB(db *dbpg.DB) {
	if db == nil {
		return
	}

	if err := db.Master.Close(); err != nil {
		zlog.Logger.Error().
			Err(err).
			Msg("Error closing master DB")
	} else {
		zlog.Logger.Info().Msg("Master DB connection closed")
	}

	for i, slave := range db.Slaves {
		if err := slave.Close(); err != nil {
			zlog.Logger.Error().
				Err(err).
				Int("index", i).
				Msg("Error closing slave DB")
		} else {
			zlog.Logger.Info().
				Int("index", i).
				Msg("Slave DB connection closed")
		}
	}
}
