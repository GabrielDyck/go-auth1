package internal

import (
	"auth1/pkg/routes/internal"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type TemplateInflator struct {
	signInTemplate         []byte
	signUpTemplate         []byte
	editProfileTemplate    []byte
	profileInfoTemplate    []byte
	forgotPasswordTemplate []byte
	resetPasswordTemplate  []byte
}

func NewTemplateInflator() TemplateInflator {
	loginTemplate, err := ioutil.ReadFile("resources/html/signin.html")
	panicIfError(err)
	signUpTemplate, err := ioutil.ReadFile("resources/html/signup.html")
	panicIfError(err)

	editProfileTemplate, err := ioutil.ReadFile("resources/html/edit-profile.html")
	panicIfError(err)

	profileInfoTemplate, err := ioutil.ReadFile("resources/html/profile-info.html")
	panicIfError(err)

	forgotPasswordTemplate, err := ioutil.ReadFile("resources/html/forgot-password.html")
	panicIfError(err)

	resetPasswordTemplate, err := ioutil.ReadFile("resources/html/reset-password.html")
	panicIfError(err)
	return TemplateInflator{
		signInTemplate:         loginTemplate,
		signUpTemplate:         signUpTemplate,
		editProfileTemplate:    editProfileTemplate,
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

func (t *TemplateInflator) AddTemplates(router *mux.Router) {
	router.HandleFunc("/", t.inflateSignInTemplate)
	router.HandleFunc("/signin", t.inflateSignInTemplate)
	router.HandleFunc("/signup", t.inflateSignUpTemplate)
	router.HandleFunc("/edit-profile/{id}", t.inflateEditProfileTemplate)
	router.HandleFunc("/profile-info/{id}", t.inflateProfileInfoTemplate)
	router.HandleFunc("/forgot-password", t.inflateForgotPasswordTemplate)
	router.HandleFunc("/reset-password", t.inflateResetPasswordTemplate)

}

func (t *TemplateInflator) inflateSignInTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.signInTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateInflator) inflateSignUpTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.signUpTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateInflator) inflateEditProfileTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.editProfileTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateInflator) inflateProfileInfoTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.profileInfoTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateInflator) inflateForgotPasswordTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.forgotPasswordTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}

func (t *TemplateInflator) inflateResetPasswordTemplate(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(t.resetPasswordTemplate)
	if err != nil {
		internal.WrapInternalErrorResponse(rw, err)
		return
	}
}
