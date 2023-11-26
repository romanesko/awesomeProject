package postgres

import (
	"context"
	"log"
	"server/config"
	"server/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPoolConnection(ctx context.Context, cfg *config.Config) (pool *pgxpool.Pool) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.GetPostgresConnString())
	if err != nil {
		log.Fatalf("unable to parse connection string, error: %v", err.Error())
	}

	pgxConfig.MaxConns = 10

	err = utils.ConnectionAttemps(func() error {
		var err error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		pool, err = pgxpool.NewWithConfig(ctx, pgxConfig)

		return err
	}, 3, time.Duration(2)*time.Second)

	if err != nil {
		log.Fatalf("didn't manage to make connection with database, error: %v", err.Error())
	}

	log.Println("database connection is established successfully")

	return pool
}
