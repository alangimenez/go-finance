package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gofinance/conexion"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Respuesta JSON para el GET
type GetResponse struct {
	Key   string  `json:"Key"`
	Value float64 `json:"Value"`
	Price float64 `json:"Price"`
}

type Bond struct {
	Simbolo      string
	UltimoPrecio float64
	Moneda       string
}

type Quotes struct {
	Quotes      []Bond    `bson:"quotes"`
	Id          string    `bson:"_id"`
	OtherQuotes PreDolars `bson:"otherQuotes"`
}

type PreDolars struct {
	Quotes Dolars `bson:"quotes"`
}

type Dolars struct {
	Mep float64 `bson:"dolarMep"`
}

func main() {
	errEnv := godotenv.Load("config.env")
	if errEnv != nil {
		fmt.Println("Error cargando el archivo de configuración:", errEnv)
		return
	}

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

	// Configurar los manejadores para los endpoints
	http.HandleFunc("/handler", handler)
	http.HandleFunc("/tir/price", calculateTirWithGivenPrice)

	// Iniciar el servidor en el puerto 5050
	fmt.Printf("Servidor escuchando en el puerto %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsHandler)
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
	/* response := GetResponse{
		Key:   "TIR",
		Value: 0.0,
	} */

	responsesList := connectToMongo()

	// Codificar la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(responsesList)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}

func calculateTirWithGivenPrice(w http.ResponseWriter, r *http.Request) {
	response := "hola"

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
	// connectToMongo()

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
	tasaDescuentoInferior := 0.0711
	tasaDescuentoSuperior := 0.2

	// Iterar hasta alcanzar la convergencia deseada
	for i := 0; i < 1; i++ {
		tasaDescuentoInterpolada := calcularTIRInterpolacion(tasaDescuentoInferior, tasaDescuentoSuperior, cashFlow)

		npvInterpolado := calcularNPV(tasaDescuentoInterpolada, cashFlow)

		// Actualizar los límites para la próxima iteración
		if npvInterpolado < 0 {
			tasaDescuentoInferior = tasaDescuentoInterpolada
		} else {
			tasaDescuentoSuperior = tasaDescuentoInterpolada
		}

		// Verificar convergencia
		if math.Abs(npvInterpolado) < 5 {
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

func connectToMongo() []GetResponse {
	/* // Establecer información de conexión
	uri, present := os.LookupEnv("MONGO_DB_URI")
	if !present {
		fmt.Printf("No esta definido el URI de Mongo.")
		return []GetResponse{}
	}
	clientOptions := options.Client().ApplyURI(uri)

	// Crear un cliente de MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Comprobar la conexión
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión a MongoDB establecida.") */

	client, err := conexion.EstablecerConexion()
	if err != nil {
		fmt.Errorf("La conexion no pudo ser establecida")
	}
	conexion.CheckConection(client)

	// Obtener una referencia a la colección
	collection := client.Database("investment-project").Collection("cashflows")

	excludedBonds := []string{"Test 202312101819"}

	// Consultar todos los documentos en la colección
	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "ticket", Value: bson.D{{Key: "$nin", Value: excludedBonds}}}})
	if err != nil {
		// log.Fatal("Aca es el error")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Recorrer los documentos y mostrarlos
	var personas []Account
	var tickets []string
	for cursor.Next(context.Background()) {
		var persona Account
		err := cursor.Decode(&persona)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(time.Parse("2006/01/02", persona.Start))
		personas = append(personas, persona)
		tickets = append(tickets, persona.Ticket)
	}
	fmt.Println("Documentos en la colección:")
	fmt.Println(personas)

	collectionValues := client.Database("investment-project").Collection("lastvalues")
	filter := bson.D{}

	result := collectionValues.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	var quotesList Quotes
	err = result.Decode(&quotesList)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Quotes: ", quotesList)

	var responses []GetResponse

	for _, bond := range personas {
		fmt.Printf("Este es el bono %s \n", bond.Ticket)
		array := createArray(bond.Finish)
		fmt.Printf("Este es el bono %s after createArray. Length of array: %d \n", bond.Ticket, len(array))
		actualPrice, errActualPrice := getActualPrice(quotesList.Quotes, bond.Ticket, quotesList.OtherQuotes.Quotes.Mep)
		if errActualPrice != nil {
			fmt.Printf("Error obteniendo el precio actual de %s \n", bond.Ticket)
		}
		if actualPrice == 0.0 {
			fmt.Printf("El precio del bono %s es 0", bond.Ticket)
			continue
		}
		fmt.Printf("Para bono %s precio %f", bond.Ticket, actualPrice)

		arrayTwo := addPaymentsToArray(
			bond.DateOfPayment,
			bond.AmountInterest,
			bond.AmountAmortization,
			array,
			actualPrice,
		)
		fmt.Printf("Este es el bono %s after addPaymentsToArray \n", bond.Ticket)
		if arrayTwo == nil {
			fmt.Printf("hola")
		}
		for pos, value := range arrayTwo {
			if value != 0 {
				fmt.Printf("Pos %d, value %f \n", pos, value)
			}
		}
		/* tir := calcularTIR(arrayTwo)
		tirAnual := tasaEfectivaAnual(tir / 100)
		fmt.Printf("La tir es de %0.2f", tirAnual) */

		// cashflow := []float64{-10000.0, 5000.0, 9000.0}
		tir := calculoTirByInterpolation(arrayTwo)
		tirAnual := tasaEfectivaAnual(tir)
		fmt.Printf("La tirAnual es de %f", tirAnual)

		response := GetResponse{
			Key:   bond.Ticket,
			Value: tirAnual,
			Price: actualPrice,
		}
		responses = append(responses, response)
	}

	// Desconectar el cliente
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión a MongoDB cerrada.")
	return responses
}

type Account struct {
	Company            string
	Start              time.Time
	Finish             time.Time
	Rate               primitive.Decimal128 `bson:"Rate"`
	DateOfPayment      []time.Time
	AmountInterest     []float64
	AmountAmortization []float64
	Ticket             string
}

func createArray(endDate time.Time) []float64 {
	fmt.Printf("The endDate is %s", endDate)
	// parsedEndTime := parseDate(endDate)
	length := diferenciaEnDias(time.Now(), endDate)
	fmt.Printf("The length is %d", length)
	var array []float64
	for i := 0; i < length; i++ {
		array = append(array, 0)
	}
	return array
}

func addPaymentsToArray(paymentDays []time.Time, paymentInterest, paymentAmortization, array []float64, actualPrice float64) []float64 {
	for i := 0; i < len(paymentDays); i++ {
		// parsedEndTime := parseDate(paymentDays[i])
		diffDays := diferenciaEnDias(time.Now(), paymentDays[i])
		if diffDays > 0 {
			array[diffDays-1] = paymentInterest[i] + paymentAmortization[i]
		}
	}
	array[0] = -actualPrice
	return array
}

func parseDate(date string) time.Time {
	parsedEndTime, err := time.Parse("2006/01/02", date)
	if err != nil {
		fmt.Printf("Error parseando %s a time.Time", date)
	}
	return parsedEndTime
}

func tasaEfectivaAnual(tasaEfectivaDiaria float64) float64 {
	// Convertir la tasa efectiva diaria a anual
	tea := math.Pow(1+tasaEfectivaDiaria, 365) - 1
	return tea
}

func calculoTirByInterpolation(cashflow []float64) float64 {
	rate := 0.000001
	var tir float64
	for {
		previousNpv := calcularNPV(rate-0.000001, cashflow)
		actualNPV := calcularNPV(rate, cashflow)
		previousNpvNegative := calcularNPV(-rate+0.000001, cashflow)
		actualNPVNegative := calcularNPV(-rate, cashflow)
		if previousNpv >= 0.0 && actualNPV < 0.0 {
			fmt.Printf("npvPositivo %f \n", previousNpv)
			fmt.Printf("npvPositivo %f \n", actualNPV)
			fmt.Printf("tasa previa %f", rate-0.000001)
			fmt.Printf("tasa actual %f", rate)
			tir = interpolation(rate, previousNpv, actualNPV)
			break
		}
		if previousNpvNegative <= 0.0 && actualNPVNegative > 0.0 {
			fmt.Printf("npvPositivo %f \n", previousNpvNegative)
			fmt.Printf("npvPositivo %f \n", actualNPVNegative)
			fmt.Printf("tasa previa %f", -rate+0.000001)
			fmt.Printf("tasa actual %f", -rate)
			tir = interpolation(-rate, previousNpvNegative, actualNPVNegative)
			break
		}
		rate += 0.000001
	}
	return tir
}

func interpolation(rate, npvPositive, npvNegative float64) float64 {
	previousRate := rate - 0.000001
	return previousRate + ((rate - previousRate) * (npvPositive / (npvPositive - npvNegative)))
}

func getActualPrice(quotesList []Bond, ticket string, mep float64) (float64, error) {
	fmt.Printf("dolar mep %f", mep)
	for _, p := range quotesList {
		if p.Simbolo == ticket && (p.Moneda == "1" || p.Moneda == "AR$") {
			return p.UltimoPrecio / mep, nil
		}
		if p.Simbolo == ticket {
			return p.UltimoPrecio, nil
		}
	}
	return 0.0, fmt.Errorf("No se encontró el ticket con el nombre %s", ticket)
}
