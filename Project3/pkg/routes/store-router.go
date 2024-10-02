package routes

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/sahilsnghai/bookstore/pkg/controllers"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {

	bookRoute := "/book/"
	bookRouteId := "{Id}"
	fmt.Printf("\nHere are the details of Routes\n\nCreate Book %v POST\nGet Book %v GET\nGet Book by ID %v GET\nUpdate Book %v POST\nDelete Book %v POST\n\n",
		bookRoute, bookRoute, bookRoute+bookRouteId, bookRoute+bookRouteId, bookRoute+bookRouteId)

	router.HandleFunc(bookRoute, controllers.CreateBook).Methods("POST")
	router.HandleFunc(bookRoute, controllers.GetBook).Methods("GET")
	router.HandleFunc(bookRoute+bookRouteId, controllers.GetBookById).Methods("GET")
	router.HandleFunc(bookRoute+bookRouteId, controllers.UpdateBook).Methods("PUT")
	router.HandleFunc(bookRoute+bookRouteId, controllers.DeleteBook).Methods("DELETE")

}
