package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sahilsnghai/Project6/service/auth"
	"github.com/sahilsnghai/Project6/types"
	"github.com/sahilsnghai/Project6/utils"
)

type Handler struct {
	productStore types.ProductStore
	userStore    types.UserStore
}

func Newhandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{productStore: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(route *mux.Router) {
	route.HandleFunc("/products", auth.WithJWTAuth(h.handleGetProducts, h.userStore)).Methods("GET")
	route.HandleFunc("/products/{productId}", auth.WithJWTAuth(h.handleGetProduct, h.userStore)).Methods("GET")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {

	ps, err := h.productStore.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, ps)

}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	str, ok := vars["productId"]

	if !ok {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return

	}

	productId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	product, err := h.productStore.GetProductByID(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, product)
}
