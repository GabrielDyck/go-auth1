package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mysql"
	"github.com/gorilla/mux"
	"net/http"
)


func AddRoutes(router *mux.Router,client mysql.Client, configuration config.Configuration){
	router.Use(commonMiddleware)
	healthCheck(router)
	signIn(router,client, configuration.ExpirationDateInMin)
	signUp(router, client)
	getProfileInfo(router,client)
	logout(router)
	forgotPassword(router)
	resetPassword(router)
}


func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

