package cart

import (
	"fmt"

	"github.com/sahilsnghai/Project6/types"
)

func getCartItemsIds(items []types.CartItem) ([]int, error) {

	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductId)
		}
		productIds[i] = item.ProductId
	}
	return productIds, nil

}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)

	for _, product := range ps {
		productMap[product.Id] = product
	}
	if err := checIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil

	}

	totalPrice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductId]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	orderId, err := h.store.CreateOrder(types.Order{
		UserId:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Some Address",
	})

	if err != nil{
		return 0, 0, nil
	}

	return orderId, totalPrice, nil

}

func calculateTotalPrice(items []types.CartItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range items {
		product := products[item.ProductId]
		total += product.Price * float64(item.Quantity)
	}
	return total
}

func checIfCartIsInStock(items []types.CartItem, products map[int]types.Product) error {

	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := products[item.ProductId]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductId)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)

		}

	}
	return nil

}
