package dto

import (
	"errors"
	"mime/multipart"

	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY      = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER           = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER           = "failed get list user"
	MESSAGE_FAILED_GET_USER_TOKEN          = "failed get user token"
	MESSAGE_FAILED_TOKEN_NOT_VALID         = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND         = "token not found"
	MESSAGE_FAILED_GET_USER                = "failed get user"
	MESSAGE_FAILED_LOGIN                   = "failed login"
	MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD = "wrong email or password"
	MESSAGE_FAILED_UPDATE_USER             = "failed update user"
	MESSAGE_FAILED_DELETE_USER             = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST          = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS           = "denied access"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER           = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER           = "success get list user"
	MESSAGE_SUCCESS_GET_USER                = "success get user"
	MESSAGE_SUCCESS_LOGIN                   = "success login"
	MESSAGE_SUCCESS_UPDATE_USER             = "success update user"
	MESSAGE_SUCCESS_DELETE_USER             = "success delete user"
)

var (
	ErrCreateUser             = errors.New("failed to create user")
	ErrGetAllUser             = errors.New("failed to get all user")
	ErrGetUserById            = errors.New("failed to get user by id")
	ErrGetByEmail             = errors.New("failed to get user by email")
	ErrEmailAlreadyExists     = errors.New("email already exist")
	ErrUpdateUser             = errors.New("failed to update user")
	ErrUserNotAdmin           = errors.New("user not admin")
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailNotFound          = errors.New("email not found")
	ErrDeleteUser             = errors.New("failed to delete user")
	ErrPasswordNotMatch       = errors.New("password not match")
	ErrEmailOrPassword        = errors.New("wrong email or password")
	ErrTokenInvalid           = errors.New("token invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrUploadProfileImage     = errors.New("failed to upload profile image")
	ErrGetPelayanan           = errors.New("failed to get pelayanan")
	ErrUserInactive           = errors.New("user account is inactive")
	ErrUserNoPelayanan        = errors.New("user has no pelayanan assignments")
)

type (
	UserCreateRequest struct {
		Name       string                `json:"name" form:"name"`
		TelpNumber string                `json:"telp_number" form:"telp_number"`
		Email      string                `json:"email" form:"email"`
		Image      *multipart.FileHeader `json:"image" form:"image"`
		Password   string                `json:"password" form:"password"`
		PersonID   uuid.UUID             `json:"person_id" form:"person_id"`
	}

	UserResponse struct {
		ID                        uuid.UUID      `json:"id"`
		Email                     string         `json:"email"`
		ImageUrl                  string         `json:"image_url"`
		IsActive                  bool           `json:"is_active"`
		HasChangedDefaultPassword bool           `json:"has_changed_default_password"`
		LastLoginAt               *time.Time     `json:"last_login_at"`
		PersonID                  uuid.UUID      `json:"person_id"`
		Person                    PersonResponse `json:"person"`
		CreatedAt                 time.Time      `json:"created_at"`
		UpdatedAt                 time.Time      `json:"updated_at"`
	}

	UserPaginationResponse struct {
		Data []UserResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User `json:"users"`
		PaginationResponse
	}

	UserUpdateRequest struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		ImageUrl string `json:"image_url" form:"image_url"`
	}

	UserUpdateResponse struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		ImageUrl   string `json:"image_url"`
		IsActive   bool   `json:"is_active"`
	}




	UserLoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	// UserLoginResponse struct {
	// 	Token  string               `json:"token"`
	// 	Email  string               `json:"email"`
	// 	Person SimplePersonResponse `json:"person"`
	// }

	UserLoginResponse struct {
		Token                   string                       `json:"token"`
		Pelayanan               []PersonHasPelayananResponse `json:"pelayanan"`
		Nama                    string                       `json:"nama"`
		ImageUrl                string                       `json:"image_url"`
		IsFirstTimeLogin        bool                         `json:"is_first_time_login"`
		RequiresPasswordSetup   bool                         `json:"requires_password_setup"`
		DefaultPasswordHint     string                       `json:"default_password_hint,omitempty"`
		ExpiredAt               time.Time                    `json:"expired_at"`
	}


	UserRequest struct {
		Email      string    `json:"email" binding:"required,email"`
		Password   string    `json:"password" binding:"required,min=6"`
		ImageUrl   string    `json:"image_url"`
		PersonID   uuid.UUID `json:"person_id" binding:"required"`
	}

	RegisterRequest struct {
		Email    string    `json:"email" binding:"required,email"`
		Password string    `json:"password" binding:"required,min=6"`
		ImageUrl string    `json:"image_url"`
		PersonID uuid.UUID `json:"person_id" binding:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	LoginResponse struct {
		Token      string                       `json:"token"`
		Pelayanan  []PersonHasPelayananResponse `json:"pelayanan"`
		Nama       string                       `json:"nama"`
		ImageUrl   string                       `json:"image_url"`
		ExpiredAt  time.Time                    `json:"expired_at"`
	}


	PasswordSetupRequest struct {
		Action      string `json:"action" binding:"required,oneof=change keep"` // "change" or "keep"
		NewPassword string `json:"new_password,omitempty"`
	}

	PasswordSetupResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)
