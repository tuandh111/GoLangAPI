package user

import (
	"awesomeProject2/configs"
	"awesomeProject2/services/auth"
	"awesomeProject2/types"
	"awesomeProject2/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	//admin
	router.HandleFunc("/users/{userId}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/delete-by-user-id/{userId}", auth.WithJWTAuth(h.handleDeleteByUser, h.store)).Methods(http.MethodDelete)
	//router.HandleFunc("/get-all-user", auth.WithJWTAuth(h.handleGetAllUser, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/get-all-user", auth.WithJWTAuth(h.handleGetAllUserPage, h.store)).Queries("page", "{page}").Methods(http.MethodGet)
	router.HandleFunc("/get-search-user/{lastname}", auth.WithJWTAuth(h.handleSearchName, h.store)).Methods(http.MethodGet)
}
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
func (h *Handler) handleSearchName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["lastname"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user lastname"))
		return
	}
	user, err := h.store.FindBySearchName(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found lastname"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		fmt.Println(err)
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		fmt.Println(errors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	//check email exits
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}
	//hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)

}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}
	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}
	user, err := h.store.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
func (h *Handler) handleDeleteByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}
	userId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.DeleteUserByID(userId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("delete fail by id: %d", userId))
		return
	}
	utils.WriteJSON(w, http.StatusOK, "delete successfully with user id: "+strconv.Itoa(userId))
}
func (h *Handler) handleGetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUserId()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
func (h *Handler) handleGetAllUserPage(w http.ResponseWriter, r *http.Request) {
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
	users, err := h.store.GetAllUserIdPage(page, 7)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
