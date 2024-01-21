package routes

import (
	"gofinance/controllers"
	"net/http"
)

func TirRoutes() {
	http.HandleFunc("/tir", tir)
	http.HandleFunc("/tir/price", calculateTirWithGivenPrice)
}

func tir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		controllers.GetTirs(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func calculateTirWithGivenPrice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		controllers.CalculateTirWithGivenPrice(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
