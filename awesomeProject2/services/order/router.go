package order

import (
	"awesomeProject2/services/auth"
	"awesomeProject2/services/order/types_order"
	"awesomeProject2/services/product/types_product"
	"awesomeProject2/services/status"
	"awesomeProject2/services/user/types_user"
	"awesomeProject2/types"
	"awesomeProject2/utils"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	productStore types_product.ProductStore
	userStore    types_user.UserStore
	orderStore   types_order.OrderStore
}

func NewHandler(productStore types_product.ProductStore, userStore types_user.UserStore, orderStore types_order.OrderStore) *Handler {
	return &Handler{
		productStore: productStore,
		userStore:    userStore,
		orderStore:   orderStore,
	}
}
func (h *Handler) RegisterOrder(router *mux.Router) {
	router.HandleFunc("/create-order", auth.WithJWTAuth(h.CreateOrder, h.userStore)).Methods("POST")
	router.HandleFunc("/find-all-order", auth.WithJWTAuth(h.FindAllOrderWithAdmin, h.userStore)).Methods("GET")
	router.HandleFunc("/order/{userId}", auth.WithJWTAuth(h.FindByOrderByUser, h.userStore)).Methods("GET")
	router.HandleFunc("/update-order", auth.WithJWTAuth(h.UpdateOrder, h.userStore)).Methods("POST")
	router.HandleFunc("/delete-order/{id}", auth.WithJWTAuth(h.deleteOrder, h.userStore)).Methods(http.MethodDelete)
}
func (h *Handler) FindAllOrderWithAdmin(w http.ResponseWriter, r *http.Request) {

	order, err := h.orderStore.FindAllOrderWithAdmin()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, order)
}
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	var odpayload types_order.OrderPayload
	if err := utils.ParseJSON(r, &odpayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error json order"))
		return
	}
	if err := utils.Validate.Struct(odpayload); err != nil {
		erros := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, erros)
		return
	}
	id, err := h.orderStore.CreateOrder(
		types_order.OrderPayload{
			UserId:  strconv.Itoa(userID),
			Total:   odpayload.Total,
			Status:  odpayload.Status,
			Address: odpayload.Address,
		})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "successfully create order with id: "+strconv.Itoa(id))
}
func (s *Handler) FindByOrderByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["userId"]
	userId, err := strconv.Atoi(Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	order, erors := s.orderStore.FindByOrderUserId(userId)
	if erors != nil {
		utils.WriteError(w, http.StatusBadRequest, erors)
		return
	}
	utils.WriteJSON(w, http.StatusOK, order)
}
func (s *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIDFromContext(r.Context())
	var orderPayload types_order.OrderUpdateUserID
	if err := utils.ParseJSON(r, &orderPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err := utils.Validate.Struct(orderPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	message, errs := s.orderStore.UpdateOrderByUserId(orderPayload, userId, status.GetPending)
	if errs != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf(message))
		return
	}
	utils.WriteJSON(w, http.StatusOK, message)
}
func (h *Handler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	orderId, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	message, errDeleteOrder := h.orderStore.DeleteOrder(orderId)
	if errDeleteOrder != nil {
		utils.WriteError(w, http.StatusBadRequest, errDeleteOrder)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.JsonResponse{Message: message})

}
func GetUserIDFromContext(ctx context.Context) int {
	var userID, ok = ctx.Value(auth.UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
