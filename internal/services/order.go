package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/clients"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type OrderService struct {
	repo   repository.Orders
	users  repository.Users
	client *clients.AccrualSystemClient
}

func NewOrderService(repo repository.Orders, users repository.Users, client *clients.AccrualSystemClient) *OrderService {
	return &OrderService{
		repo:   repo,
		users:  users,
		client: client,
	}
}

func (os *OrderService) Create(ctx context.Context, number string, user *models.User) (*models.Order, error) {
	if err := os.repo.Create(ctx, number, user); err != nil {
		return nil, err
	}

	order, err := os.FindByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	if _, err = os.PollStatus(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (os *OrderService) FindByNumber(ctx context.Context, number string) (*models.Order, error) {
	return os.repo.FindByNumber(ctx, number)
}

func (os *OrderService) All(ctx context.Context, user *models.User) ([]*models.Order, error) {
	return os.repo.All(ctx, user)
}

func (os *OrderService) RunPollingStatuses(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := os.pollStatuses(ctx); err != nil {
				return err
			}
		}
	}
}

func (os *OrderService) PollStatus(ctx context.Context, order *models.Order) (bool, error) {
	resp, err := os.client.GetOrderInfo(ctx, order.Number)
	if err != nil {
		return false, fmt.Errorf("get order info: %w", err)
	}
	if resp == nil {
		return false, nil
	}

	if order.Status == resp.Status {
		return false, nil
	}

	order.Status = resp.Status
	order.Accrual = resp.Accrual

	if err = os.repo.Update(ctx, order); err != nil {
		return false, err
	}

	user, err := os.users.Find(ctx, order.UserID)
	if err != nil {
		return false, err
	}
	user.Balance += *order.Accrual

	if err = os.users.Update(ctx, user); err != nil {
		return false, err
	}

	return true, nil
}

func (os *OrderService) pollStatuses(ctx context.Context) error {
	if err := os.client.CanMakeRequest(); err != nil {
		log.Debug().Err(err).Msg("Client cannot make a request")
		return nil
	}

	collection, err := os.repo.AllPending(ctx)
	if err != nil {
		return err
	}

	for _, order := range collection {
		updated, err := os.PollStatus(ctx, order)
		if err != nil {
			return err
		}
		if updated {
			log.Info().Int("order_id", order.ID).Msg("Order has been updated")
		}
	}

	return nil
}
