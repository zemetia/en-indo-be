package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/service"
)

type KabupatenController struct {
	kabupatenService *service.KabupatenService
}

func NewKabupatenController(kabupatenService *service.KabupatenService) *KabupatenController {
	return &KabupatenController{
		kabupatenService: kabupatenService,
	}
}

func (c *KabupatenController) GetAll(ctx *gin.Context) {
	kabupaten, err := c.kabupatenService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, kabupaten)
}

func (c *KabupatenController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	kabupaten, err := c.kabupatenService.GetByID(uint(id)) // Mengubah id menjadi uint
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, kabupaten)
}

func (c *KabupatenController) GetByProvinsiID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	kabupaten, err := c.kabupatenService.GetByProvinsiID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, kabupaten)
}
