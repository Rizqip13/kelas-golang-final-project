package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UUID      string `gorm:"not null" json:"uuid"`
	Name      string `gorm:"not null" json:"name" valid:"required~Product name is required"`
	ImageUrl  string `gorm:"no null"`
	AdminID   uint
	Admin     *Admin
	Variants  *[]Variant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"variants"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(product)

	if err != nil {
		return
	}

	product.UUID = uuid.New().String()

	err = nil
	return
}
