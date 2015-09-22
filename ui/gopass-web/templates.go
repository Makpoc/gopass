package main

import (
	"html/template"
	"path"
)

var templatesDir = "templates/"

var templatePage, formPage, resultPage, aboutPage, errorPage *template.Template

func prepareTemplates() {
	templatePage = template.Must(
		template.ParseFiles(
			path.Join(templatesDir, "base.html"),
			path.Join(templatesDir, "headers.html"),
			path.Join(templatesDir, "nav.html")))

	formPage = prepareTemplate("input-form.html")
	resultPage = prepareTemplate("result.html")
	aboutPage = prepareTemplate("about.html")
	errorPage = prepareTemplate("error.html")
}

func prepareTemplate(file string) *template.Template {
	page := template.Must(templatePage.Clone())
	return template.Must(page.ParseFiles(path.Join(templatesDir, file)))
}
