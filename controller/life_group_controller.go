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
	UpdateMembers(ctx *gin.Context)
	AddToLifeGroup(ctx *gin.Context)
	RemoveFromLifeGroup(ctx *gin.Context)
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

func (c *lifeGroupController) UpdateMembers(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.UpdateMembersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.lifeGroupService.UpdateMembers(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *lifeGroupController) AddToLifeGroup(ctx *gin.Context) {
	personID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID format"})
		return
	}

	lifeGroupID, err := uuid.Parse(ctx.Param("life_group_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid life group ID format"})
		return
	}

	err = c.lifeGroupService.AddToLifeGroup(ctx, personID, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person added to life group successfully"})
}

func (c *lifeGroupController) RemoveFromLifeGroup(ctx *gin.Context) {
	personID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID format"})
		return
	}

	lifeGroupID, err := uuid.Parse(ctx.Param("life_group_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid life group ID format"})
		return
	}

	err = c.lifeGroupService.RemoveFromLifeGroup(ctx, personID, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person removed from life group successfully"})
}
