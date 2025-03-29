package controller

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/zemetia/en-indo-be/dto"
// 	"github.com/zemetia/en-indo-be/service"
// )

// type ChurchController struct {
// 	churchService *service.ChurchService
// }

// func NewChurchController(churchService *service.ChurchService) *ChurchController {
// 	return &ChurchController{
// 		churchService: churchService,
// 	}
// }

// func (c *ChurchController) Create(ctx *gin.Context) {
// 	var req dto.ChurchRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	church, err := c.churchService.Create(&req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, church)
// }

// func (c *ChurchController) GetAll(ctx *gin.Context) {
// 	churches, err := c.churchService.GetAll()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, churches)
// }

// func (c *ChurchController) GetByID(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	church, err := c.churchService.GetByID(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, church)
// }

// func (c *ChurchController) GetByKabupatenID(ctx *gin.Context) {
// 	kabupatenID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	churches, err := c.churchService.GetByKabupatenID(kabupatenID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, churches)
// }

// func (c *ChurchController) GetByProvinsiID(ctx *gin.Context) {
// 	provinsiID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	churches, err := c.churchService.GetByProvinsiID(provinsiID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, churches)
// }

// func (c *ChurchController) Update(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	var req dto.ChurchRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	church, err := c.churchService.Update(id, &req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, church)
// }

// func (c *ChurchController) Delete(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	if err := c.churchService.Delete(id); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Gereja berhasil dihapus"})
// }
