package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRegister(ctx *gin.Context) {
	db := database.GetDB()
	contentType := ctx.Request.Header.Get("Content-Type")
	Admin := models.Admin{}

	if contentType == "application/json" {
		ctx.ShouldBindJSON(&Admin)
	} else {
		ctx.ShouldBind(&Admin)
	}

	err := db.Create(&Admin).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    Admin,
	})
}

func AdminLogin(ctx *gin.Context) {
	var err error
	db := database.GetDB()
	contentType := ctx.Request.Header.Get("Content-Type")
	Admin := models.Admin{}

	if contentType == "application/json" {
		err = ctx.ShouldBindJSON(&Admin)
	} else {
		err = ctx.ShouldBind(&Admin)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	rawPassword := Admin.Password

	err = db.Where("email = ?", Admin.Email).Take(&Admin).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(Admin.Password), []byte(rawPassword))
	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid password",
		})
		return
	}

	token := helpers.GenerateToken(Admin.ID, Admin.Email)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
