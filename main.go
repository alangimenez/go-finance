package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"
)

// Respuesta JSON para el GET
type GetResponse struct {
	Key   string  `json:"Key"`
	Value float64 `json:"Value"`
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
		Key:   "TIR",
		Value: 0.0,
	}

	fecha1 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	fecha2 := time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)

	diferencia := diferenciaEnDias(fecha1, fecha2)

	fmt.Printf("Diferencia en días entre %v y %v: %d días\n", fecha1, fecha2, diferencia)

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
	tir := calcularTIR(postData.Message) * 100

	// Procesar los datos
	fmt.Fprintf(w, "Handler POST en /handler. Mensaje recibido. Tir calculada: %0.2f%%\n", tir)
}

type PostData struct {
	Message []float64 `json:"message"`
}

func calcularNPV(tasaDescuento float64, cashFlow []float64) float64 {
	cashFlowWithoutInitialPayment := cashFlow[1:]
	var calculatedValues []float64
	for i, v := range cashFlowWithoutInitialPayment {
		calculatedValues = append(calculatedValues, (v / math.Pow(1+tasaDescuento, float64(i+1))))
	}
	var calculatedValuesWithInitialPayment []float64
	calculatedValuesWithInitialPayment = append([]float64{cashFlow[0]}, calculatedValues...)

	sumatoria := 0.0
	for _, valor := range calculatedValuesWithInitialPayment {
		sumatoria += valor
	}

	return sumatoria
}

func calcularTIRInterpolacion(tasaDescuentoInferior, tasaDescuentoSuperior float64, cashFlow []float64) float64 {
	npvInferior := calcularNPV(tasaDescuentoInferior, cashFlow)
	npvSuperior := calcularNPV(tasaDescuentoSuperior, cashFlow)

	// Verificar si ya estamos lo suficientemente cerca de la solución
	if math.Abs(npvInferior) < 0.0001 {
		return tasaDescuentoInferior
	}
	if math.Abs(npvSuperior) < 0.0001 {
		return tasaDescuentoSuperior
	}

	// Interpolación lineal
	tasaDescuentoInterpolada := tasaDescuentoInferior - (npvInferior*(tasaDescuentoSuperior-tasaDescuentoInferior))/(npvSuperior-npvInferior)

	return tasaDescuentoInterpolada
}

func calcularTIR(cashFlow []float64) float64 {
	// Definir tasas de descuento inicial y final
	tasaDescuentoInferior := 0.05
	tasaDescuentoSuperior := 0.1

	// Iterar hasta alcanzar la convergencia deseada
	for i := 0; i < 1000; i++ {
		tasaDescuentoInterpolada := calcularTIRInterpolacion(tasaDescuentoInferior, tasaDescuentoSuperior, cashFlow)

		npvInterpolado := calcularNPV(tasaDescuentoInterpolada, cashFlow)

		// Actualizar los límites para la próxima iteración
		if npvInterpolado < 0 {
			tasaDescuentoInferior = tasaDescuentoInterpolada
		} else {
			tasaDescuentoSuperior = tasaDescuentoInterpolada
		}

		// Verificar convergencia
		if math.Abs(npvInterpolado) < 0.0001 {
			return tasaDescuentoInterpolada
		}
	}

	// Manejar el caso en el que no converge
	panic("No se pudo converger a una TIR en el número máximo de iteraciones")
}

func diferenciaEnDias(fecha1, fecha2 time.Time) int {
	// Truncar las fechas para ignorar la información de la hora
	fecha1 = fecha1.Truncate(24 * time.Hour)
	fecha2 = fecha2.Truncate(24 * time.Hour)

	// Calcular la diferencia en días
	diferencia := fecha2.Sub(fecha1) / (24 * time.Hour)

	// Convertir la diferencia a un entero
	return int(diferencia)

	// Ejemplo de uso
	/* fecha1 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	fecha2 := time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)

	diferencia := diferenciaEnDias(fecha1, fecha2)

	fmt.Printf("Diferencia en días entre %v y %v: %d días\n", fecha1, fecha2, diferencia) */
}
