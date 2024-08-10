package services

import (
	"gofinance/model"
	"gofinance/repositories"
)

func GetQuotesByAssetType(assetType string) []model.Bond {
	quotesList := repositories.GetAllLastvalues()
	filteredQuotesList := filterByAssetType(quotesList.Quotes, assetType)
	return filteredQuotesList
}

func filterByAssetType(quotesList []model.Bond, assetType string) []model.Bond {
	var filteredAssets []model.Bond

	// Iterar sobre el listado de assets
	for _, asset := range quotesList {
		// Si el assetType coincide, añadirlo al listado filtrado
		if asset.Type == assetType {
			filteredAssets = append(filteredAssets, asset)
		}
	}

	return filteredAssets
}
