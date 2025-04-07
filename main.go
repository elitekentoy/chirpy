package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/asset/", http.FileServer(http.Dir("./asset/")))
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
