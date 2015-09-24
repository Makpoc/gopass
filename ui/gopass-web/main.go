package main

import (
	"log"
	"net/http"
	"os"
)

const defaultPort = "8000"
const portEnvKey = "PORT"

func getPort() string {
	port := os.Getenv(portEnvKey)
	if port == "" {
		port = defaultPort
	}

	return port
}

func main() {
	prepareTemplates()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", logger(homeHandler))
	http.HandleFunc("/about", logger(aboutHandler))
	http.HandleFunc("/generate_pass", logger(generateHandler))
	//	http.HandleFunc("/", logger(generate, "Generate"))
	err := http.ListenAndServe(":"+getPort(), nil)
	if err != nil {
		log.Fatal(err)
	}

}
