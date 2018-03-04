package interfaces

import (
	"domain"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"usecases"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type OrderInteractor interface {
	Items(userId, orderId int) ([]usecases.Item, error)
	GetOrder(userId, orderId int) (domain.Order, error)
	Add(userId, orderId, itemId int) error
}

type WebserviceHandler struct {
	OrderInteractor OrderInteractor
}

func (handler WebserviceHandler) ShowOrder(res http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.FormValue("userId"))
	orderId, _ := strconv.Atoi(req.FormValue("orderId"))
	items, _ := handler.OrderInteractor.Items(userId, orderId)
	for _, item := range items {
		io.WriteString(res, fmt.Sprintf("item id: %d\n", item.Id))
		io.WriteString(res, fmt.Sprintf("item name: %v\n", item.Name))
		io.WriteString(res, fmt.Sprintf("item value: %f\n", item.Value))
	}
}

func (handler WebserviceHandler) PlaceOrder(res http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.FormValue("userId"))
	orderId, _ := strconv.Atoi(req.FormValue("orderId"))
	token := req.FormValue("stripeToken")
	order, _ := handler.OrderInteractor.GetOrder(userId, orderId)

	charge := ProcessPayment(order, token)
	chargeAmount := float64(charge.Amount) / 100
	// order := interactor.OrderRepository.FindById(orderId)
	io.WriteString(res, fmt.Sprintf("order: %d, costs $%f and charged you %f %s: %s\n", order.Id, order.Value(), chargeAmount, charge.Currency, charge.Status))
}

func ProcessPayment(order domain.Order, token string) *stripe.Charge {
	dollaz := uint64(order.Value() * 100)
	// TODO: Set your secret key: remember to change this to your live secret key in production
	// See your keys here: https://dashboard.stripe.com/account/apikeys

	stripe.Key = "sk_test_BQokikJOvBiI2HlWgH4olfQ2"

	// Token is created using Checkout or Elements!
	// Get the payment token ID submitted by the form:

	// Charge the user's card:
	params := &stripe.ChargeParams{
		Amount:   dollaz,
		Currency: "usd",
	}
	params.SetSource(token)

	charge, _ := charge.New(params)
	return charge
}
