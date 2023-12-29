package model

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
