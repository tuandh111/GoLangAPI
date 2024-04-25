package api

import (
	"awesomeProject2/services/product"
	"awesomeProject2/services/user"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRouter(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRouter(subrouter)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
