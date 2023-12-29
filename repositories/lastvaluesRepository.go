package repositories

import (
	"context"
	"gofinance/conexion"
	"gofinance/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var lastvaluesCollection *mongo.Collection = conexion.Client.Database("investment-project").Collection("lastvalues")

func GetAllLastvalues() model.Quotes {
	filter := bson.D{}

	result := lastvaluesCollection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	var quotesList model.Quotes
	err := result.Decode(&quotesList)
	if err != nil {
		log.Fatal(err)
	}

	return quotesList
}
