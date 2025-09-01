package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type PelayananController interface {
	GetMyPelayanan(ctx *gin.Context)
	GetAllAssignments(ctx *gin.Context)
	GetAllPelayanan(ctx *gin.Context)
	AssignPelayanan(ctx *gin.Context)
	UnassignPelayanan(ctx *gin.Context)
	UpdatePelayananAssignment(ctx *gin.Context)
	GetAssignmentByID(ctx *gin.Context)
}

type pelayananController struct {
	pelayananService service.PelayananService
}

func NewPelayananController(pelayananService service.PelayananService) PelayananController {
	return &pelayananController{
		pelayananService: pelayananService,
	}
}

func (c *pelayananController) GetMyPelayanan(ctx *gin.Context) {
	// Get person ID from JWT claims (assuming it's stored in context)
	personIDStr, exists := ctx.Get("person_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Person ID not found in token",
			"error":   "unauthorized",
		})
		return
	}

	personID, err := uuid.Parse(personIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid person ID format",
			"error":   err.Error(),
		})
		return
	}

	pelayanan, err := c.pelayananService.GetMyPelayanan(ctx, personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get pelayanan assignments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get pelayanan assignments",
		"data":    pelayanan,
	})
}

func (c *pelayananController) GetAllAssignments(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	search := ctx.DefaultQuery("search", "")

	req := dto.PaginationRequest{
		Page:    page,
		PerPage: perPage,
		Search:  search,
	}

	assignments, err := c.pelayananService.GetAllAssignments(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get pelayanan assignments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get all pelayanan assignments",
		"data":    assignments,
	})
}

func (c *pelayananController) GetAllPelayanan(ctx *gin.Context) {
	pelayanan, err := c.pelayananService.GetAllPelayanan(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get all pelayanan",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get all pelayanan",
		"data":    pelayanan,
	})
}

func (c *pelayananController) AssignPelayanan(ctx *gin.Context) {
	var req dto.AssignPelayananRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	err := c.pelayananService.AssignPelayanan(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to assign pelayanan",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success assign pelayanan",
	})
}

func (c *pelayananController) UnassignPelayanan(ctx *gin.Context) {
	assignmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid assignment ID format",
			"error":   err.Error(),
		})
		return
	}

	err = c.pelayananService.UnassignPelayanan(ctx, assignmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to unassign pelayanan",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success unassign pelayanan",
	})
}

func (c *pelayananController) UpdatePelayananAssignment(ctx *gin.Context) {
	assignmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid assignment ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.UpdatePelayananAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	err = c.pelayananService.UpdatePelayananAssignment(ctx, assignmentID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update pelayanan assignment",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success update pelayanan assignment",
	})
}

func (c *pelayananController) GetAssignmentByID(ctx *gin.Context) {
	assignmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid assignment ID format",
			"error":   err.Error(),
		})
		return
	}

	assignment, err := c.pelayananService.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get pelayanan assignment",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get pelayanan assignment",
		"data":    assignment,
	})
}