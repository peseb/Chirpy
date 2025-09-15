package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	fileHandler := http.FileServer(http.Dir("."))
	serveMux.Handle("/", fileHandler)
	server.ListenAndServe()
}
