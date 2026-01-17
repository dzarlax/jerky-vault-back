package utils

import (
	"errors"
	"strconv"
	"strings"
)

// CalculateIngredientCost recalculates the ingredient price taking into account units of measurement
func CalculateIngredientCost(price float64, priceQuantity int, priceUnit string, recipeQuantityStr string, recipeUnit string) (float64, error) {
	recipeQuantity, err := strconv.ParseFloat(strings.Replace(recipeQuantityStr, ",", ".", 1), 64)
	if err != nil {
		return 0, errors.New("invalid recipe quantity")
	}

	var unitPrice float64

	switch {
	case priceUnit == "kg" && recipeUnit == "g":
		unitPrice = price / (float64(priceQuantity) * 1000) // price per gram
	case priceUnit == "g" && recipeUnit == "g":
		unitPrice = price / float64(priceQuantity) // price per gram
	case priceUnit == "l" && recipeUnit == "ml":
		unitPrice = price / (float64(priceQuantity) * 1000) // price per milliliter
	case priceUnit == "ml" && recipeUnit == "ml":
		unitPrice = price / float64(priceQuantity) // price per milliliter
	default:
		unitPrice = price / float64(priceQuantity) // price per unit
	}

	cost := unitPrice * recipeQuantity
	return cost, nil
}