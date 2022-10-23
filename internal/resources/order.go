package resources

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"time"
)

type OrderResource struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

func MakeOrderResource(order *models.Order) OrderResource {
	return OrderResource{
		Number:     order.Number,
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: order.CreatedAt.Format(time.RFC3339),
	}
}

func MakeOrderResourceCollection(collection []*models.Order) (result []OrderResource) {
	for _, order := range collection {
		result = append(result, MakeOrderResource(order))
	}
	return result
}
