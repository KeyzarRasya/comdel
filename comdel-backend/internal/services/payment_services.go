package services

import (
	"context"
	"os"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func Pay(cookie string, plan dto.PremiumPlan) dto.Response {
	var client snap.Client;
	var userName string;

	conn := config.LoadDatabase();

	userId, err := helper.VerifyAndGet(cookie);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT name FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&userName)

	tx, err := conn.Begin(context.Background())

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		}
	}


	defer tx.Rollback(context.Background());


	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	var orderId string = uuid.NewString();
	var price int64;
	var planStr string;

	if plan == dto.Creator {
		planStr = "CREATOR";
		price = 20000;
	} else {
		planStr = "NEWBIE";
		price = 15000;
	}

	req := & snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: orderId,
			GrossAmt: price,
		},

		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: userName,
		},

		
	}

	resp, _ := client.CreateTransaction(req)


	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO transaction(user_id, order_id, premium_plan) VALUES($1, $2, $3)",
		userId, orderId, planStr,
	)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		}
	}

	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Creating a transaction",
		Data: resp.RedirectURL,
	}
}

func FinishPayment(cookie string, transaction dto.TransactionStatus) dto.Response {
	conn := config.LoadDatabase()
	var premiumPlan string;
	var subsId string;

	userId, err := helper.VerifyAndGet(cookie)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get your information",
			Data: err.Error(),
		}
	}


	statusResponse, err := dto.Status(transaction.OrderID);
	
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message:"Failed to parse into status response",
			Data: err.Error(),
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT premium_plan FROM transaction WHERE order_id=$1",
		statusResponse.OrderID,
	).Scan(&premiumPlan)

	tx, err := conn.Begin(context.Background());

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get plan information",
			Data: err.Error(),
		}
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"UPDATE transaction SET transaction_status=$1 WHERE order_id=$2",
		statusResponse.TransactionStatus, statusResponse.OrderID,
	)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update transaction status",
			Data: err.Error(),
		}
	}

	if statusResponse.TransactionStatus == "capture" {
		subscriptionEnd := time.Now().Add((time.Hour * 24) * 30);

		
		err = tx.QueryRow(
			context.Background(),
			"INSERT INTO subscription(user_id, bank, transaction_time, payment_type, fraud_status, status_code, settlement_time, expiry_time, transaction_status, premium_plan) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING subs_id",
			userId, statusResponse.Bank, statusResponse.TransactionTime.Time(), statusResponse.PaymentType, statusResponse.FraudStatus, statusResponse.StatusCode, statusResponse.SettlementTime.Time(), subscriptionEnd, statusResponse.TransactionStatus, premiumPlan,
		).Scan(&subsId)

		statusResponse.ExpiryTime = dto.MidtransTime(subscriptionEnd);
		statusResponse.PremiumPlan = premiumPlan;


		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to insert subscription data",
				Data: err.Error(),
			}
		}

		_, err = tx.Exec(
			context.Background(),
			"UPDATE user_info SET subscription = 'ACTIVE', premium_plan = 'CREATOR', subs_id = $1 WHERE user_id=$2",
			subsId, userId,
		)

		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to update subscription status",
				Data: err.Error(),
			}
		}
	}


	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Successfully subcribe",
		Data: nil,
	}

}

func Unsubscribe(cookie string) dto.Response {
	conn := config.LoadDatabase()

	var expiryTime time.Time;
	var subsId string;

	userId, err := helper.VerifyAndGet(cookie);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: err.Error(),
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT subs_id from user_info WHERE user_id=$1",
		userId,
	).Scan(&subsId)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get subscription information",
			Data: err.Error(),
		}
	}

	tx, err := conn.Begin(context.Background())

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to start transaction",
			Data: err.Error(),
		}
	}

	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(),
		"SELECT expiry_time FROM subscription WHERE subs_id=$1",
		subsId,
	).Scan(&expiryTime);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get expiry information",
			Data: err.Error(),
		}
	}

	if time.Now().Before(expiryTime) {
		log.Info("True")

		return dto.Response{
			Status: fiber.StatusOK,
			Message: "Subscription still valid",
			Data: nil,
		}

	}

	_, err = tx.Exec(
		context.Background(),
		"UPDATE user_info SET subscription = 'NONE', premium_plan='NONE', subs_id=null WHERE user_id=$1",
		userId,
	)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update user information",
		}
	}

	tx.Commit(context.Background())

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success unsubscribing",
		Data: nil,
	}

}