package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("chirpy go")
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
