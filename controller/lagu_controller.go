package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/service"
)

type LaguController interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type laguController struct {
	laguService service.LaguService
}

func NewLaguController(laguService service.LaguService) LaguController {
	return &laguController{
		laguService: laguService,
	}
}

func (c *laguController) Create(ctx *gin.Context) {
	var req dto.LaguRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var originalLaguID *uuid.UUID
	if req.OriginalLaguID != "" {
		id, err := uuid.Parse(req.OriginalLaguID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OriginalLaguID format: " + err.Error()})
			return
		}
		originalLaguID = &id
	}

	lagu := entity.Lagu{
		Judul:          req.Judul,
		Artis:          req.Artis,
		YoutubeLink:    req.YoutubeLink,
		Genre:          req.Genre,
		Lirik:          req.Lirik,
		Tags:           req.Tags,
		NadaDasar:      req.NadaDasar,
		TahunRilis:     req.TahunRilis,
		OriginalLaguID: originalLaguID,
	}

	if err := c.laguService.Create(ctx, &lagu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, lagu)
}

func (c *laguController) FindAll(ctx *gin.Context) {
	lagus, err := c.laguService.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lagus)
}

func (c *laguController) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	lagu, err := c.laguService.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	ctx.JSON(http.StatusOK, lagu)
}

func (c *laguController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.LaguRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var originalLaguID *uuid.UUID
	if req.OriginalLaguID != "" {
		oid, err := uuid.Parse(req.OriginalLaguID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OriginalLaguID format"})
			return
		}
		originalLaguID = &oid
	}

	lagu := entity.Lagu{
		ID:             id,
		Judul:          req.Judul,
		Artis:          req.Artis,
		YoutubeLink:    req.YoutubeLink,
		Genre:          req.Genre,
		Lirik:          req.Lirik,
		Tags:           req.Tags,
		NadaDasar:      req.NadaDasar,
		TahunRilis:     req.TahunRilis,
		OriginalLaguID: originalLaguID,
	}

	if err := c.laguService.Update(ctx, &lagu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lagu)
}

func (c *laguController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.laguService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
