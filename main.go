package main

import (
	"fmt"
	"gofinance/conexion"
	"gofinance/routes"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	// Set connection to database
	conexion.EstablecerConexion()

	// Configurar CORS middleware con opciones específicas
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                              // Origen permitido (ajusta según tu necesidad)
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "HEAD"}, // Métodos permitidos
		AllowedHeaders:   []string{"Authorization", "Content-Type"},  // Encabezados permitidos
		AllowCredentials: true,                                       // Permitir credenciales (cookies, autenticación, etc.)
		Debug:            true,                                       // Habilitar modo de depuración
	})

	// Configurar CORS middleware
	corsHandler := corsOptions.Handler(http.DefaultServeMux)

	// Configurar puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routes.TirRoutes()

	// Iniciar el servidor en el puerto 5050
	fmt.Printf("Servidor escuchando en el puerto %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
