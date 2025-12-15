package routes

import (
	"product-api/controllers"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, product *controllers.ProductController) {
	r.POST("/products", product.Create)
	r.GET("/products", product.List)
	r.GET("/products/:id", product.FindByID)
	r.PUT("/products/:id", product.Update)
	r.DELETE("/products/:id", product.Delete)
}
