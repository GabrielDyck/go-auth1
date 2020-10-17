package internal

import (
	"auth1/api"
	"auth1/pkg/routes/internal"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const errorTemplate string = `
<html>
<body>
<h1> Error ocurred </h1>
</body>
</html>`

type TemplateBuilder struct {
	signInTemplate         []byte
	signUpTemplate         []byte
	profileInfoTemplate    []byte
	forgotPasswordTemplate []byte
	resetPasswordTemplate  []byte
}

func NewTemplateInflator() TemplateBuilder {
	loginTemplate, err := ioutil.ReadFile("resources/html/signin.html")
	panicIfError(err)
	signUpTemplate, err := ioutil.ReadFile("resources/html/signup.html")
	panicIfError(err)

	profileInfoTemplate, err := ioutil.ReadFile("resources/html/profile-info.html")
	panicIfError(err)

	forgotPasswordTemplate, err := ioutil.ReadFile("resources/html/forgot-password.html")
	panicIfError(err)

	resetPasswordTemplate, err := ioutil.ReadFile("resources/html/reset-password.html")
	panicIfError(err)
	return TemplateBuilder{
		signInTemplate:         loginTemplate,
		signUpTemplate:         signUpTemplate,
		profileInfoTemplate:    profileInfoTemplate,
		forgotPasswordTemplate: forgotPasswordTemplate,
		resetPasswordTemplate:  resetPasswordTemplate,
	}

}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *TemplateBuilder) AddTemplates(router *mux.Router) {
	router.HandleFunc("/", t.inflateSignInTemplate)
	router.HandleFunc("/signin", t.inflateSignInTemplate)
	router.HandleFunc("/signup", t.inflateSignUpTemplate)
	router.HandleFunc("/edit-profile", t.inflateEditProfileTemplate)
	router.HandleFunc("/profile-info/{id}", t.inflateProfileInfoTemplate)
	router.HandleFunc("/forgot-password", t.inflateForgotPasswordTemplate)
	router.HandleFunc("/reset-password", t.inflateResetPasswordTemplate)

}

func (t *TemplateBuilder) inflateSignInTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.signInTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateBuilder) inflateSignUpTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.signUpTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateBuilder) parseEditProfileTemplate(req *http.Request) ([]byte, int) {

	log.Print(req)
	user, err := req.Cookie("User")
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}
	authCookie, err := req.Cookie("Authorization")
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}
	client := &http.Client{}
	newReq, err := http.NewRequest("GET", os.Getenv("DOMAIN")+os.Getenv("PORT")+"/auth/profile-info/"+user.Value, nil)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}
	newReq.Header.Set("Authorization", authCookie.Value)
	newReq.Header.Set("Content-Type", "application/json")

	response, err := client.Do(newReq)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	if response.StatusCode != 200 {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
		}

		return data, response.StatusCode
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}
	return data, 0
}
func (t *TemplateBuilder) inflateEditProfileTemplate(rw http.ResponseWriter, req *http.Request) {

	var account api.Account
	data, status := t.parseEditProfileTemplate(req)
	if status != 0 {
		rw.WriteHeader(status)
		rw.Write([]byte(errorTemplate))

	} else {

		err := json.Unmarshal(data, &account)

		tmpl := template.Must(template.New("edit-profile.html").Funcs(template.FuncMap{"showIfNotNil": func(value *string) string {

			if value != nil {
				return *value
			}
			return ""
		}}).ParseFiles("pkg/routes/front/internal/templates/edit-profile.html"))
		err = tmpl.Execute(rw, account)
		if err != nil {
			data, httpStatus := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			internal.WrapResponse(rw, data, httpStatus)
			return
		}
	}
}

func (t *TemplateBuilder) inflateProfileInfoTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.profileInfoTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateBuilder) inflateForgotPasswordTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.forgotPasswordTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateBuilder) inflateResetPasswordTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.resetPasswordTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}
