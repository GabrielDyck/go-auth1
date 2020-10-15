package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/front"
	"auth1/pkg/routes/internal"
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

func (c *CustomRouter) AddFrontendRoutes() {
	htmlRouter := mux.NewRouter()
	customRouter:= front.NewFrontRouter()
	customRouter.AddRoutes(htmlRouter)
}
func (c *CustomRouter) AddBackendRoutes(backendRouter *mux.Router, expirationDateInMin int, emailSender mail.Sender) {
	backendRouter.Use(c.commonMiddleware)

	internal.HealthCheck(backendRouter)
	signInService := internal.NewSignInService(c.client)

	internal.SignIn(backendRouter, signInService)
	internal.SignUp(backendRouter, c.client)
	internal.Logout(backendRouter)
	internal.ForgotPassword(backendRouter,c.client,expirationDateInMin, emailSender)
	internal.ResetPassword(backendRouter,c.client)
	http.Handle("/backend/",backendRouter)
}
func (c *CustomRouter) AddAuthRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	router.Use(c.secureMiddleware)
	profileService := internal.NewProfileInfoService(c.client)
	profileService.GetProfileInfo(router)
	profileService.EditProfileInfo(router)
	internal.Logout(router)


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
		token := r.Header.Get("Authorization")
		authenticated, err := c.client.IsAuthenticated(token)

		if err !=nil{
			internal.WrapInternalErrorResponse(w,err)
			return
		}

		if !authenticated{
			internal.WrapNotAllowedRequestResponse(w,errors.New("not Authenticated"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
