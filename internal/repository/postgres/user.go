package postgres

import (
	"context"
	"errors"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/jackc/pgx/v4"
)

type UserRepository struct {
	db *Postgres
}

func NewUserRepository(db *Postgres) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, login, hashedPassword string) error {
	sql := `INSERT INTO users (login, password) VALUES ($1, $2);`

	if _, err := ur.db.Exec(ctx, sql, login, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	sql := `SELECT id, login, password, created_at, updated_at FROM users WHERE login = $1;`

	user := new(models.User)

	row := ur.db.QueryRow(ctx, sql, login)
	if err := row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
