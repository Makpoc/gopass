package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Makpoc/gopass/generator"
)

type errorMessage struct {
	Message string
	Err     error
}

func homePageHandler(w http.ResponseWriter, _ *http.Request) {
	mustExecuteTemplate(w, func() error { return formPage.Execute(w, nil) })
}

func generatePageHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := parseForm(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, errorMessage{Message: "Failed to parse form parameters.", Err: err})
		return
	}
	pass, err := generator.GeneratePassword(settings)
	if err != nil {
		handleError(w, http.StatusInternalServerError, errorMessage{Message: "Failed to parse form parameters.", Err: err})
		return
	}

	mustExecuteTemplate(w, func() error { return resultPage.Execute(w, string(pass)) })
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	mustExecuteTemplate(w, func() error { return aboutPage.Execute(w, nil) })
}

func parseForm(r *http.Request) (settings generator.Settings, err error) {
	settings = generator.DefaultSettings()

	masterPass := r.PostFormValue("password")
	settings.MasterPhrase = masterPass

	settings.Domain = r.PostFormValue("domain")
	settings.AdditionalInfo = r.PostFormValue("additional-info")

	pLength := r.PostFormValue("password-length")
	if pLength != "" {
		var passwordLength int
		if passwordLength, err = strconv.Atoi(pLength); err != nil {
			return
		}
		settings.PasswordLength = passwordLength
	}

	sChars := r.PostFormValue("special-chars")
	settings.AddSpecialCharacters = sChars == "on"

	return
}

func handleError(w http.ResponseWriter, status int, err errorMessage) {
	w.WriteHeader(status)
	templateErr := errorPage.Execute(w, err)
	if templateErr != nil {
		log.Printf("Failed to execute error template: %v", templateErr)
	}
}

func mustExecuteTemplate(w http.ResponseWriter, execTemplateFunc func() error) {
	var err = execTemplateFunc()
	if err != nil {
		handleError(w, http.StatusInternalServerError, errorMessage{Message: "Failed to instantiate template", Err: err})
	}
}
