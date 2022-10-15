package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, address string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, address)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		pool.Close()
	}()

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
