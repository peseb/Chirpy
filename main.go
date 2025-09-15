package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	fileHandler := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", http.StripPrefix("/app", fileHandler))
	serveMux.HandleFunc("/healthz", healthHandler)
	log.Fatal(server.ListenAndServe())
}

func healthHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}
