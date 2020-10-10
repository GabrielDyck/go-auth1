package routes

import (
	"github.com/gorilla/mux"
)


func AddRoutes(router *mux.Router){
	router.Use(commonMiddleware)
	addHealthCheck(router)
	addSignUp(router)
}



