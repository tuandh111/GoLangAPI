package checkout

import (
	"awesomeProject2/services/auth"
	"awesomeProject2/services/checkout/types_checkout"
	"awesomeProject2/services/user/types_user"
	"awesomeProject2/types"
	"awesomeProject2/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	CheckOutHandler types_checkout.CheckOutStore
	UserStore       types_user.UserStore
}

func NewHandler(checkOutStore types_checkout.CheckOutStore, UserStore types_user.UserStore) *Handler {
	return &Handler{
		CheckOutHandler: checkOutStore,
		UserStore:       UserStore,
	}
}
func (h *Handler) RegisterCheckout(router *mux.Router) {
	router.HandleFunc("/update-checkout-status-admin", auth.WithJWTAuth(h.updateStatusAmin, h.UserStore)).Methods(http.MethodPost)
}
func (h *Handler) updateStatusAmin(w http.ResponseWriter, r *http.Request) {
	var orderUpdateStatusAdmin types_checkout.OrderCheckout
	if err := utils.ParseJSON(r, &orderUpdateStatusAdmin); err != nil {
		utils.WriteError(w, http.StatusOK, err)
		return
	}
	if err := utils.Validate.Struct(orderUpdateStatusAdmin); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusOK, errors)
		return
	}
	message, err := h.CheckOutHandler.UpdateStatusAdmin(orderUpdateStatusAdmin.ID, orderUpdateStatusAdmin.Status)
	if err != nil {
		utils.WriteError(w, http.StatusOK, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.JsonResponse{Message: message})
}
