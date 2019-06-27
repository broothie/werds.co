package main

import (
	"log"
	"net/http"
	"time"
)

func newLogger(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggerResponseWriter := newResponseWriter(w)

			start := time.Now()
			next.ServeHTTP(loggerResponseWriter, r)
			elapsed := time.Now().Sub(start)

			log.Printf("%s %s | %s | %d\n", r.Method, r.URL.String(), elapsed.String(), loggerResponseWriter.statusCode)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *responseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
