package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mysql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomRouter struct {
	client mysql.Client
	configuration config.Configuration
}

func NewCustomRouter(client mysql.Client, configuration config.Configuration) CustomRouter{
	return CustomRouter{
		client: client,
		configuration: configuration,
	}
}

func (c *CustomRouter) AddRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	healthCheck(router)
	signIn(router, c.client, c.configuration.ExpirationDateInMin)
	signUp(router, c.client)
	getProfileInfo(router, c.client)
	logout(router)
	forgotPassword(router)
	resetPassword(router)
	http.Handle("/",router)
}
func (c *CustomRouter) AddAuthRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	router.Use(c.secureMiddleware)
	editProfileInfo(router, c.client)
	logout(router)


	http.Handle("/auth/",router)
}

func (c *CustomRouter)commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (c *CustomRouter) secureMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := w.Header().Get("AUTHORIZATION")
		authenticated, err := c.client.IsAuthenticated(token)

		if err !=nil{
			wrapInternalErrorResponse(w,err)
			return
		}

		if !authenticated{
			wrapBadRequestResponse(w,errors.New("not Authenticated"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
