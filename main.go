package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"cloud.google.com/go/firestore"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log := log.New(os.Stdout, "[werds] ", log.LstdFlags)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "werds-241615")
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	werds := client.Collection("werds")
	keyRegexp := regexp.MustCompile(`^/[a-z]{4}$`)
	mainTmpl := template.Must(template.ParseFiles("views/main.tmpl.html"))
	fileServer := http.FileServer(http.Dir("public"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("t")
		if text != "" {
			if len(text) > 1000 {
				text = text[:1000]
			}

			key := generateKey()
			if _, err := werds.Doc(key).Set(ctx, map[string]interface{}{"text": text}); err != nil {
				log.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/%s", key), http.StatusPermanentRedirect)
			return
		}

		if keyRegexp.MatchString(r.URL.Path) {
			key := r.URL.Path[1:]
			doc, err := werds.Doc(key).Get(ctx)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			mainTmpl.ExecuteTemplate(w, "main.tmpl.html", doc.Data()["text"])
			return
		}

		fileServer.ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("serving @ %s", addr)
	log.Panic(http.ListenAndServe(addr, newLogger(log)(handler)))
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

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func generateKey() string {
	alphabet := []rune(alphabet)

	runes := make([]rune, 4)
	for i := range runes {
		runes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(runes)
}
