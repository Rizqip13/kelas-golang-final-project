package database

import (
	"final-project/configs"
	"final-project/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	db, err = gorm.Open(mysql.Open(configs.EnvDBConn()), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to databse: ", err)
	}

	db.AutoMigrate(models.Admin{}, models.Product{}, models.Variant{})
}

func GetDB() *gorm.DB {
	return db
}
