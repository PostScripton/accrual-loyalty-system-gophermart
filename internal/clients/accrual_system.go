package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Order struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual"`
}

type AccrualSystemClient struct {
	client  *resty.Client
	retryAt *time.Time
}

var ErrAccrualSystemIsUnavailable = errors.New("accrual system is unavailable")

func NewAccrualSystemClient(baseURL string) *AccrualSystemClient {
	return &AccrualSystemClient{
		client: resty.New().SetBaseURL(baseURL),
	}
}

func (asc *AccrualSystemClient) GetOrderInfo(ctx context.Context, number string) (*Order, error) {
	order := new(Order)

	resp, err := asc.client.R().
		SetContext(ctx).
		SetResult(order).
		Get(fmt.Sprintf("/api/orders/%s", number))
	if err != nil {
		return nil, ErrAccrualSystemIsUnavailable
	}
	if err = asc.blocked(resp); err != nil {
		return nil, err
	}

	log.Debug().
		Str("status", resp.Status()).
		Str("msg", string(resp.Body())).
		Str("order", number).
		Msg("Polling order status")

	if resp.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	return order, nil
}

func (asc *AccrualSystemClient) CanMakeRequest() error {
	if asc.retryAt == nil {
		return nil
	}

	if time.Now().After(*asc.retryAt) {
		asc.retryAt = nil
		return nil
	}

	return fmt.Errorf("client unblocks in %s", time.Until(*asc.retryAt))
}

func (asc *AccrualSystemClient) blocked(resp *resty.Response) error {
	if !resp.IsError() && resp.StatusCode() != http.StatusTooManyRequests {
		return nil
	}

	retryAfter, err := time.ParseDuration(fmt.Sprintf("%ss", resp.Header().Get("Retry-After")))
	if err != nil {
		return err
	}

	body := string(resp.Body())
	log.Debug().Str("resp", body).Msg("Client got blocked")
	*asc.retryAt = time.Now().Add(retryAfter)

	return fmt.Errorf(body)
}
