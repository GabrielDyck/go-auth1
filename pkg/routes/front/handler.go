package front

import (
	"auth1/pkg/mysql"
	"auth1/pkg/routes/front/internal"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type FrontRouter struct {
	db mysql.Auth
}

func NewFrontRouter() FrontRouter {
	return FrontRouter{}
}

func (c *FrontRouter) AddRoutes(router *mux.Router) {
	http.Handle("/",router)
	router.Use(c.commonMiddleware)
	internal.HealthCheck(router)
	internal.AddResources(router)
	inflator:= internal.NewTemplateInflator()
	inflator.AddTemplates(router)
}

func (c *FrontRouter)commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}

func (c *FrontRouter) secureMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := w.Header().Get("AUTHORIZATION")
		authenticated, err := c.db.IsAuthenticated(token)

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