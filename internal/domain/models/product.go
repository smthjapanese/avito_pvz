package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductType string

const (
	ProductTypeElectronics ProductType = "электроника"
	ProductTypeClothes     ProductType = "одежда"
	ProductTypeShoes       ProductType = "обувь"
)

type Product struct {
	ID          uuid.UUID   `json:"id"`
	DateTime    time.Time   `json:"date_time"`
	Type        ProductType `json:"type"`
	ReceptionID uuid.UUID   `json:"reception_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

func NewProduct(productType ProductType, receptionID uuid.UUID) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New(),
		DateTime:    now,
		Type:        productType,
		ReceptionID: receptionID,
		CreatedAt:   now,
	}
}

func IsValidProductType(productType ProductType) bool {
	return productType == ProductTypeElectronics ||
		productType == ProductTypeClothes ||
		productType == ProductTypeShoes
}
