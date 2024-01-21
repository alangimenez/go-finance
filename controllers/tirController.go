package controllers

import (
	"encoding/json"
	"gofinance/services"
	"net/http"
)

func GetTirs(w http.ResponseWriter, r *http.Request) {
	listOfTirResponse := services.GetTirs()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(listOfTirResponse)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}

func CalculateTirWithGivenPrice(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	price := GetPrice(queryParams)
	if price == 0 {
		http.Error(w, "price must be a float64 type", http.StatusInternalServerError)
		return
	}
	ticket := GetTicket(queryParams)
	if ticket == "" {
		http.Error(w, "ticket must not be empty", http.StatusInternalServerError)
		return
	}

	tir, err := services.CalculateTirWithGivenPrice(price, ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(tir)
	if errEncode != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}
