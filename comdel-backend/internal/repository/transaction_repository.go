package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/model"
)

type TransactionRepository interface {
	Create(tx config.DBTx, transaction dto.Transaction) 						error
	GetPremiumPlan(orderId string)										(string, error)
	UpdateTransactionStatus(tx config.DBTx, status string, orderId string)	error
	Status(orderid string)												(*model.Subscription, error)
}

type TransactionRepositoryImpl struct {
	conn config.DBConn;
}

func NewTransactionRepository(pgxConn config.DBConn) TransactionRepository {
	return &TransactionRepositoryImpl{conn: pgxConn}
}

func (tr *TransactionRepositoryImpl) Create(tx config.DBTx, transaction dto.Transaction) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO transaction(user_id, order_id, premium_plan) VALUES($1, $2, $3)",
		transaction.UserId, transaction.OrderId, transaction.Plan,
	)

	return err;
}

func (tr *TransactionRepositoryImpl) GetPremiumPlan(orderId string)	(string, error) {
	var premiumPlan string;
	err := tr.conn.QueryRow(
		context.Background(),
		"SELECT premium_plan FROM transaction WHERE order_id=$1",
		orderId,
	).Scan(&premiumPlan)

	if err != nil {
		return "", err;
	}

	return premiumPlan, nil
}

func (tr *TransactionRepositoryImpl) UpdateTransactionStatus(tx config.DBTx, status string, orderId string) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE transaction SET transaction_status=$1 WHERE order_id=$2",
		status, orderId,
	)

	return err;
}
 
func (tr *TransactionRepositoryImpl) Status(orderId string) (*model.Subscription, error) {
	var statusResponse model.Subscription;
	var client http.Client;
	var endpoint = fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/status", orderId);

	var serverKey = fmt.Sprintf("%s:", os.Getenv("MIDTRANS_SERVER_KEY"))
	var encodedServerKey = base64.StdEncoding.EncodeToString([]byte(serverKey))

	req, err := http.NewRequest("GET", endpoint, nil);

	if err != nil {
		return nil, err;
	}

	req.Header.Add("Accept", "application/json");
	req.Header.Add("Content-Type", "application/json");
	req.Header.Add("Authorization", encodedServerKey);

	resp, err := client.Do(req)

	if err != nil {
		return nil, err;
	}

	byteBody, err := io.ReadAll(resp.Body);

	if err != nil {
		return nil, err;
	}

	if err := json.Unmarshal(byteBody, &statusResponse); err != nil {
		return nil, err;
	}

	return &statusResponse, nil;
}