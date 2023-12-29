package controllers

import (
	"encoding/json"
	"gofinance/services"
	"net/http"
)

func calculateTirWithGivenPrice(w http.ResponseWriter, r *http.Request) {

}

func GetTirs(w http.ResponseWriter, r *http.Request) {
	listOfTirResponse := services.GetTirs()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(listOfTirResponse)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}
