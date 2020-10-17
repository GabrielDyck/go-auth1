package internal

import (
	"auth1/pkg/routes/internal"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)


type TemplateBuilder struct {
	signInTemplate         []byte
	signUpTemplate         []byte
	forgotPasswordTemplate []byte
	resetPasswordTemplate  []byte
}

func NewTemplateInflator() TemplateBuilder {
	loginTemplate, err := ioutil.ReadFile("resources/html/signin.html")
	panicIfError(err)
	signUpTemplate, err := ioutil.ReadFile("resources/html/signup.html")
	panicIfError(err)


	forgotPasswordTemplate, err := ioutil.ReadFile("resources/html/forgot-password.html")
	panicIfError(err)

	resetPasswordTemplate, err := ioutil.ReadFile("resources/html/reset-password.html")
	panicIfError(err)
	return TemplateBuilder{
		signInTemplate:         loginTemplate,
		signUpTemplate:         signUpTemplate,
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
