package main

import (
	"log"
	"net/http"
	"time"
)

// loggingResponseWriter embeds the default responseWriter and provides a status field to get the status of a request
// after the response is received
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = 200
	}
	return l.ResponseWriter.Write(b)
}

// logger wrapps a HandlerFunc and logs a request/response information
func logger(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		loggingRW := &loggingResponseWriter{ResponseWriter: w}
		inner.ServeHTTP(loggingRW, r)

		log.Printf("%s:\t%s\t%s\t%d\t%s\n", time.Now().Format(time.RFC822Z), r.Method, r.RequestURI, loggingRW.status, time.Since(start))
	})
}
