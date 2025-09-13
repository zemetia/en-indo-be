package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type UserController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetByEmail(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Login(ctx *gin.Context)
	UploadProfileImage(ctx *gin.Context)
	SetupPassword(ctx *gin.Context)
	ToggleActivationStatus(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var req dto.UserCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	res, err := c.userService.Register(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_REGISTER_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": dto.MESSAGE_SUCCESS_REGISTER_USER,
		"data":    res,
	})
}

func (c *userController) GetAll(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	res, err := c.userService.GetAllUserWithPagination(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_GET_LIST_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_GET_LIST_USER,
		"data":    res,
	})
}

func (c *userController) GetByID(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   "user id is required",
		})
		return
	}

	res, err := c.userService.GetUserById(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": dto.MESSAGE_FAILED_GET_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_GET_USER,
		"data":    res,
	})
}

func (c *userController) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   "email is required",
		})
		return
	}

	res, err := c.userService.GetByEmail(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": dto.MESSAGE_FAILED_GET_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_GET_USER,
		"data":    res,
	})
}

func (c *userController) Update(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   "user id is required",
		})
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	res, err := c.userService.Update(ctx, req, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_UPDATE_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_UPDATE_USER,
		"data":    res,
	})
}

func (c *userController) Delete(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   "user id is required",
		})
		return
	}

	err := c.userService.Delete(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_DELETE_USER,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_DELETE_USER,
	})
}

func (c *userController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	res, err := c.userService.Verify(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": dto.MESSAGE_FAILED_LOGIN,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_LOGIN,
		"data":    res,
	})
}

// func (c *userController) UpdatePassword(ctx *gin.Context) {
// 	id, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	var req dto.UpdatePasswordRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := c.userService.UpdatePassword(id, &req); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
// }

func (c *userController) UploadProfileImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := c.userService.UploadProfileImage(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func (c *userController) SetupPassword(ctx *gin.Context) {
	// Get user ID from JWT token context
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in token",
			"error":   "unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.PasswordSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	if err := c.userService.SetupPassword(ctx, userID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to setup password",
			"error":   err.Error(),
		})
		return
	}

	var message string
	if req.Action == "change" {
		message = "Password updated successfully"
	} else {
		message = "Password setup completed successfully"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": dto.PasswordSetupResponse{
			Success: true,
			Message: message,
		},
	})
}

func (c *userController) ToggleActivationStatus(ctx *gin.Context) {
	personID := ctx.Param("person_id")

	// Debug logging
	fmt.Printf("[DEBUG] ToggleActivationStatus - PersonID: %s\n", personID)
	fmt.Printf("[DEBUG] Request Method: %s\n", ctx.Request.Method)
	fmt.Printf("[DEBUG] Content-Type: %s\n", ctx.GetHeader("Content-Type"))
	fmt.Printf("[DEBUG] Authorization: %s\n", ctx.GetHeader("Authorization"))

	if personID == "" {
		fmt.Printf("[ERROR] Person ID is empty\n")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Person ID is required",
			"error":   "person_id parameter is missing",
		})
		return
	}

	personUUID, err := uuid.Parse(personID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse person ID: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid person ID format",
			"error":   err.Error(),
		})
		return
	}

	// Read raw body for debugging
	body, _ := ctx.GetRawData()
	fmt.Printf("[DEBUG] Raw request body: %s\n", string(body))

	// Reset body reader for ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[ERROR] Failed to bind JSON: %v\n", err)
		fmt.Printf("[ERROR] Request body was: %s\n", string(body))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("[DEBUG] Parsed request - IsActive: %v\n", req.IsActive)

	// Get current user's person ID from JWT token to prevent self-deactivation
	currentPersonIDStr, exists := ctx.Get("person_id")
	if exists {
		currentPersonID, err := uuid.Parse(currentPersonIDStr.(string))
		if err == nil && currentPersonID == personUUID && !req.IsActive {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "You cannot deactivate your own account",
				"error":   "self_deactivation_not_allowed",
			})
			return
		}
	}

	// Handle three-state logic: no account, inactive account, active account
	if err := c.userService.ToggleUserActivationStatus(ctx, personUUID, req.IsActive); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user activation status",
			"error":   err.Error(),
		})
		return
	}

	var actionText string
	if req.IsActive {
		actionText = "Account has been activated successfully"
	} else {
		actionText = "Account has been deactivated successfully"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": actionText,
		"data": gin.H{
			"person_id":        personID,
			"is_active":        req.IsActive,
			"has_user_account": true, // After this operation, user account will always exist
		},
	})
}
