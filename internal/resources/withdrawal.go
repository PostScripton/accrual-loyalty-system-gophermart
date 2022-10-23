package resources

import (
	"github.com/PostScripton/accrual-loyalty-system-gophermart/internal/models"
	"time"
)

type WithdrawalResource struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func MakeWithdrawalResource(withdrawal *models.Withdrawal) WithdrawalResource {
	return WithdrawalResource{
		Order:       withdrawal.Order,
		Sum:         withdrawal.Sum,
		ProcessedAt: withdrawal.CreatedAt.Format(time.RFC3339),
	}
}

func MakeWithdrawalResourceCollection(collection []*models.Withdrawal) (result []WithdrawalResource) {
	for _, withdrawal := range collection {
		result = append(result, MakeWithdrawalResource(withdrawal))
	}
	return result
}
