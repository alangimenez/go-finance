package conexion

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func Find(filter interface{}, collection *mongo.Collection) *mongo.Cursor {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		// log.Fatal("Aca es el error")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	return cursor
}
