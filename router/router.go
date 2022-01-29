package router

import (
	"simple_api_postgres/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/fallback/todo/{id}", middleware.GetTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/fallback/todo", middleware.GetAllTodo).Methods("GET", "OPTIONS")
	router.HandleFunc("/fallback/todo", middleware.CreateTodo).Methods("POST", "OPTIONS")
	router.HandleFunc("/fallback/todo/{id}", middleware.UpdateTodo).Methods("PUT", "OPTIONS")
	router.HandleFunc("/fallback/deleteTodo/{id}", middleware.DeleteTodo).Methods("DELETE", "OPTIONS")

	return router
}
