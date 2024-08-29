package repository

type Products struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}
type UpdateProducts struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"required"`
	Quantity    *int     `json:"quantity" binding:"required"`
}
