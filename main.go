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
func removePrefix(prefix string, handle http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		handle(res, req)
	}
}

// Handles the "/" route
func handleIndex() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	}
}

// Proxys a request to the given origin
func handleProxy(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(res http.ResponseWriter, req *http.Request) {
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host
		proxy.ServeHTTP(res, req)
	}
}

// Handles GitHub requests by adding the GitHub auth token to the request
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
