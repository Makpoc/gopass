package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = "8000"
	localhost   = "localhost" // public
	publichost  = ""

	portEnvKey = "GOPASS_PORT"
	hostEnvKey = "GOPASS_LOCALHOST"
)

// getPort checks if the PORT env variable is set and uses it instead of the defaultPort
func getPort() string {
	port := os.Getenv(portEnvKey)
	if port == "" {
		port = defaultPort
	}

	return port
}

func getHost() string {
	_, isSet := os.LookupEnv(hostEnvKey)
	if isSet {
		return localhost
	}
	return publichost
}

func main() {
	prepareTemplates()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", logger(homePageHandler))
	http.HandleFunc("/about", logger(aboutPageHandler))
	http.HandleFunc("/generate_pass", logger(generatePageHandler))

	port := getHost() + ":" + getPort()
	fmt.Println("Listening on " + port + "...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
