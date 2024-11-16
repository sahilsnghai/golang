package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sahilsnghai/Project6/service/auth"
	"github.com/sahilsnghai/Project6/types"
	"github.com/sahilsnghai/Project6/utils"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func Newhandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(route *mux.Router) {
	route.HandleFunc("/products", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {

	var cart types.CartCheckoutPayload

	if err := utils.ParseJson(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	productId, err := getCartItemsIds(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByID(productId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userId := 0

	orderId, totalPrice, err := h.createOrder(ps, cart.Items, userId)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderId,
	})

}
