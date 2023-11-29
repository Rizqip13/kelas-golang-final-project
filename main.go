package main

import (
	"final-project/database"
	"final-project/router"
)

var (
	PORT = ":8080"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(PORT)
}
