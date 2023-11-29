package models

import (
	"final-project/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UUID      string     `gorm:"not null" json:"uuid"`
	Name      string     `gorm:"not null" json:"name" valid:"required~Your name is required"`
	Email     string     `gorm:"not null;unique" json:"email" valid:"required~Your email is required"`
	Password  string     `gorm:"not null" json:"password" valid:"required~Your password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Products  *[]Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"products"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (admin *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	_, err = govalidator.ValidateStruct(admin)

	if err != nil {
		return
	}

	admin.Password = helpers.HashPass(admin.Password)

	// generate a new UUID
	admin.UUID = uuid.New().String()

	err = nil
	return
}
