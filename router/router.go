package router

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", controllers.AdminRegister)
		authRouter.POST("/login", controllers.AdminLogin)
	}

	productRouter := router.Group("/products")
	{
		productRouter.Use(middlewares.Authentication())
		productRouter.POST("", controllers.CreateProduct)
		productRouter.GET("", controllers.GetProducts)
		productRouter.GET("/:productUUID", controllers.GetProductByUUID)
		productRouter.Use(middlewares.ProductAuthorization())
		productRouter.PUT("/:productUUID", controllers.UpdateProductByUUID)
		productRouter.DELETE("/:productUUID", controllers.DeleteProductByUUID)

	}

	variantRouter := router.Group("/products/variants")
	{
		variantRouter.Use(middlewares.Authentication())
		variantRouter.POST("", middlewares.VariantProductAuthorization(), controllers.CreateVariant)
		variantRouter.GET("", controllers.GetVariants)
		variantRouter.GET("/:variantUUID", controllers.GetVariantByUUID)
		variantRouter.Use(middlewares.VariantAuthorization())
		variantRouter.PUT("/:variantUUID", controllers.UpdateVariantByUUID)
		variantRouter.DELETE("/:variantUUID", controllers.DeleteVariantByUUID)

	}

	return router
}
