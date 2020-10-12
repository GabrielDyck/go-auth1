package internal

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AddResources(router *mux.Router){
	router.Handle("/resources/", http.StripPrefix("/resources/",
		http.FileServer(http.Dir("./../../resources/"))))

}