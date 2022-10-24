package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/jackc/pgx/v4"
	"time"
)

type OrderRepository struct {
	db *Postgres
}

func NewOrderRepository(db *Postgres) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (or *OrderRepository) Create(ctx context.Context, number string, user *models.User) error {
	sql := `INSERT INTO orders (user_id, number, status) VALUES ($1, $2, $3);`

	if _, err := or.db.Exec(ctx, sql, user.ID, number, models.NewOrderStatus); err != nil {
		return err
	}

	return nil
}

func (or *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	sql := `UPDATE orders SET status = $1, accrual = $2, updated_at = $3 WHERE number = $4;`

	if _, err := or.db.Exec(ctx, sql, order.Status, order.Accrual, time.Now(), order.Number); err != nil {
		return err
	}

	return nil
}

func (or *OrderRepository) FindByNumber(ctx context.Context, number string) (*models.Order, error) {
	sql := `SELECT id, user_id, number, status, accrual, created_at, updated_at FROM orders WHERE number = $1;`

	order := new(models.Order)

	row := or.db.QueryRow(ctx, sql, number)
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return order, nil
}

func (or *OrderRepository) All(ctx context.Context, user *models.User) ([]*models.Order, error) {
	count := 0
	if err := or.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, user.ID).Scan(&count); err != nil {
		return []*models.Order{}, err
	}
	if count == 0 {
		return []*models.Order{}, nil
	}

	sql := `SELECT id, user_id, number, status, accrual, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC;`

	rows, err := or.db.Query(ctx, sql, user.ID)
	if err != nil {
		return []*models.Order{}, err
	}

	collection, err := or.scanOrdersCollection(rows, count)
	if err != nil {
		return []*models.Order{}, err
	}

	return collection, nil
}

func (or *OrderRepository) AllPending(ctx context.Context) ([]*models.Order, error) {
	rawSQL := `SELECT %s FROM orders WHERE status IN ($1, $2);`
	sql := fmt.Sprintf(rawSQL, "id, user_id, number, status, accrual, created_at, updated_at")
	sqlCount := fmt.Sprintf(rawSQL, "COUNT(*)")

	count := 0
	if err := or.db.QueryRow(ctx, sqlCount, models.NewOrderStatus, models.ProcessingOrderStatus).Scan(&count); err != nil {
		return []*models.Order{}, err
	}
	if count == 0 {
		return []*models.Order{}, nil
	}

	rows, err := or.db.Query(ctx, sql, models.NewOrderStatus, models.ProcessingOrderStatus)
	if err != nil {
		return []*models.Order{}, err
	}

	collection, err := or.scanOrdersCollection(rows, count)
	if err != nil {
		return []*models.Order{}, nil
	}

	return collection, nil
}

func (or *OrderRepository) scanOrder(row pgx.Row, order *models.Order) error {
	if err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (or *OrderRepository) scanOrdersCollection(rows pgx.Rows, count int) ([]*models.Order, error) {
	collection := make([]*models.Order, 0, count)

	for rows.Next() {
		order := new(models.Order)

		if err := or.scanOrder(rows, order); err != nil {
			return []*models.Order{}, err
		}

		collection = append(collection, order)
	}

	return collection, nil
}
