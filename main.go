package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

func main() {
	log := log.New(os.Stdout, "[werds] ", log.LstdFlags)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "werds-241615")
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()
	werds := client.Collection("werds")

	mainTmpl, err := template.ParseFiles("views/main.tmpl.html")
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()

	// Werds display
	router.
		Methods(http.MethodGet).
		Path("/{key}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, ok := mux.Vars(r)["key"]
			if !ok {
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			doc, err := werds.Doc(key).Get(ctx)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusPermanentRedirect)
				return
			}
			text := doc.Data()["text"]

			mainTmpl.ExecuteTemplate(w, "main.tmpl.html", text)
		})

	// Index
	router.
		Methods(http.MethodGet).
		Path("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			text := r.URL.Query().Get("t")

			// Serve index view
			if text == "" {
				http.ServeFile(w, r, "views/index.html")
				return
			}

			key, err := generateKey(8)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			if len(text) > 1000 {
				text = text[:1000]
			}
			if _, err := werds.Doc(key).Set(ctx, map[string]interface{}{"text": text}); err != nil {
				log.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/%s", key), http.StatusPermanentRedirect)
		})

	// Fileserver
	router.
		Methods(http.MethodGet).
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("public")))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/", http.StatusPermanentRedirect) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("serving @ %s", addr)
	log.Panic(http.ListenAndServe(addr, newLogger(log)(router)))
}

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggerResponseWriter(w http.ResponseWriter) *loggerResponseWriter {
	return &loggerResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggerResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func newLogger(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggerResponseWriter := newLoggerResponseWriter(w)

			start := time.Now()
			next.ServeHTTP(loggerResponseWriter, r)
			elapsed := time.Now().Sub(start)

			log.Printf("%s %s | %s | %d\n", r.Method, r.URL.String(), elapsed.String(), loggerResponseWriter.statusCode)
		})
	}
}

// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func generateKey(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
