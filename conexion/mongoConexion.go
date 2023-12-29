package conexion

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EstablecerConexion() (*mongo.Client, error) {
	mongoURI, err := getURi()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getURi() (string, error) {
	uri, present := os.LookupEnv("MONGO_DB_URI")
	if !present {
		return "", errors.New("no está definido el URI de mongo")
	}
	return uri, nil
}

func CheckConection(client *mongo.Client) {
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión a MongoDB establecida.")
}
