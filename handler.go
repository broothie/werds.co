package main

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const keyLength = 6

var (
	keyRegexp  = regexp.MustCompile(fmt.Sprintf(`^/[a-zA-Z]{%d}$`, keyLength))
	mainTmpl   = template.Must(template.ParseFiles("views/main.tmpl.html"))
	fileServer = http.FileServer(http.Dir("public"))
	alphabet   = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var werds *firestore.CollectionRef

func init() {
	rand.Seed(time.Now().UnixNano())

	client, err := firestore.NewClient(context.Background(), "werds-241615")
	if err != nil {
		logger.Panic(err)
	}

	werds = client.Collection("werds")
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Write post
	text := r.URL.Query().Get("t")
	if text != "" {
		// Truncate posts to 1000 characters
		if len(text) > 1000 {
			text = text[:1000]
		}

		// Generate key and add to db
		key := generateKey()
		if _, err := werds.Doc(key).Set(ctx, map[string]interface{}{"text": text}); err != nil {
			logger.Println(err)
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}

		// Redirect to post
		http.Redirect(w, r, fmt.Sprintf("/%s", key), http.StatusPermanentRedirect)
		return
	}

	// Read post
	if keyRegexp.MatchString(r.URL.Path) {
		// Get key and pull post out of db
		key := strings.ToUpper(strings.Split(r.URL.Path, "/")[1])
		doc, err := werds.Doc(key).Get(ctx)
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

		return
	}

	// File server
	notFoundInterceptor := newInterceptor(w, func(i *Interceptor) {
		if i.statusCode == http.StatusNotFound {
			http.Redirect(i, r, "/", http.StatusPermanentRedirect)
		}
	})
	fileServer.ServeHTTP(notFoundInterceptor, r)
}

func generateKey() string {
	runes := make([]rune, keyLength)
	for i := range runes {
		runes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(runes)
}
