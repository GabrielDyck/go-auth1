package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/front"
	"auth1/pkg/routes/internal"
	"auth1/pkg/routes/internal/singin"
	"auth1/pkg/routes/internal/singup"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomRouter struct {
	client mysql.Client
	configuration config.Configuration
	authService internal.AuthService
}

func NewCustomRouter(client mysql.Client, configuration config.Configuration) CustomRouter{

	return CustomRouter{
		client: client,
		configuration: configuration,
		authService: internal.NewAuthService(client),
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
	signInService := singin.NewSignInService(c.client)

	singin.SignIn(backendRouter, signInService)
	singup.SignUp(backendRouter, c.client)
	internal.Logout(backendRouter)
	internal.ForgotPassword(backendRouter,c.client,expirationDateInMin, emailSender)
	internal.ResetPassword(backendRouter,c.client)
	http.Handle("/backend/",backendRouter)
}
func (c *CustomRouter) AddAuthRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	router.Use(c.secureMiddleware)
	profileService := internal.NewProfileInfoService(c.client, c.authService)
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
		authenticated, err := c.authService.IsAuthorized(token)

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
