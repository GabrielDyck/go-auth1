package internal

import (
	"auth1/api"
	"auth1/pkg/routes/front/internal/templates"
	"auth1/pkg/routes/internal"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func (t *TemplateBuilder) parseEditProfileTemplate(req *http.Request) ([]byte, int) {

	log.Print(req)
	user, err := req.Cookie("User")
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusMethodNotAllowed)
	}
	authCookie, err := req.Cookie("Authorization")
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusMethodNotAllowed)
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
	if status == http.StatusMethodNotAllowed {
		rw.Write([]byte(templates.RedirectAuthenticationError))
	} else if status != 0 {
		tmpl, err := template.New("error.html").Parse(templates.ErrorTemplate)
		if err != nil {
			log.Println(err)
		}
		var errorMsg api.ErrorMSG

		err = json.Unmarshal(data, &errorMsg)
		if err != nil {
			log.Println(err)
		}

		err = tmpl.Execute(rw, errorMsg)
		if err != nil {
			log.Println(err)
		}

	} else {
		err := json.Unmarshal(data, &account)

		tmpl := template.Must(template.New("edit-profile.html").Funcs(
			template.FuncMap{
				"showIfNotNil": func(value *string) string {

			if value != nil {
				return *value
			}
			return ""
		},
		"disabled": func(value api.AccountType) string{
			if value != api.Basic{
				return "disabled"
			}
			return ""
		},
			}).ParseFiles("pkg/routes/front/internal/templates/edit-profile.html"))
		err = tmpl.Execute(rw, account)
		if err != nil {
			log.Println(err)
		}
	}
}
