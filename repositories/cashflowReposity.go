package repositories

import (
	"context"
	"gofinance/conexion"
	"gofinance/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cashflowCollection *mongo.Collection = conexion.Client.Database("investment-project").Collection("cashflows")

func GetAllCashflowsWithTickets() ([]model.Cashflow, []string) {
	excludedBonds := []string{"Test 202312101819"}

	// Consultar todos los documentos en la colección
	cursor, err := cashflowCollection.Find(context.TODO(), bson.D{{Key: "ticket", Value: bson.D{{Key: "$nin", Value: excludedBonds}}}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var cashflows []model.Cashflow
	var tickets []string
	for cursor.Next(context.Background()) {
		var persona model.Cashflow
		err := cursor.Decode(&persona)
		if err != nil {
			log.Fatal(err)
		}
		cashflows = append(cashflows, persona)
		tickets = append(tickets, persona.Ticket)
	}

	return cashflows, tickets
}
