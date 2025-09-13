package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/service"
	"github.com/zemetia/en-indo-be/utils"
)

// RequireLifeGroupManageAccess middleware ensures user can manage the lifegroup
// This checks if user is PIC Lifegroup for the church OR leader/co-leader of this specific lifegroup
func RequireLifeGroupManageAccess(lifeGroupService service.LifeGroupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from authentication middleware
		userIDInterface, exists := ctx.Get("user_id")
		if !exists {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "User not authenticated", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID type", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Get lifegroup ID from URL parameter
		lifeGroupIDStr := ctx.Param("id")
		if lifeGroupIDStr == "" {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Lifegroup ID is required", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		lifeGroupID, err := uuid.Parse(lifeGroupIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid lifegroup ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if user can manage this lifegroup
		canManage, err := lifeGroupService.CheckUserCanManageLifeGroup(ctx, userID, lifeGroupID)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Failed to check permissions", nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		if !canManage {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "You are not authorized to manage this lifegroup. Only PIC Lifegroup or leaders/co-leaders can perform this action.", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Store lifegroup ID in context for use in handlers
		ctx.Set("lifegroup_id", lifeGroupID)
		ctx.Next()
	}
}

// RequireLifeGroupViewAccess middleware ensures user can view the lifegroup
// This checks if user is PIC, leader, co-leader, or member of this lifegroup
func RequireLifeGroupViewAccess(lifeGroupService service.LifeGroupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from authentication middleware
		userIDInterface, exists := ctx.Get("user_id")
		if !exists {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "User not authenticated", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID type", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Get lifegroup ID from URL parameter
		lifeGroupIDStr := ctx.Param("id")
		if lifeGroupIDStr == "" {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Lifegroup ID is required", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		lifeGroupID, err := uuid.Parse(lifeGroupIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid lifegroup ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if user can view this lifegroup
		canView, err := lifeGroupService.CheckUserCanViewLifeGroup(ctx, userID, lifeGroupID)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Failed to check permissions", nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		if !canView {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "You are not authorized to view this lifegroup", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Store lifegroup ID in context for use in handlers
		ctx.Set("lifegroup_id", lifeGroupID)
		ctx.Next()
	}
}

// RequireLifeGroupEditAccess middleware ensures user can edit the lifegroup
// This checks if user is PIC Lifegroup for the church OR leader/co-leader of this specific lifegroup
func RequireLifeGroupEditAccess(lifeGroupService service.LifeGroupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from authentication middleware
		userIDInterface, exists := ctx.Get("user_id")
		if !exists {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "User not authenticated", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID type", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Get lifegroup ID from URL parameter
		lifeGroupIDStr := ctx.Param("id")
		if lifeGroupIDStr == "" {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Lifegroup ID is required", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		lifeGroupID, err := uuid.Parse(lifeGroupIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid lifegroup ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if user can edit this lifegroup
		canEdit, err := lifeGroupService.CheckUserCanEditLifeGroup(ctx, userID, lifeGroupID)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Failed to check permissions", nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		if !canEdit {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "You are not authorized to edit this lifegroup. Only PIC Lifegroup or leaders/co-leaders can perform this action.", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Store lifegroup ID in context for use in handlers
		ctx.Set("lifegroup_id", lifeGroupID)
		ctx.Next()
	}
}

// RequireLifeGroupDeleteAccess middleware ensures user can delete the lifegroup
// This checks if user is PIC Lifegroup for the church OR leader (NOT co-leader) of this specific lifegroup
func RequireLifeGroupDeleteAccess(lifeGroupService service.LifeGroupService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get user ID from authentication middleware
		userIDInterface, exists := ctx.Get("user_id")
		if !exists {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "User not authenticated", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID type", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid user ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Get lifegroup ID from URL parameter
		lifeGroupIDStr := ctx.Param("id")
		if lifeGroupIDStr == "" {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Lifegroup ID is required", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		lifeGroupID, err := uuid.Parse(lifeGroupIDStr)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Invalid lifegroup ID format", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Check if user can delete this lifegroup
		canDelete, err := lifeGroupService.CheckUserCanDeleteLifeGroup(ctx, userID, lifeGroupID)
		if err != nil {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "Failed to check permissions", nil)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		if !canDelete {
			response := utils.BuildResponseFailed("ACCESS_DENIED", "You are not authorized to delete this lifegroup. Only PIC Lifegroup or leaders can perform this action.", nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// Store lifegroup ID in context for use in handlers
		ctx.Set("lifegroup_id", lifeGroupID)
		ctx.Next()
	}
}
