package controller

import (
	"net/http"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	MeUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	jwtService  services.JWTService
	userService services.UserService
}

func NewUserController(us services.UserService, jwt services.JWTService) UserController {
	return &userController{
		jwtService:  jwt,
		userService: us,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	
	result, err := uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetAllUser(ctx *gin.Context) {
	result, err := uc.userService.GetAllUser(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_USER, err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) MeUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userId, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_TOKEN, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	result, err := uc.userService.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	
	res, _ := uc.userService.Verify(ctx.Request.Context(), req.Email, req.Password)
	if !res {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, dto.MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD, utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user, err := uc.userService.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	userResponse := entities.Authorization{
		Token: token,
		Role:  user.Role,
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, userResponse)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userId, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_TOKEN, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err = uc.userService.UpdateUser(ctx.Request.Context(), req, userId); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) DeleteUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_TOKEN, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err = uc.userService.DeleteUser(ctx.Request.Context(), userID); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_USER, nil)
	ctx.JSON(http.StatusOK, res)
}