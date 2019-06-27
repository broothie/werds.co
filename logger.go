package main

import (
	"log"
	"net/http"
	"time"
)

func newLoggerMiddleware(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			interceptor := newInterceptor(w, nil)

			start := time.Now()
			next.ServeHTTP(interceptor, r)
			elapsed := time.Now().Sub(start)

			log.Printf("%s %s | %s | %d\n", r.Method, r.URL.String(), elapsed.String(), interceptor.statusCode)
		})
	}
}
