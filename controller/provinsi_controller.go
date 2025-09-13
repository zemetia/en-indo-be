package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/service"
)

type ProvinsiController interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}

type provinsiController struct {
	provinsiService service.ProvinsiService
}

func NewProvinsiController(provinsiService service.ProvinsiService) ProvinsiController {
	return &provinsiController{
		provinsiService: provinsiService,
	}
}

func (c *provinsiController) GetAll(ctx *gin.Context) {
	provinsi, err := c.provinsiService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success", "data": provinsi})
}

func (c *provinsiController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	provinsi, err := c.provinsiService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success", "data": provinsi})
}
