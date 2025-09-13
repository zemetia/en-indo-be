package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type VisitorInformationController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetByVisitorID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type visitorInformationController struct {
	visitorInfoService service.VisitorInformationService
}

func NewVisitorInformationController(visitorInfoService service.VisitorInformationService) VisitorInformationController {
	return &visitorInformationController{
		visitorInfoService: visitorInfoService,
	}
}

func (c *visitorInformationController) Create(ctx *gin.Context) {
	var req dto.VisitorInformationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorInfoService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create visitor information",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Visitor information created successfully",
		"data":    res,
	})
}

func (c *visitorInformationController) GetAll(ctx *gin.Context) {
	res, err := c.visitorInfoService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitor information",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor information retrieved successfully",
		"data":    res,
	})
}

func (c *visitorInformationController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor information ID",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorInfoService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Visitor information not found",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor information retrieved successfully",
		"data":    res,
	})
}

func (c *visitorInformationController) GetByVisitorID(ctx *gin.Context) {
	visitorID, err := uuid.Parse(ctx.Param("visitor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor ID",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorInfoService.GetByVisitorID(ctx, visitorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitor information",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor information retrieved successfully",
		"data":    res,
	})
}

func (c *visitorInformationController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor information ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.VisitorInformationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.visitorInfoService.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update visitor information",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor information updated successfully",
		"data":    res,
	})
}

func (c *visitorInformationController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor information ID",
			"error":   err.Error(),
		})
		return
	}

	err = c.visitorInfoService.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete visitor information",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor information deleted successfully",
	})
}
