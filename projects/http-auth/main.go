package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	// "fmt"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func ParamsList(u *url.URL) string {
	values := u.Query()
	var list string
	for key, value := range values {
		list += "<li>" + html.EscapeString(key) + ": [" + html.EscapeString(strings.Join(value, ", ")) + "] </li>"
	}

	return "<ul>" + list + "</ul>"
}

var limiter = rate.NewLimiter(100, 30)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data string
		if r.Method == "POST" {

			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)

			if err != nil {
				data = "Error reading data"
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				data = html.EscapeString(string(bodyBytes)) + ParamsList(r.URL)
			}

		} else {
			data = "<h1>Hello Tomi</h1>"
		}
		w.Header().Add("content-type", "text/html")
		w.Write([]byte(data))
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200"))
	})

	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Internal server error"))
	})

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		expectedUsername := os.Getenv("AUTH_USERNAME")
		expectedPassword := os.Getenv("AUTH_PASSWORD")

		if !ok || username != expectedUsername || password != expectedPassword {
			w.Header().Add("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Add("content-type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte("<h1>Hello " + username))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusServiceUnavailable)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<!DOCTYPE html>\n<html>\nHello world!"))
	})

	http.ListenAndServe(":8080", nil)
}
