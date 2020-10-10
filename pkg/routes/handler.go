package routes

import (
	"auth1/pkg/mysql"
	"github.com/gorilla/mux"
	"net/http"
)


func AddRoutes(router *mux.Router,client mysql.Client){
	router.Use(commonMiddleware)
	addHealthCheck(router)
	addSignIn(router)
	addSignUp(router, client)
	addProfileRoutes(router)
	addLogout(router)
	addForgotPassword(router)
	addResetPassword(router)
}


func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

