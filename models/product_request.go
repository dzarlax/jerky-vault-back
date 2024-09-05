// Убедитесь, что в начале файла указано название пакета
package models

// Определение структуры запроса для продукта
type ProductRequest struct {
	Product Product         `json:"product"`
	Options []ProductOption `json:"options"`
}
