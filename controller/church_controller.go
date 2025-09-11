package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type ChurchController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetByKabupatenID(ctx *gin.Context)
	GetByProvinsiID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type churchController struct {
	churchService *service.ChurchService
}

func NewChurchController(churchService *service.ChurchService) ChurchController {
	return &churchController{
		churchService: churchService,
	}
}

func (c *churchController) Create(ctx *gin.Context) {
	log.Printf("[INFO] Church Controller: Received create church request")

	var req dto.ChurchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Church Controller: Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	// Log the incoming request data in detail
	log.Printf("[INFO] Church Controller: Creating church with data:")
	log.Printf("  - Name: %s", req.Name)
	log.Printf("  - Address: %s", req.Address)
	log.Printf("  - ChurchCode: '%s'", req.ChurchCode)
	log.Printf("  - Phone: %s", req.Phone)
	log.Printf("  - Email: %s", req.Email)
	log.Printf("  - Website: %s", req.Website)
	log.Printf("  - Latitude: %f", req.Latitude)
	log.Printf("  - Longitude: %f", req.Longitude)
	log.Printf("  - KabupatenID: %d", req.KabupatenID)

	// Validate required fields
	if req.Name == "" {
		log.Printf("[ERROR] Church Controller: Name is required but empty")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Name is required",
			"error":   "name field cannot be empty",
		})
		return
	}

	if req.Address == "" {
		log.Printf("[ERROR] Church Controller: Address is required but empty")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Address is required",
			"error":   "address field cannot be empty",
		})
		return
	}

	// Church code is now optional, but if provided, it should not be just whitespace
	if req.ChurchCode != "" && len(strings.TrimSpace(req.ChurchCode)) == 0 {
		log.Printf("[ERROR] Church Controller: ChurchCode contains only whitespace")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Church code cannot be only whitespace",
			"error":   "church_code field cannot contain only whitespace",
		})
		return
	}

	if req.KabupatenID == 0 {
		log.Printf("[ERROR] Church Controller: KabupatenID is required but zero")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Kabupaten ID is required",
			"error":   "kabupaten_id field cannot be zero",
		})
		return
	}

	church, err := c.churchService.Create(&req)
	if err != nil {
		log.Printf("[ERROR] Church Controller: Service error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create church",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[INFO] Church Controller: Successfully created church: %s", church.Name)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success create church",
		"data":    church,
	})
}

func (c *churchController) GetAll(ctx *gin.Context) {
	// Check if user wants all records without pagination
	all := ctx.Query("all")
	perPageStr := ctx.DefaultQuery("per_page", "10")

	if all == "true" || perPageStr == "0" {
		churches, err := c.churchService.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get churches",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success get all churches",
			"data":    churches,
			"count":   len(churches),
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(perPageStr)

	churches, err := c.churchService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get churches",
			"error":   err.Error(),
		})
		return
	}

	// Calculate pagination
	total := len(churches)
	startIdx := (page - 1) * perPage
	endIdx := startIdx + perPage

	if startIdx >= total {
		startIdx = 0
		endIdx = 0
		churches = []dto.ChurchResponse{}
	} else {
		if endIdx > total {
			endIdx = total
		}
		churches = churches[startIdx:endIdx]
	}

	maxPage := (total + perPage - 1) / perPage
	if maxPage == 0 {
		maxPage = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Success get all churches",
		"data":     churches,
		"page":     page,
		"per_page": perPage,
		"max_page": maxPage,
		"count":    total,
	})
}

func (c *churchController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	church, err := c.churchService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Church not found",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get church",
		"data":    church,
	})
}

func (c *churchController) GetByKabupatenID(ctx *gin.Context) {
	kabupatenIDStr := ctx.Param("id")
	kabupatenID, err := strconv.ParseUint(kabupatenIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid kabupaten ID format",
			"error":   err.Error(),
		})
		return
	}

	churches, err := c.churchService.GetByKabupatenID(uint(kabupatenID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get churches by kabupaten",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get churches by kabupaten",
		"data":    churches,
	})
}

func (c *churchController) GetByProvinsiID(ctx *gin.Context) {
	provinsiIDStr := ctx.Param("id")
	provinsiID, err := strconv.ParseUint(provinsiIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid provinsi ID format",
			"error":   err.Error(),
		})
		return
	}

	churches, err := c.churchService.GetByProvinsiID(uint(provinsiID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get churches by provinsi",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get churches by provinsi",
		"data":    churches,
	})
}

func (c *churchController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.ChurchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	church, err := c.churchService.Update(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update church",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success update church",
		"data":    church,
	})
}

func (c *churchController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.churchService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete church",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success delete church",
	})
}
