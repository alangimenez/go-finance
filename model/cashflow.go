package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cashflow struct {
	Company            string
	Start              time.Time
	Finish             time.Time
	Rate               primitive.Decimal128 `bson:"Rate"`
	DateOfPayment      []time.Time
	AmountInterest     []float64
	AmountAmortization []float64
	Ticket             string
}
