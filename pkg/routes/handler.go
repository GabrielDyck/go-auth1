package routes

import (
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
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

func (c *CustomRouter) AddRoutes(router *mux.Router, expirationDateInMin int, emailSender mail.Sender) {
	router.Use(c.commonMiddleware)
	internal.HealthCheck(router)
	signInService := internal.NewSignInService(c.client)

	internal.SignIn(router, signInService)
	internal.SignUp(router, c.client)
	internal.GetProfileInfo(router, c.client)
	internal.Logout(router)
	internal.ForgotPassword(router,c.client,expirationDateInMin, emailSender)
	internal.ResetPassword(router,c.client)
	http.Handle("/",router)
}
func (c *CustomRouter) AddAuthRoutes(router *mux.Router) {
	router.Use(c.commonMiddleware)
	router.Use(c.secureMiddleware)
	profileService := internal.NewProfileInfoService(c.client)
	internal.EditProfileInfo(router, profileService )
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
		token := w.Header().Get("AUTHORIZATION")
		authenticated, err := c.client.IsAuthenticated(token)

		if err !=nil{
			internal.WrapInternalErrorResponse(w,err)
			return
		}

		if !authenticated{
			internal.WrapBadRequestResponse(w,errors.New("not Authenticated"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
