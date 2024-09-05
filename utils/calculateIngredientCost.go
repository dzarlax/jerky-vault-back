package utils

import (
	"errors"
	"strconv"
	"strings"
)

// CalculateIngredientCost пересчитывает цену ингредиента с учетом единиц измерения
func CalculateIngredientCost(price float64, priceQuantity int, priceUnit string, recipeQuantityStr string, recipeUnit string) (float64, error) {
	recipeQuantity, err := strconv.ParseFloat(strings.Replace(recipeQuantityStr, ",", ".", 1), 64)
	if err != nil {
		return 0, errors.New("invalid recipe quantity")
	}

	var unitPrice float64

	switch {
	case priceUnit == "kg" && recipeUnit == "g":
		unitPrice = price / (float64(priceQuantity) * 1000) // цена за грамм
	case priceUnit == "g" && recipeUnit == "g":
		unitPrice = price / float64(priceQuantity) // цена за грамм
	case priceUnit == "l" && recipeUnit == "ml":
		unitPrice = price / (float64(priceQuantity) * 1000) // цена за миллилитр
	case priceUnit == "ml" && recipeUnit == "ml":
		unitPrice = price / float64(priceQuantity) // цена за миллилитр
	default:
		unitPrice = price / float64(priceQuantity) // цена за единицу
	}

	cost := unitPrice * recipeQuantity
	return cost, nil
}