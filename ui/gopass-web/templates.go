package main

import (
	"html/template"
	"path"
)

var templatesDir = "templates/"

var templatePage, formPage, resultPage, aboutPage, errorPage *template.Template

// prepareTemplates ensures all pages have a valid template and parses it for each of them
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

// prepareTemplate clones the main structure and parses the template for the concrete page
func prepareTemplate(file string) *template.Template {
	page := template.Must(templatePage.Clone())
	return template.Must(page.ParseFiles(path.Join(templatesDir, file)))
}
