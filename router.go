package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequest()  {
	// Router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CustomResponse(200, "Tunaiku Test API", w)
	}).Methods("GET")

	// Testing
	myRouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		CustomResponse(200, "success connect to api", w)
	})

	// Loan
	myRouter.HandleFunc("/loans", AllLoan).Methods("GET")
	myRouter.HandleFunc("/loan/create", NewLoan).Methods("POST")
	myRouter.HandleFunc("/loan/list/{name}", ListLoan).Methods("POST")
	myRouter.HandleFunc("/loan/installment", InstallmentLoan).Methods("POST")

	// Run Port
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}