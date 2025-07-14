package services

import (
	"context"
	"os"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/helper"
	"github.com/KeyzarRasya/comdel-server/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	Pay(cookie string, plan dto.PremiumPlan) dto.Response;
	Finish(cookie string, transactionStatus dto.TransactionStatus) dto.Response;
	Unsubscribe(cookies string)	dto.Response;
}

type PaymentServiceImpl struct {
	UserRepository repository.UserRepository;
	TransactionRepository repository.TransactionRepository;
	SubscriptionRepository repository.SubscriptionRepository;
}

/*
	Constructor for Creating
	PaymentService
	=====
	Also injecting dependency
*/
func NewPaymentService(
	userRepository repository.UserRepository,
	transactionRepository repository.TransactionRepository,
	subscriptionRepository repository.SubscriptionRepository,
) PaymentService {
	return &PaymentServiceImpl{
		UserRepository: userRepository,
		TransactionRepository: transactionRepository,
		SubscriptionRepository: subscriptionRepository,
	}
}

func (ps *PaymentServiceImpl) Pay(cookie string, plan dto.PremiumPlan) dto.Response {
	var client snap.Client;
	var transaction dto.Transaction;

	conn := config.LoadDatabase();

	userId, err := helper.VerifyAndGet(cookie);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		}
	}

	username, err := ps.UserRepository.GetNameById(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get name",
			Data: err.Error(),
		}
	}

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

	transaction.OrderId = uuid.NewString()

	if plan == dto.Creator {
		transaction.Plan = "CREATOR";
		transaction.Price = 20000;
	} else {
		transaction.Plan = "NEWBIE";
		transaction.Price = 15000;
	}

	req := & snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: transaction.OrderId,
			GrossAmt: transaction.Price,
		},

		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},

		CustomerDetail: &midtrans.CustomerDetails{
			FName: username,
		},

		
	}

	resp, _ := client.CreateTransaction(req)

	err = ps.TransactionRepository.Create(tx, transaction)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to create transaction",
			Data: err.Error(),
		}
	}

	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Creating a transaction",
		Data: resp.RedirectURL,
	}
}

func (ps *PaymentServiceImpl) Finish(cookie string, transaction dto.TransactionStatus) dto.Response {
	conn := config.LoadDatabase()
	
	userId, err := helper.VerifyAndGet(cookie)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get your information",
			Data: err.Error(),
		}
	}

	subscription, err := ps.TransactionRepository.Status(transaction.OrderID);
	
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message:"Failed to parse into status response",
			Data: err.Error(),
		}
	}

	premiumPlan, err := ps.TransactionRepository.GetPremiumPlan(subscription.OrderID)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get plan information",
			Data: err.Error(),
		}
	}

	tx, err := conn.Begin(context.Background());
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to begin transaction",
			Data: err.Error(),
		}
	}
	defer tx.Rollback(context.Background())

	err = ps.TransactionRepository.UpdateTransactionStatus(tx, subscription.TransactionStatus, subscription.OrderID);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update transaction status",
			Data: err.Error(),
		}
	}

	if subscription.TransactionStatus == "capture" {
		subsId, err := ps.SubscriptionRepository.SaveReturningSubsId(tx, userId, subscription, premiumPlan)
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to insert subscription data",
				Data: err.Error(),
			}
		}

		err = ps.SubscriptionRepository.Activate(tx, premiumPlan, subsId, userId)
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
		Data: subscription,
	}

}

func (ps *PaymentServiceImpl) Unsubscribe(cookie string) dto.Response {
	conn := config.LoadDatabase()

	userId, err := helper.VerifyAndGet(cookie);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: err.Error(),
		}
	}

	subsId, err := ps.UserRepository.GetSubsIdById(userId);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get subscription information",
			Data: err.Error(),
		}
	}

	expiryTime, err := ps.SubscriptionRepository.GetExpiryTimeBySubsId(subsId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get expiry information",
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


	if time.Now().Before(expiryTime) {
		log.Info("True")

		return dto.Response{
			Status: fiber.StatusOK,
			Message: "Subscription still valid",
			Data: nil,
		}

	}

	err = ps.UserRepository.DeactivateSubscription(tx, userId)
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