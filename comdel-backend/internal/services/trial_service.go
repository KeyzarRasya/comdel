package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

/*
	Restricted Function Area
	Only On When Trial Test time
*/

 func GrantSubscriptionAccess(tx pgx.Tx, userId string) error {

	var subscriptionExpiry time.Time = time.Now().Add((time.Hour * 24) * 30);

	_, err := tx.Exec(
		context.Background(),
		"UPDATE user_info SET premium_plan = 'NEWBIE', subscription = 'ACTIVE', subscription_expiry = $1 WHERE user_id=$2",
		subscriptionExpiry, userId,
	)

	return err

 } 