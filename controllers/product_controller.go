package controllers

import (
	"net/http"
	"strconv"

	"product-api/facade"
	"product-api/logger"
	"product-api/mappers"

	"product-api/dto/request"
	"product-api/dto/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ProductController struct {
	facade *facade.ProductFacade
}

func NewProductController(f *facade.ProductFacade) *ProductController {
	return &ProductController{facade: f}
}

// Create godoc
// @Summary Criar produto
// @Description Cria um novo produto
// @Tags Products
// @Accept json
// @Produce json
// @Param product body request.ProductRequestDTO true "Produto a ser criado"
// @Success 201 {object} response.ProductResponseDTO
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (c *ProductController) Create(ctx *gin.Context) {
	var req request.ProductRequestDTO

	// Bind do JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponseDTO{
			Status: http.StatusBadRequest,
			Info:   "Erro ao criar produto",
		})
		logger.Logger.WithFields(log.Fields{
			"method": "Create",
			"body":   req,
			"error":  err,
		}).Error("Falha ao bindar JSON")
		return
	}

	product := mappers.ToProductModel(req)

	// Criação do produto
	created, err := c.facade.Create(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponseDTO{
			Status: http.StatusInternalServerError,
			Info:   "Erro ao criar produto",
		})
		logger.Logger.WithFields(log.Fields{
			"method":  "Create",
			"product": product,
			"error":   err,
		}).Error("Erro ao salvar produto no banco")
		return
	}

	resp := mappers.ToProductResponse(created)
	ctx.JSON(http.StatusCreated, resp)
}

// List godoc
// @Summary Listar produtos
// @Description Retorna lista de produtos
// @Tags Products
// @Produce json
// @Success 200 {array} response.ProductResponseDTO
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (c *ProductController) List(ctx *gin.Context) {
	products, err := c.facade.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponseDTO{
			Status: http.StatusInternalServerError,
			Info:   "Erro ao listar produto",
		})
		logger.Logger.WithFields(log.Fields{
			"method": "Create",
			"error":  err,
		}).Error("Erro ao listar produtos do banco")
		return
	}

	resp := mappers.ToProductResponseList(products)
	ctx.JSON(http.StatusOK, resp)
}

// FindByID godoc
// @Summary Buscar produto por ID
// @Description Retorna um produto específico pelo ID
// @Tags Products
// @Produce json
// @Param id path int true "ID do produto"
// @Success 200 {object} response.ProductResponseDTO
// @Failure 500 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (c *ProductController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	p, err := c.facade.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resp := mappers.ToProductResponse(p)
	ctx.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary Atualizar produto
// @Description Atualiza os dados de um produto
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "ID do produto"
// @Param product body request.ProductRequestDTO true "Dados do produto"
// @Success 200 {object} response.ProductResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req request.ProductRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := mappers.ToProductModel(req)

	updated, err := c.facade.Update(id, product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := mappers.ToProductResponse(updated)
	ctx.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary Deletar produto
// @Description Remove um produto pelo ID
// @Tags Products
// @Param id path int true "ID do produto"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := c.facade.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
