package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

const keyLength = 6

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logger := log.New(os.Stdout, "[werds] ", log.LstdFlags)

	client, err := firestore.NewClient(context.Background(), "werds-241615")
	if err != nil {
		logger.Panic(err)
	}
	werds := client.Collection("werds")

	router := mux.NewRouter()

	// Submit post
	router.
		Methods(http.MethodGet).
		Queries("t", "{t}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Truncate posts to 1000 characters
			text := mux.Vars(r)["t"]
			if len(text) > 1000 {
				text = text[:1000]
			}

			// Generate key and add to db
			key := generateKey()
			if _, err := werds.Doc(key).Set(context.Background(), map[string]interface{}{"text": text}); err != nil {
				logger.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			// Redirect to post
			http.Redirect(w, r, fmt.Sprintf("/%s", key), http.StatusPermanentRedirect)
		})

	// Show post
	mainTmpl := template.Must(template.ParseFiles("views/main.tmpl.html"))
	router.
		Methods(http.MethodGet).
		Path("/{key:[0-9a-zA-Z=]{4,}}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get key and pull post out of db
			key := mux.Vars(r)["key"]
			doc, err := werds.Doc(key).Get(context.Background())
			if err != nil {
				logger.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}

			// Render view
			text := doc.Data()["text"]
			if err := mainTmpl.ExecuteTemplate(w, "main.tmpl.html", text); err != nil {
				logger.Println(err)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			}
		})

	// File server at /public/
	router.
		PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Index
	router.
		Methods(http.MethodGet).
		Path("/").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/index.html") })

	router.NotFoundHandler = http.RedirectHandler("/", http.StatusPermanentRedirect)

	loggerMiddleware := newLoggerMiddleware(logger)
	handler := loggerMiddleware(router)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Printf("serving @ %s", addr)
	logger.Panic(http.ListenAndServe(addr, handler))
}

var alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateKey() string {
	runes := make([]rune, keyLength)
	for i := range runes {
		runes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(runes)
}

func newLoggerMiddleware(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			statusRecorder := statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			start := time.Now()
			next.ServeHTTP(&statusRecorder, r)
			elapsed := time.Now().Sub(start)

			log.Printf("%s %s | %s | %d\n", r.Method, r.URL.String(), elapsed.String(), statusRecorder.statusCode)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	sr.statusCode = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}
