package services

import (
	"fmt"
	"gofinance/model"
	"gofinance/repositories"
	"gofinance/responses"
	"math"
	"time"
)

func GetTirs() []responses.TirResponse {
	cashflows, tickets := repositories.GetAllCashflowsWithTickets()
	quotesList := repositories.GetAllLastvalues()

	fmt.Print(tickets)

	var listOfTirResponse []responses.TirResponse

	for _, bond := range cashflows {
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

		response := responses.TirResponse{
			Key:   bond.Ticket,
			Value: tirAnual,
			Price: actualPrice,
		}
		listOfTirResponse = append(listOfTirResponse, response)
	}

	return listOfTirResponse
}

func CalculateTirWithGivenPrice(price float64, ticket string) (float64, error) {
	cashflow, err := repositories.GetCashflowByTicket(ticket)
	if err != nil {
		return 0.0, err
	}

	array := createArray(cashflow.Finish)
	secondArray := addPaymentsToArray(
		cashflow.DateOfPayment,
		cashflow.AmountInterest,
		cashflow.AmountAmortization,
		array,
		price,
	)
	tir := calculoTirByInterpolation(secondArray)
	tirAnual := tasaEfectivaAnual(tir)
	return tirAnual, nil
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

func diferenciaEnDias(fecha1, fecha2 time.Time) int {
	// Truncar las fechas para ignorar la información de la hora
	fecha1 = fecha1.Truncate(24 * time.Hour)
	fecha2 = fecha2.Truncate(24 * time.Hour)

	// Calcular la diferencia en días
	diferencia := fecha2.Sub(fecha1) / (24 * time.Hour)

	// Convertir la diferencia a un entero
	return int(diferencia)
}

func getActualPrice(quotesList []model.Bond, ticket string, mep float64) (float64, error) {
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

func calcularNPV(tasaDescuento float64, cashFlow []float64) float64 {
	cashFlowWithoutInitialPayment := cashFlow[1:]
	var calculatedValues []float64
	for i, v := range cashFlowWithoutInitialPayment {
		calculatedValues = append(calculatedValues, (v / math.Pow(1+tasaDescuento, float64(i+1))))
	}
	var calculatedValuesWithInitialPayment = append([]float64{cashFlow[0]}, calculatedValues...)

	sumatoria := 0.0
	for _, valor := range calculatedValuesWithInitialPayment {
		sumatoria += valor
	}

	return sumatoria
}

func tasaEfectivaAnual(tasaEfectivaDiaria float64) float64 {
	// Convertir la tasa efectiva diaria a anual
	tea := math.Pow(1+tasaEfectivaDiaria, 365) - 1
	return tea
}

func interpolation(rate, npvPositive, npvNegative float64) float64 {
	previousRate := rate - 0.000001
	return previousRate + ((rate - previousRate) * (npvPositive / (npvPositive - npvNegative)))
}
