package controllers

import (
	"net/url"
	"strconv"
)

func GetPrice(queryParams url.Values) float64 {
	stringPrice := queryParams.Get("price")
	price, err := strconv.ParseFloat(stringPrice, 64)
	if err != nil {
		return 0.0
	}
	return price
}

func GetTicket(queryParams url.Values) string {
	return queryParams.Get("ticket")
}
