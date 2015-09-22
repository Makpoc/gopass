package main

import (
	"log"
	"net/http"
)

func main() {
	prepareTemplates()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", logger(homeHandler))
	http.HandleFunc("/about", logger(aboutHandler))
	http.HandleFunc("/generate_pass", logger(generateHandler))
	//	http.HandleFunc("/", logger(generate, "Generate"))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
