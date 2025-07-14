package handlers

import (
	"fmt"
	"strings"

	"comdel-backend/internal/dto"
	"comdel-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandlers interface{
	CreatePayment(c *fiber.Ctx)		error;
	FinishPayment(c *fiber.Ctx)		error
}

type PaymentHandlersImpl struct {
	PaymentService services.PaymentService
}

func NewPaymentHandlers(paymentService services.PaymentService) PaymentHandlers {
	return &PaymentHandlersImpl{PaymentService: paymentService}
}

func (ph *PaymentHandlersImpl) CreatePayment(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	planQuery := c.Query("plan")
	var premiumPlan dto.PremiumPlan;

	if strings.ToUpper(planQuery) == "CREATOR" {
		premiumPlan = dto.Creator;
	} else {
		premiumPlan = dto.Newbie
	}

	var response dto.Response = ph.PaymentService.Pay(cookie, premiumPlan);

	return c.Redirect(fmt.Sprint(response.Data))
}

func (ph *PaymentHandlersImpl) FinishPayment(c *fiber.Ctx) error {
	queries := c.Queries();
	cookie := c.Cookies("jwt");

	transactionStatus := dto.TransactionStatus {
		StatusCode: queries["status_code"],
		OrderID: queries["order_id"],
		Status: queries["transaction_status"],
	}

	var response dto.Response = ph.PaymentService.Finish(cookie, transactionStatus);

	if response.Status != fiber.StatusOK {
		return c.JSON(response.JSON())
	}

	var redirectUrl string = fmt.Sprintf("http://localhost:5173/payment/finish/?status_code=%s&transaction_status=%s&order_id=%s", transactionStatus.StatusCode, transactionStatus.Status, transactionStatus.OrderID);

	return c.Redirect(redirectUrl)
}

// func UnsubscribePlan(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")

// 	var response dto.Response = 
// }