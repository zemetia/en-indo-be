package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type PersonController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetByChurchID(ctx *gin.Context)
	GetByKabupatenID(ctx *gin.Context)
	GetByUserID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	AddToLifeGroup(ctx *gin.Context)
	RemoveFromLifeGroup(ctx *gin.Context)
}

type personController struct {
	personService service.PersonService
}

func NewPersonController(personService service.PersonService) PersonController {
	return &personController{
		personService: personService,
	}
}

func (c *personController) Create(ctx *gin.Context) {
	var req dto.PersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.personService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *personController) GetAll(ctx *gin.Context) {
	var search dto.PersonSearchDto

	// Cek parameter pencarian
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
		res, err := c.personService.Search(ctx, &search)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
		return
	}

	// Jika tidak ada parameter pencarian, ambil semua data
	res, err := c.personService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	res, err := c.personService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) GetByChurchID(ctx *gin.Context) {
	churchID, err := uuid.Parse(ctx.Param("church_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid church ID format"})
		return
	}

	res, err := c.personService.GetByChurchID(ctx, churchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) GetByKabupatenID(ctx *gin.Context) {
	kabupatenID, err := uuid.Parse(ctx.Param("kabupaten_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kabupaten ID format"})
		return
	}

	res, err := c.personService.GetByKabupatenID(ctx, kabupatenID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) GetByUserID(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	res, err := c.personService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req dto.PersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.personService.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *personController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = c.personService.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}

func (c *personController) AddToLifeGroup(ctx *gin.Context) {
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

	err = c.personService.AddToLifeGroup(ctx, personID, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person added to life group successfully"})
}

func (c *personController) RemoveFromLifeGroup(ctx *gin.Context) {
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

	err = c.personService.RemoveFromLifeGroup(ctx, personID, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Person removed from life group successfully"})
}
