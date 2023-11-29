package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Respuesta JSON para el GET
type GetResponse struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func main() {
	// Configurar puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configurar los manejadores para los endpoints
	http.HandleFunc("/handler", handler)

	// Iniciar el servidor en el puerto 5050
	fmt.Printf("Servidor escuchando en el puerto %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Manejar la solicitud GET
	// Crear la respuesta JSON
	response := GetResponse{
		Key:   "Key",
		Value: "Hola Alan",
	}

	// Codificar la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Manejar la solicitud POST
	// Decodificar el cuerpo JSON
	var postData PostData
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Procesar los datos
	fmt.Fprintf(w, "Handler POST en /handler. Mensaje recibido: %s\n", postData.Message)
}

type PostData struct {
	Message string `json:"message"`
}
