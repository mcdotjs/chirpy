package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("chirpy go")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("listen", err)
	}
}
