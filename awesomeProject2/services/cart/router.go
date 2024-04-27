package cart

import (
	"awesomeProject2/services/auth"
	"awesomeProject2/services/cart/types_cart"
	"awesomeProject2/services/order"
	"awesomeProject2/services/order/types_order"
	"awesomeProject2/services/product/types_product"
	"awesomeProject2/services/status"
	"awesomeProject2/services/user/types_user"
	"awesomeProject2/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	CartStore    types_cart.CartStore
	ProductStore types_product.ProductStore
	OrderStore   types_order.OrderStore
	UserStore    types_user.UserStore
}

func NewHandler(CartStore types_cart.CartStore, ProductStore types_product.ProductStore, OrderStore types_order.OrderStore, UserStore types_user.UserStore) *Handler {
	return &Handler{
		CartStore:    CartStore,
		ProductStore: ProductStore,
		OrderStore:   OrderStore,
		UserStore:    UserStore,
	}
}
func (h *Handler) RegisterCart(router *mux.Router) {
	router.HandleFunc("/get-order-id", auth.WithJWTAuth(h.getOrderId, h.UserStore)).Methods(http.MethodGet)
	router.HandleFunc("/update-or-save-order-and-product", auth.WithJWTAuth(h.updateCartOrderAndProduct, h.UserStore)).Methods(http.MethodPost)
	//testAPI
	router.HandleFunc("/check-cart-order-and-product", auth.WithJWTAuth(h.checkCartOrderAndProduct, h.UserStore)).Methods(http.MethodGet)
}
func (h *Handler) getOrderId(w http.ResponseWriter, r *http.Request) {
	userId := order.GetUserIDFromContext(r.Context())
	order, err := h.OrderStore.FindByOrderUserIdAndStatus(userId, status.GetPending)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	cart, errs := h.CartStore.CartItems(order.ID)
	if errs != nil {
		utils.WriteError(w, http.StatusBadRequest, errs)
		return
	}
	utils.WriteJSON(w, http.StatusOK, cart)
}
func (h *Handler) updateCartOrderAndProduct(w http.ResponseWriter, r *http.Request) {
	userId := order.GetUserIDFromContext(r.Context())
	var orderItemUpdateOrSave types_cart.CartItemUpdate
	if err := utils.ParseJSON(r, &orderItemUpdateOrSave); err != nil {
		fmt.Println(err)
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	order, err := h.OrderStore.FindByOrderUserIdAndStatus(userId, status.GetPending)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println(orderItemUpdateOrSave)
	cart, number := h.CartStore.CheckCartOrderIdAndProductId(order.ID, orderItemUpdateOrSave.ProductID)
	if number == 0 {
		if er := h.OrderStore.CreateOrderItem(types_order.OrderItem{
			OrderID:   orderItemUpdateOrSave.OrderID,
			ProductID: orderItemUpdateOrSave.ProductID,
			Quantity:  orderItemUpdateOrSave.Quantity,
			Price:     orderItemUpdateOrSave.Price,
		}); er != nil {
			utils.WriteError(w, http.StatusBadRequest, er)
			return
		}
		utils.WriteJSON(w, http.StatusOK, "save successfully")
	} else {
		message, errs := h.CartStore.UpdateOrSaveOrderIdAndProductId(types_cart.CartItemUpdate{
			OrderID:   orderItemUpdateOrSave.OrderID,
			ProductID: orderItemUpdateOrSave.ProductID,
			Quantity:  orderItemUpdateOrSave.Quantity,
			Price:     orderItemUpdateOrSave.Price,
		}, cart.Id)
		fmt.Println(errs)
		if errs != nil {
			utils.WriteError(w, http.StatusBadRequest, errs)
			return
		}
		utils.WriteJSON(w, http.StatusOK, message)
	}

}

// testAPI
func (h *Handler) checkCartOrderAndProduct(w http.ResponseWriter, r *http.Request) {
	userId := order.GetUserIDFromContext(r.Context())
	order, err := h.OrderStore.FindByOrderUserIdAndStatus(userId, status.GetPending)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	cart, number := h.CartStore.CheckCartOrderIdAndProductId(order.ID, 9)
	if number == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("The product already exists in the shopping cart"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, cart)
}
