package domain

// MeatCategory represents the type of meat.
type MeatCategory string

const (
	CategoryBoeuf       MeatCategory = "boeuf"
	CategoryAgneau      MeatCategory = "agneau"
	CategoryPoulet      MeatCategory = "poulet"
	CategoryVeau        MeatCategory = "veau"
	CategoryCharcuterie MeatCategory = "charcuterie"
)

// ValidCategories lists all valid meat categories.
var ValidCategories = []MeatCategory{
	CategoryBoeuf, CategoryAgneau, CategoryPoulet, CategoryVeau, CategoryCharcuterie,
}

// Product represents a meat product sold by the butcher.
type Product struct {
	ID         string       `json:"id"`
	Name       string       `json:"name" validate:"required,min=2"`
	Category   MeatCategory `json:"category" validate:"required"`
	PricePerKg float64      `json:"pricePerKg" validate:"required,gt=0"`
	Image      string       `json:"image,omitempty"`
	InStock    bool         `json:"inStock"`
}

// CreateProductRequest represents the payload to create a product.
type CreateProductRequest struct {
	Name       string       `json:"name" validate:"required,min=2"`
	Category   MeatCategory `json:"category" validate:"required"`
	PricePerKg float64      `json:"pricePerKg" validate:"required,gt=0"`
	Image      string       `json:"image,omitempty"`
}

// UpdateProductRequest represents the payload to update a product.
type UpdateProductRequest struct {
	Name       *string       `json:"name,omitempty" validate:"omitempty,min=2"`
	Category   *MeatCategory `json:"category,omitempty"`
	PricePerKg *float64      `json:"pricePerKg,omitempty" validate:"omitempty,gt=0"`
	Image      *string       `json:"image,omitempty"`
	InStock    *bool         `json:"inStock,omitempty"`
}
