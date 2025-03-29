package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type UserController interface {
	Register(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetByEmail(ctx *gin.Context)
	SendVerificationEmail(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Login(ctx *gin.Context)
	UploadProfileImage(ctx *gin.Context)
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

func (c *userController) SendVerificationEmail(ctx *gin.Context) {
	var req dto.SendVerificationEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	err := c.userService.SendVerificationEmail(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_VERIFY_EMAIL,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SEND_VERIFICATION_EMAIL_SUCCESS,
	})
}

func (c *userController) VerifyEmail(ctx *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": dto.MESSAGE_FAILED_GET_DATA_FROM_BODY,
			"error":   err.Error(),
		})
		return
	}

	res, err := c.userService.VerifyEmail(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": dto.MESSAGE_FAILED_VERIFY_EMAIL,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": dto.MESSAGE_SUCCESS_VERIFY_EMAIL,
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
