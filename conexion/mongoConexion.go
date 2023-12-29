package conexion

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EstablecerConexion() *mongo.Client {
	mongoURI := getURi()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, errConnect := mongo.Connect(context.Background(), clientOptions)
	if errConnect != nil {
		log.Fatal(errConnect)
	}

	errPing := client.Ping(context.Background(), nil)
	if errPing != nil {
		log.Fatal(errPing)
	}
	fmt.Println("Conexión a MongoDB establecida.")

	return client
}

var Client *mongo.Client = EstablecerConexion()

func getURi() string {
	// Load environment configs
	errEnv := godotenv.Load("config.env")
	if errEnv != nil {
		fmt.Println("Error cargando el archivo de configuración:", errEnv)
	}
	return os.Getenv("MONGO_DB_URI")
}
