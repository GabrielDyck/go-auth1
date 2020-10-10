package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	profile="/profile/{id}"
)

func addProfileRoutes(router *mux.Router){
	router.HandleFunc(profile, func(writer http.ResponseWriter, request *http.Request) {
		id := mux.Vars(request)["id"]

		fmt.Println(id)

	}).Methods("GET")


}