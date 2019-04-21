package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/rs/cors"
)

// Remove the prefix from the URL Path
func removePrefix(prefix string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		h(w, r)
	}
}

// Handle the "/" route
func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}

// Proxy a request to the given origin
func handleProxy(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}

// Handle GitHub requests by adding the GitHub auth token to the request
func handleGitHub(target string, authToken string) http.HandlerFunc {
	proxyHandler := handleProxy(target)
	return func(res http.ResponseWriter, req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+authToken)
		proxyHandler(res, req)
	}
}

// Run the proxy handlers
func main() {
	port := os.Getenv("PORT")
	ghURL := os.Getenv("GITHUB_API_URL")
	ghToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	log.Println("PORT=" + port)
	log.Println("GITHUB_API_URL=" + ghURL)
	log.Println("GITHUB_ACCESS_TOKEN=" + ghToken)
	log.Println("ALLOWED_ORIGINS=" + allowedOrigins)

	router := http.NewServeMux()
	router.HandleFunc("/github/", removePrefix("/github", handleGitHub(ghURL, ghToken)))
	router.HandleFunc("/", handleIndex())

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: strings.Split(allowedOrigins, ","),
	})
	http.ListenAndServe(":"+port, corsMiddleware.Handler(router))
}
