package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/backend/auth"
	"auth1/pkg/routes/backend/forgot"
	"auth1/pkg/routes/backend/healthcheck"
	"auth1/pkg/routes/backend/logout"
	"auth1/pkg/routes/backend/profile"
	"auth1/pkg/routes/backend/resetpassword"
	"auth1/pkg/routes/backend/signup"
	"auth1/pkg/routes/backend/singin"
	"auth1/pkg/routes/front"
	"auth1/pkg/routes/internal"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomRouter struct {
	client        mysql.Client
	configuration config.Configuration
	authService   auth.AuthService
}

func NewCustomRouter(client mysql.Client, configuration config.Configuration) CustomRouter {

	return CustomRouter{
		client:        client,
		configuration: configuration,
		authService:   auth.NewAuthService(client),
	}
}

func (c *CustomRouter) AddFrontendRoutes() {
	htmlRouter := mux.NewRouter()
	customRouter := front.NewFrontRouter()
	customRouter.AddRoutes(htmlRouter)
}
func (c *CustomRouter) AddBackendRoutes(backendRouter *mux.Router, expirationDateInMin int, emailSender mail.Sender) {
	backendRouter.Use(c.commonMiddleware)

	c.authService.AddRoutes(backendRouter)

	healthcheck.HealthCheck(backendRouter)

	signInService := singin.NewSignInService(c.client)
	signInService.AddRoutes(backendRouter)

	signupService := signup.NewSignUpService(c.client)
	signupService.AddRoutes(backendRouter)

	forgotService := forgot.NewForgotPasswordService(c.client, expirationDateInMin, emailSender)
	forgotService.AddRoutes(backendRouter)

	resetPasswordService := resetpassword.NewResetPasswordService(c.client)
	resetPasswordService.AddRoutes(backendRouter)

	http.Handle("/backend/", backendRouter)
}
func (c *CustomRouter) AddAuthRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	router.Use(c.secureMiddleware)
	profileService := profile.NewProfileInfoService(c.client, c.authService)
	profileService.AddRoutes(router)

	logoutService := logout.NewLogoutService(c.client)
	logoutService.AddRoutes(router)

	http.Handle("/auth/", router)
}

func (c *CustomRouter) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (c *CustomRouter) secureMiddleware(next http.Handler) http.Handler {

	secure := func(request *http.Request) ([]byte, int) {
		token := request.Header.Get("Authorization")
		authenticated, err := c.authService.IsAuthorized(token)

		if err != nil {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
		}

		if !authenticated {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("not Authenticated")), http.StatusMethodNotAllowed)
		}
		return  []byte("{}"), 0
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		resp, status := secure(request)
		if status != 0 {
			writer.WriteHeader(status)
			writer.Write(resp)
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}
