package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/makpoc/gopass/generator"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	formPage.Execute(w, nil)
}

func generatePageHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := parseForm(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, "Failed to parse form parameters")
		return
	}
	pass, err := generator.GeneratePassword(settings)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err, "Failed to generate your password")
		return
	}

	resultPage.Execute(w, string(pass))
}

func handleError(w http.ResponseWriter, status int, err error, msg string) {
	w.WriteHeader(status)
	errorPage.Execute(w, msg)
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	aboutPage.Execute(w, nil)
}

func parseForm(r *http.Request) (settings generator.Settings, err error) {
	settings = generator.DefaultSettings()

	masterPass := r.PostFormValue("password")
	confirmPass := r.PostFormValue("confirm-password")
	if masterPass != confirmPass {
		return settings, errors.New("Passwords differ")
	}
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

	sChars := r.PostFormValue("special-characters")
	if sChars != "" {
		var specialCharacters bool
		if specialCharacters, err = strconv.ParseBool(sChars); err != nil {
			return
		}
		settings.AddSpecialCharacters = specialCharacters
	}

	return
}