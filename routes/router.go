package routes

import (
	"FTBank/handlers"
	"FTBank/middleware"
	"FTBank/repository"
	"FTBank/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)


func Router(db *gorm.DB) *mux.Router {
	repo := repository.NewBankRepository(db)
	service := service.NewBankService(repo)
	handlers := handlers.NewBankHandler(*service)

	r := mux.NewRouter()

	// define public routes
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")

	// define protected routes
	protectedRoutes := r.PathPrefix("/api/v1").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	return r
}