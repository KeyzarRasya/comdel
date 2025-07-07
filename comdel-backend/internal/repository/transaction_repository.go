package repository

import "github.com/KeyzarRasya/comdel-server/internal/dto"

type TransactionRepository interface {
	Save(transaction dto.Transaction) 	error
	GetPremiumPlan(orderId string)		(string, error)
}