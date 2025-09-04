package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type LifeGroupVisitorMemberController interface {
	AddVisitorMember(ctx *gin.Context)
	GetVisitorMembers(ctx *gin.Context)
	RemoveVisitorMember(ctx *gin.Context)
	GetVisitorMemberByID(ctx *gin.Context)
	GetVisitorLifeGroups(ctx *gin.Context)
}

type lifeGroupVisitorMemberController struct {
	visitorMemberService service.LifeGroupVisitorMemberService
}

func NewLifeGroupVisitorMemberController(visitorMemberService service.LifeGroupVisitorMemberService) LifeGroupVisitorMemberController {
	return &lifeGroupVisitorMemberController{
		visitorMemberService: visitorMemberService,
	}
}

func (c *lifeGroupVisitorMemberController) AddVisitorMember(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.AddVisitorMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.visitorMemberService.AddVisitorMember(ctx, lifeGroupID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add visitor member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Visitor member added successfully",
		"data":    result,
	})
}

func (c *lifeGroupVisitorMemberController) GetVisitorMembers(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.visitorMemberService.GetVisitorMembers(ctx, lifeGroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitor members",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor members retrieved successfully",
		"data":    result,
	})
}

func (c *lifeGroupVisitorMemberController) RemoveVisitorMember(ctx *gin.Context) {
	lifeGroupID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid lifegroup ID",
			"error":   err.Error(),
		})
		return
	}

	var req dto.RemoveVisitorMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err = c.visitorMemberService.RemoveVisitorMember(ctx, lifeGroupID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove visitor member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor member removed successfully",
	})
}

func (c *lifeGroupVisitorMemberController) GetVisitorMemberByID(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("memberID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid member ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.visitorMemberService.GetVisitorMemberByID(ctx, memberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitor member",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor member retrieved successfully",
		"data":    result,
	})
}

func (c *lifeGroupVisitorMemberController) GetVisitorLifeGroups(ctx *gin.Context) {
	visitorID, err := uuid.Parse(ctx.Param("visitorID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid visitor ID",
			"error":   err.Error(),
		})
		return
	}

	result, err := c.visitorMemberService.GetVisitorLifeGroups(ctx, visitorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get visitor lifegroups",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Visitor lifegroups retrieved successfully",
		"data":    result,
	})
}