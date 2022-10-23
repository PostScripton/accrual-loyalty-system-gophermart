package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres(ctx context.Context, address string) (*Postgres, error) {
	pool, err := pgxpool.Connect(ctx, address)
	if err != nil {
		return nil, err
	}

	db := &Postgres{
		pool,
	}

	go func() {
		<-ctx.Done()
		db.Close()
	}()

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = db.MigrateAllMigrations(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
