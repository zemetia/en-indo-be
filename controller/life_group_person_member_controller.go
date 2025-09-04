package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type LifeGroupPersonMemberController interface {
	AddPersonMember(ctx *gin.Context)
	GetPersonMembers(ctx *gin.Context)
	UpdatePersonMemberPosition(ctx *gin.Context)
	RemovePersonMember(ctx *gin.Context)
	GetPersonMemberByID(ctx *gin.Context)
	GetPersonLifeGroups(ctx *gin.Context)
	GetLeadershipStructure(ctx *gin.Context)
}

type lifeGroupPersonMemberController struct {
	personMemberService service.LifeGroupPersonMemberService
}

func NewLifeGroupPersonMemberController(personMemberService service.LifeGroupPersonMemberService) LifeGroupPersonMemberController {
	return &lifeGroupPersonMemberController{
		personMemberService: personMemberService,
	}
}

func (c *lifeGroupPersonMemberController) AddPersonMember(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.AddPersonMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.AddPersonMember(ctx, lifeGroupID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add person member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Person member added successfully",
		"data":    result,
	})
}

func (c *lifeGroupPersonMemberController) GetPersonMembers(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.GetPersonMembers(ctx, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get person members",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person members retrieved successfully",
		"data":    result,
	})
}

func (c *lifeGroupPersonMemberController) UpdatePersonMemberPosition(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.UpdatePersonMemberPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.UpdatePersonMemberPosition(ctx, lifeGroupID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update person member position",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person member position updated successfully",
		"data":    result,
	})
}

func (c *lifeGroupPersonMemberController) RemovePersonMember(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.RemovePersonMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = c.personMemberService.RemovePersonMember(ctx, lifeGroupID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove person member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person member removed successfully",
	})
}

func (c *lifeGroupPersonMemberController) GetPersonMemberByID(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("memberID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid member ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.GetPersonMemberByID(ctx, memberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get person member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person member retrieved successfully",
		"data":    result,
	})
}

func (c *lifeGroupPersonMemberController) GetPersonLifeGroups(ctx *gin.Context) {
	personID, err := uuid.Parse(ctx.Param("personID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid person ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.GetPersonLifeGroups(ctx, personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get person lifegroups",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person lifegroups retrieved successfully",
		"data":    result,
	})
}

func (c *lifeGroupPersonMemberController) GetLeadershipStructure(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.personMemberService.GetLeadershipStructure(ctx, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get leadership structure",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Leadership structure retrieved successfully",
		"data":    result,
	})
}