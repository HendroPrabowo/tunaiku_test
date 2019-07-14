package main

import (
	"encoding/json"
	"net/http"
	"fmt"
)

type Message struct {
	Message string `json:"message"`
}

func RecoverFunc(w http.ResponseWriter) {
	if r := recover(); r != nil {
		text := fmt.Sprint(r)
		CustomResponse(500, text, w)
	}
}

func CustomResponse(response_code int, message string, w http.ResponseWriter){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response_code)
	json.NewEncoder(w).Encode(Message{message})
}
