package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	port := "8080"
	filepathRoot := "."
	
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//log.Fatal(srv.ListenAndServe())
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("listen", err)
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
}
