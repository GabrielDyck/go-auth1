package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)


func AddRoutes(router *mux.Router){
	router.Use(commonMiddleware)
	addHealthCheck(router)
	addSignIn(router)
	addSignUp(router)
	addProfileRoutes(router)
}


func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

