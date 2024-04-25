package product

import (
	"awesomeProject2/types"
	"awesomeProject2/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		productStore: productStore,
		userStore:    userStore,
	}
}
func (h *Handler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/get-product-by-id/{id}", h.handleGetProductByID).Methods("GET")
	router.HandleFunc("/GetProductsByID", h.handleGetProductsByID).Methods("POST")
	router.HandleFunc("/GetProductsPage", h.handleGetProductsPage).Queries("page", "{page}").Methods(http.MethodGet)
	router.HandleFunc("/GetProducts", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/CreateProduct", h.handleCreateProduct).Methods("POST")
	router.HandleFunc("/UpdateProduct", h.handleUpdateProduct).Methods("POST")
	router.HandleFunc("/DeleteProduct", h.handleDeleteProduct).Methods("POST")

}
func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error"))
		return
	}
	fmt.Println(str)
	productId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error convert product with id: %d", productId))
		return
	}
	product, errs := h.productStore.GetProductByID(productId)
	if errs != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error: %v", errs))
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleGetProductsByID(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) handleGetProductsPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageStr := vars["page"]
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}
	products, err := h.productStore.GetProductsPage(page, 4)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productStore.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {

}
