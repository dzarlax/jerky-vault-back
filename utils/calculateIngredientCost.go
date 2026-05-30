package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CalculateIngredientCost recalculates the ingredient price taking into account units of measurement
func CalculateIngredientCost(price float64, priceQuantity int, priceUnit string, recipeQuantityStr string, recipeUnit string) (float64, error) {
	recipeQuantity, err := strconv.ParseFloat(strings.Replace(recipeQuantityStr, ",", ".", 1), 64)
	if err != nil {
		return 0, errors.New("invalid recipe quantity")
	}
	if priceQuantity <= 0 {
		return 0, errors.New("price quantity must be greater than zero")
	}
	if recipeQuantity < 0 {
		return 0, errors.New("recipe quantity cannot be negative")
	}

	priceDimension, priceFactor, err := normalizeIngredientUnit(priceUnit)
	if err != nil {
		return 0, err
	}
	recipeDimension, recipeFactor, err := normalizeIngredientUnit(recipeUnit)
	if err != nil {
		return 0, err
	}
	if priceDimension != recipeDimension {
		return 0, fmt.Errorf("incompatible units: %s and %s", priceUnit, recipeUnit)
	}

	basePriceQuantity := float64(priceQuantity) * priceFactor
	baseRecipeQuantity := recipeQuantity * recipeFactor
	unitPrice := price / basePriceQuantity

	return unitPrice * baseRecipeQuantity, nil
}

func normalizeIngredientUnit(unit string) (string, float64, error) {
	switch strings.ToLower(strings.TrimSpace(unit)) {
	case "kg":
		return "mass", 1000, nil
	case "g":
		return "mass", 1, nil
	case "l":
		return "volume", 1000, nil
	case "ml":
		return "volume", 1, nil
	case "pcs":
		return "count", 1, nil
	default:
		return "", 0, fmt.Errorf("unsupported unit: %s", unit)
	}
}
