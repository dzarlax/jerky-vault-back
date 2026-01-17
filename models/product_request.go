// Ensure the package name is specified at the beginning of the file
package models

// Definition of product request structure
type ProductRequest struct {
	Product Product         `json:"product"`
	Options []ProductOption `json:"options"`
}
