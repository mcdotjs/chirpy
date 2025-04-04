package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("c")
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	type returnVals struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}

	statusCode := 200
	respBody := returnVals{}

	bodyLength := len(params.Body)
	if bodyLength > 140 {
		respBody.Error = "Chirp is too long"
		statusCode = 400
	} else {
		respBody.Valid = true
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
