package services

import (
	"context"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/clients"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type OrderService struct {
	repo   repository.Orders
	client *clients.AccrualSystemClient
}

func NewOrderService(repo repository.Orders, client *clients.AccrualSystemClient) *OrderService {
	return &OrderService{
		repo:   repo,
		client: client,
	}
}

func (os *OrderService) Create(ctx context.Context, number string, user *models.User) (*models.Order, error) {
	if err := os.repo.Create(ctx, number, user); err != nil {
		return nil, err
	}

	return os.FindByNumber(ctx, number, user)
}

func (os *OrderService) FindByNumber(ctx context.Context, number string, user *models.User) (*models.Order, error) {
	return os.repo.FindByNumber(ctx, number, user)
}

func (os *OrderService) All(ctx context.Context, user *models.User) ([]*models.Order, error) {
	return os.repo.All(ctx, user)
}

func (os *OrderService) RunPollingStatuses(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
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
		resp, err := os.client.GetOrderInfo(ctx, order.Number)
		if err != nil {
			return err
		}
		if resp == nil {
			continue
		}

		if order.Status == resp.Status {
			continue
		}

		order.Status = resp.Status
		order.Accrual = resp.Accrual

		if err := os.repo.Update(ctx, order); err != nil {
			return err
		}

		log.Info().Int("order_id", order.ID).Msg("Order has been updated")
	}

	return nil
}
