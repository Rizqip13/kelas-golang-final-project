package controllers

import (
	"final-project/database"
	"final-project/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type VariantRequest struct {
	VariantName string `json:"variant_name" form:"variant_name" binding:"required"`
	Quantity    uint   `json:"quantity" form:"quantity" binding:"required"`
	ProductID   string `json:"product_id" form:"product_id"`
}

func GetVariants(ctx *gin.Context) {
	// Todo: pagination
	db := database.GetDB()

	var variants []models.Variant

	err := db.Preload("Product").Find(&variants).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": variants})
}

func CreateVariant(ctx *gin.Context) {
	db := database.GetDB()

	userData := ctx.MustGet("userData").(jwt5.MapClaims)

	var variantReq VariantRequest
	var p models.Product
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find product id
	err := db.Where("uuid = ?", variantReq.ProductID).Where("admin_id = ?", uint(userData["id"].(float64))).First(&p).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant := models.Variant{
		VariantName: variantReq.VariantName,
		Quantity:    variantReq.Quantity,
		ProductID:   p.ID,
	}

	err = db.Create(&variant).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": variant,
	})
}

func GetVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	variantUUID := ctx.Param("variantUUID")

	var variant models.Variant

	err := db.Preload("Product").Where("uuid = ?", variantUUID).First(&variant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Variant Not Found",
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

	ctx.JSON(http.StatusOK, gin.H{"data": variant})
}

func UpdateVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	variantUUID := ctx.Param("variantUUID")

	var variantDB models.Variant
	var variantReq VariantRequest

	err := ctx.ShouldBind(&variantReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	err = db.Preload("Product").Where("uuid = ?", variantUUID).First(&variantDB).Error
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

	variantDB.VariantName = variantReq.VariantName
	variantDB.Quantity = variantReq.Quantity

	err = db.Save(&variantDB).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": variantDB,
	})
}

func DeleteVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	variantUUID := ctx.Param("variantUUID")

	res := db.Where("uuid = ?", variantUUID).Delete(&models.Variant{})

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
