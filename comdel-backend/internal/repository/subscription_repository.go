package repository

import (
	"context"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository interface {
	SaveReturningSubsId(tx pgx.Tx, userId string, subscription *model.Subscription, premiumPlan string)		(string, error);
	Activate(tx pgx.Tx, plan string, subsId string, userId string)							error;
	GetExpiryTimeBySubsId(subsId string)													(time.Time, error)
}

type SubscriptionRepositoryImpl struct {
	conn *pgx.Conn;
}

func NewSubscriptionRepository(pgxConn *pgx.Conn) SubscriptionRepository {
	return &SubscriptionRepositoryImpl{conn: pgxConn}
}

func (sr *SubscriptionRepositoryImpl) SaveReturningSubsId(tx pgx.Tx, userId string, subscription *model.Subscription, premiumPlan string) (string, error) {
	var subsId string;
	subscription.End = time.Now().Add((time.Hour * 24) * 30);
		
	err := tx.QueryRow(
		context.Background(),
		"INSERT INTO subscription(user_id, bank, transaction_time, payment_type, fraud_status, status_code, settlement_time, expiry_time, transaction_status, premium_plan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING subs_id",
		userId, subscription.Bank, subscription.TransactionTime.Time(), subscription.PaymentType, subscription.FraudStatus, subscription.StatusCode, subscription.SettlementTime.Time(), subscription.End, subscription.TransactionStatus, premiumPlan,
	).Scan(&subsId)

	if err != nil {
		return "", err;
	}

	subscription.ExpiryTime = model.MidtransTime(subscription.End);
	subscription.PremiumPlan = premiumPlan;

	return subsId, nil
}

func (sr *SubscriptionRepositoryImpl) Activate(tx pgx.Tx, plan string, subsId string, userId string) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE user_info SET subscription = 'ACTIVE', premium_plan = $1, subs_id = $2 WHERE user_id=$3",
		plan, subsId, userId,
	)

	return err;
}

func (sr *SubscriptionRepositoryImpl) GetExpiryTimeBySubsId(subsId string) (time.Time, error) {
	var expiryTime time.Time;

	err := sr.conn.QueryRow(
		context.Background(),
		"SELECT expiry_time FROM subscription WHERE subs_id=$1",
		subsId,
	).Scan(&expiryTime);

	if err != nil {
		return time.Time{}, err
	}

	return expiryTime, nil
}