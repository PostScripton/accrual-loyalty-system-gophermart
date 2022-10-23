package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (ur *UserRepository) Update(ctx context.Context, user *models.User) error {
	sql := `UPDATE users SET balance = $1 WHERE id = $2;`

	if _, err := ur.db.Exec(ctx, sql, user.Balance, user.ID); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Find(ctx context.Context, id int) (*models.User, error) {
	return ur.find(ctx, "id", id)
}

func (ur *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	return ur.find(ctx, "login", login)
}

func (ur *UserRepository) find(ctx context.Context, column string, value interface{}) (*models.User, error) {
	sql := `SELECT id, login, password, balance, created_at, updated_at FROM users WHERE %s = $1;`

	user := new(models.User)

	row := ur.db.QueryRow(ctx, fmt.Sprintf(sql, column), value)
	if err := ur.scanUser(row, user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) scanUser(row pgx.Row, user *models.User) error {
	if err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}

	return nil
}
