package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variant struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UUID        string `gorm:"not null" json:"uuid"`
	VariantName string `gorm:"not null" json:"variant_name" valid:"required~Variant Name is required"`
	Quantity    uint   `gorm:"not null" json:"quantity" valid:"required~Quantity of variant is required, numeric~Invalid quantity format"`
	ProductID   uint
	Product     *Product
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (variant *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(variant)

	if err != nil {
		return
	}

	// generate a new UUID
	variant.UUID = uuid.New().String()

	err = nil
	return
}
