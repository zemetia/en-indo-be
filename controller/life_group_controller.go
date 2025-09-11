package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type LifeGroupController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UpdateLeader(ctx *gin.Context)
	GetByChurch(ctx *gin.Context)
	GetByUser(ctx *gin.Context)
	GetMyLifeGroup(ctx *gin.Context)
	GetDaftarLifeGroup(ctx *gin.Context)
	GetByMultipleChurches(ctx *gin.Context)
	GetLifeGroupsByPICRole(ctx *gin.Context)
}

type lifeGroupController struct {
	lifeGroupService service.LifeGroupService
}

func NewLifeGroupController(lifeGroupService service.LifeGroupService) LifeGroupController {
	return &lifeGroupController{
		lifeGroupService: lifeGroupService,
	}
}

func (c *lifeGroupController) Create(ctx *gin.Context) {
	var req dto.LifeGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.lifeGroupService.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *lifeGroupController) GetAll(ctx *gin.Context) {
	var search dto.PersonSearchDto

	if nameStr := ctx.Query("name"); nameStr != "" {
		search.Name = &nameStr
	}

	if churchIDStr := ctx.Query("church_id"); churchIDStr != "" {
		churchID, err := uuid.Parse(churchIDStr)
		if err == nil {
			search.ChurchID = &churchID
		}
	}

	if kabupatenIDStr := ctx.Query("kabupaten_id"); kabupatenIDStr != "" {
		kabupatenID, err := strconv.ParseUint(kabupatenIDStr, 10, 32)
		if err == nil {
			kabID := uint(kabupatenID)
			search.KabupatenID = &kabID
		}
	}

	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			search.UserID = &userID
		}
	}

	// Jika ada parameter pencarian, gunakan Search
	if search.ChurchID != nil || search.KabupatenID != nil || search.UserID != nil {
		res, err := c.lifeGroupService.Search(ctx, &search)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
		return
	}

	// Jika tidak ada parameter pencarian, ambil semua data
	res, err := c.lifeGroupService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *lifeGroupController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	response, err := c.lifeGroupService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.LifeGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.lifeGroupService.Update(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.lifeGroupService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "LifeGroup deleted successfully"})
}

func (c *lifeGroupController) UpdateLeader(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateLeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.lifeGroupService.UpdateLeader(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetByChurch(ctx *gin.Context) {
	churchID, err := uuid.Parse(ctx.Param("church_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID format"})
		return
	}

	response, err := c.lifeGroupService.GetByChurchID(churchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetByUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := c.lifeGroupService.GetByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetMyLifeGroup(ctx *gin.Context) {
	// Get user ID from authentication middleware
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := c.lifeGroupService.GetMyLifeGroup(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetDaftarLifeGroup(ctx *gin.Context) {
	// Get user ID from authentication middleware
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := c.lifeGroupService.GetDaftarLifeGroup(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetByMultipleChurches(ctx *gin.Context) {
	var req dto.BatchChurchLifeGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.ChurchIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one church ID is required"})
		return
	}

	response, err := c.lifeGroupService.GetByMultipleChurchIDs(req.ChurchIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) GetLifeGroupsByPICRole(ctx *gin.Context) {
	// Get user ID from authentication middleware
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	response, err := c.lifeGroupService.GetLifeGroupsByPICRole(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
