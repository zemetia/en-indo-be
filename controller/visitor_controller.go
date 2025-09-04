package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type VisitorController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Search(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type visitorController struct {
	visitorService service.VisitorService
}

func NewVisitorController(visitorService service.VisitorService) VisitorController {
	return &visitorController{
		visitorService: visitorService,
	}
}

func (c *visitorController) Create(ctx *gin.Context) {
	var req dto.VisitorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create visitor",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Visitor created successfully",
		"data":    res,
	})
}

func (c *visitorController) GetAll(ctx *gin.Context) {
	res, err := c.visitorService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitors",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitors retrieved successfully",
		"data":    res,
	})
}

func (c *visitorController) Search(ctx *gin.Context) {
	var search dto.VisitorSearchDto
	if err := ctx.ShouldBindQuery(&search); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse search parameters",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorService.Search(ctx, &search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to search visitors",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitors searched successfully",
		"data":    res,
	})
}

func (c *visitorController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor ID",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Visitor not found",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor retrieved successfully",
		"data":    res,
	})
}

func (c *visitorController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.VisitorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorService.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update visitor",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor updated successfully",
		"data":    res,
	})
}

func (c *visitorController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor ID",
			"error":   err.Error(),
		})
		return
	}

	err = c.visitorService.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete visitor",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor deleted successfully",
	})
}