package model

import "time"

type Subscription struct {
	Bank				string			`json:"bank"`
	TransactionTime 	MidtransTime	`json:"transaction_time"`
	OrderID				string			`json:"order_id"`
	PaymentType			string			`json:"payment_type"`
	FraudStatus			string			`json:"fraud_status"`
	StatusCode			string			`json:"status_code"`
	SettlementTime		MidtransTime	`json:"settlement_time"`
	TransactionStatus	string			`json:"transaction_status"`
	ExpiryTime			MidtransTime	`json:"expiry_time"`
	PremiumPlan			string			`json:"premium_plan"`
	End					time.Time		`json:"subscriptionEnd"`
}