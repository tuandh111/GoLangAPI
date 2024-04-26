package order

import (
	"awesomeProject2/services/auth"
	"awesomeProject2/services/order/types_order"
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
	productStore types.ProductStore
	userStore    types.UserStore
	orderStore   types.OrderStore
}

func NewHandler(productStore types.ProductStore, userStore types.UserStore, orderStore types.OrderStore) *Handler {
	return &Handler{
		productStore: productStore,
		userStore:    userStore,
		orderStore:   orderStore,
	}
}
func (h *Handler) RegisterOrder(router *mux.Router) {
	router.HandleFunc("/create-order", auth.WithJWTAuth(h.CreateOrder, h.userStore)).Methods("POST")
	router.HandleFunc("/find-all-order", auth.WithJWTAuth(h.FindAllOrderWithAdmin, h.userStore)).Methods("GET")
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
	fmt.Println("userId" + strconv.Itoa(userID))
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
func GetUserIDFromContext(ctx context.Context) int {
	var userID, ok = ctx.Value(auth.UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
