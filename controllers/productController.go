package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ProductRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
}

func GetProducts(ctx *gin.Context) {
	// Todo: pagination
	db := database.GetDB()

	var products []models.Product

	err := db.Preload("Admin").Preload("Variants").Find(&products).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()

	userData := ctx.MustGet("userData").(jwt5.MapClaims)

	var productReq ProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadResult, err := helpers.UploadFile(productReq.Image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:     productReq.Name,
		ImageUrl: uploadResult,
		AdminID:  uint(userData["id"].(float64)),
	}

	err = db.Create(&product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func UpdateProductByUUID(ctx *gin.Context) {
	var err error
	db := database.GetDB()
	productUUID := ctx.Param("productUUID")

	var productDB models.Product
	var productReq ProductRequest

	err = ctx.ShouldBind(&productReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	err = db.Preload("Admin").Preload("Variants").Where("uuid = ?", productUUID).First(&productDB).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Product Not Found",
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
		}
		return
	}

	productDB.Name = productReq.Name

	if productReq.Image != nil {
		uploadResult, err := helpers.UploadFile(productReq.Image)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		productDB.ImageUrl = uploadResult
	}

	err = db.Save(&productDB).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": productDB,
	})
}

func DeleteProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	productUUID := ctx.Param("productUUID")

	res := db.Where("uuid = ?", productUUID).Delete(&models.Product{})

	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": res.Error.Error(),
		})
		return
	}

	if res.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "record Not Found",
			"message": "Data Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": fmt.Sprintf("%v row(s) affected", res.RowsAffected),
	})

}

func GetProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	productUUID := ctx.Param("productUUID")

	var product models.Product

	err := db.Preload("Admin").Preload("Variants").Where("uuid = ?", productUUID).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Product Not Found",
				"message": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}
