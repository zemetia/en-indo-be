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
	GetByUserID(ctx *gin.Context)
	GetByPICLifegroupChurches(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.personService.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create person",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success create person",
		"data":    res,
	})
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
	if search.ChurchID != nil || search.KabupatenID != nil || search.UserID != nil || search.Name != nil {
		res, err := c.personService.Search(ctx, &search)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to search persons",
				"error":   err.Error(),
			})
			return
		}

		// Check if user wants all records without pagination for search results
		all := ctx.Query("all")
		perPageStr := ctx.DefaultQuery("per_page", "10")

		if all == "true" || perPageStr == "0" {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Success search persons",
				"data":    res,
				"count":   len(res),
			})
			return
		}

		// Apply pagination to search results
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		perPage, _ := strconv.Atoi(perPageStr)

		total := len(res)
		startIdx := (page - 1) * perPage
		endIdx := startIdx + perPage

		if startIdx >= total {
			startIdx = 0
			endIdx = 0
			res = []dto.SimplePersonResponse{}
		} else {
			if endIdx > total {
				endIdx = total
			}
			res = res[startIdx:endIdx]
		}

		maxPage := (total + perPage - 1) / perPage
		if maxPage == 0 {
			maxPage = 1
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success search persons",
			"data": gin.H{
				"data":     res,
				"page":     page,
				"per_page": perPage,
				"max_page": maxPage,
				"count":    total,
			},
		})
		return
	}

	// Jika tidak ada parameter pencarian, ambil semua data
	res, err := c.personService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get all persons",
			"error":   err.Error(),
		})
		return
	}

	// Check if user wants all records without pagination
	all := ctx.Query("all")
	perPageStr := ctx.DefaultQuery("per_page", "10")

	if all == "true" || perPageStr == "0" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success get all persons",
			"data":    res,
			"count":   len(res),
		})
		return
	}

	// Apply pagination
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(perPageStr)

	total := len(res)
	startIdx := (page - 1) * perPage
	endIdx := startIdx + perPage

	if startIdx >= total {
		startIdx = 0
		endIdx = 0
		res = []dto.SimplePersonResponse{}
	} else {
		if endIdx > total {
			endIdx = total
		}
		res = res[startIdx:endIdx]
	}

	maxPage := (total + perPage - 1) / perPage
	if maxPage == 0 {
		maxPage = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get all persons",
		"data": gin.H{
			"data":     res,
			"page":     page,
			"per_page": perPage,
			"max_page": maxPage,
			"count":    total,
		},
	})
}

func (c *personController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.personService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get person",
		"data":    res,
	})
}

func (c *personController) GetByUserID(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.personService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Person not found for user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get person by user",
		"data":    res,
	})
}

func (c *personController) GetByPICLifegroupChurches(ctx *gin.Context) {
	// Get person ID from JWT claims
	personIDStr, exists := ctx.Get("person_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Person ID not found in token. User may not have a Person record associated.",
			"error":   "no_person_id",
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

	persons, err := c.personService.GetByPICLifegroupChurches(ctx, personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get persons from PIC lifegroup churches",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get persons from PIC lifegroup churches",
		"data":    persons,
	})
}

func (c *personController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.PersonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.personService.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update person",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success update person",
		"data":    res,
	})
}

func (c *personController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	err = c.personService.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete person",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success delete person",
	})
}
