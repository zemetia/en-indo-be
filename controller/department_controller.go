package controller

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/zemetia/en-indo-be/dto"
// 	"github.com/zemetia/en-indo-be/service"
// )

// type DepartmentController struct {
// 	departmentService *service.DepartmentService
// }

// func NewDepartmentController(departmentService *service.DepartmentService) *DepartmentController {
// 	return &DepartmentController{
// 		departmentService: departmentService,
// 	}
// }

// func (c *DepartmentController) Create(ctx *gin.Context) {
// 	var req dto.DepartmentRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	department, err := c.departmentService.Create(&req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, department)
// }

// func (c *DepartmentController) GetAll(ctx *gin.Context) {
// 	departments, err := c.departmentService.GetAll()
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, departments)
// }

// func (c *DepartmentController) GetByID(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	department, err := c.departmentService.GetByID(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, department)
// }

// func (c *DepartmentController) GetByChurchID(ctx *gin.Context) {
// 	churchID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	departments, err := c.departmentService.GetByChurchID(churchID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, departments)
// }

// func (c *DepartmentController) Update(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	var req dto.DepartmentRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	department, err := c.departmentService.Update(id, &req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, department)
// }

// func (c *DepartmentController) Delete(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
// 		return
// 	}

// 	if err := c.departmentService.Delete(id); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Departemen berhasil dihapus"})
// }
