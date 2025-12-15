package main

// @title Product API
// @version 1.0
// @description API de produtos com CRUD genérico e Oracle
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @host localhost:8080
// @BasePath /api

import (
	"github.com/gin-gonic/gin"

	_ "product-api/docs"
	"product-api/logger"

	"product-api/controllers"
	"product-api/crud"
	"product-api/database"
	"product-api/facade"
	"product-api/repository"
	"product-api/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db := database.OpenOracle()
	defer db.Close()

	crudSvc := crud.NewCrud(db, "")

	logger.Init()
	logger.Logger.Info("Iniciando a API...")

	baseRepo := repository.NewBaseRepository(crudSvc)
	productRepo := repository.NewProductRepository(baseRepo)
	productFacade := facade.NewProductFacade(productRepo)
	productController := controllers.NewProductController(productFacade)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	routes.Register(api, productController)

	r.Run(":8080")
}
