package interfaces

import (
	"domain"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"usecases"
	//"github.com/bogdanx/go-cleanarchitecture/src/domain"
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
	order, _ := handler.OrderInteractor.GetOrder(userId, orderId)
	// order := interactor.OrderRepository.FindById(orderId)
	io.WriteString(res, fmt.Sprintf("order: %d, costs $%f\n", order.Id, order.Value()))
}
