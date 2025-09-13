package service

import (
	"context"
	"fmt"
	"time"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/helpers"
	"github.com/zemetia/en-indo-be/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	UserService interface {
		Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
		GetByEmail(ctx context.Context, email string) (dto.UserResponse, error)
		GetUserByPersonID(ctx context.Context, personID uuid.UUID) (dto.UserResponse, error)
		Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
		Delete(ctx context.Context, userId string) error
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		UploadProfileImage(ctx context.Context, file *multipart.FileHeader) (string, error)
		CreateUserFromPerson(ctx context.Context, person *entity.Person) (*entity.User, error)
		UpdateUserActivationStatus(ctx context.Context, personID uuid.UUID) error
		ToggleUserActivationStatus(ctx context.Context, personID uuid.UUID, isActive bool) error
		SetupPassword(ctx context.Context, userID uuid.UUID, request dto.PasswordSetupRequest) error
	}

	userService struct {
		userRepo        repository.UserRepository
		personRepo      repository.PersonRepository
		documentService DocumentService
		jwtService      JWTService
	}
)

func NewUserService(userRepo repository.UserRepository, personRepo repository.PersonRepository, documentService DocumentService, jwtService JWTService) UserService {
	return &userService{
		userRepo:        userRepo,
		personRepo:      personRepo,
		documentService: documentService,
		jwtService:      jwtService,
	}
}

func (s *userService) Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	// Cek apakah person exists
	person, err := s.personRepo.GetByID(ctx, req.PersonID)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	var imageUrl string
	if req.Image != nil {
		imageUrl, err = s.UploadProfileImage(ctx, req.Image)
		if err != nil {
			return dto.UserResponse{}, dto.ErrUploadProfileImage
		}
	}

	user := entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
		ImageUrl: imageUrl,
		PersonID: req.PersonID,
		Person:   *person,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:                        userReg.ID,
		Email:                     userReg.Email,
		ImageUrl:                  userReg.ImageUrl,
		IsActive:                  userReg.IsActive,
		HasChangedDefaultPassword: userReg.HasChangedDefaultPassword,
		LastLoginAt:               userReg.LastLoginAt,
		PersonID:                  userReg.PersonID,
		Person: dto.PersonResponse{
			ID:                userReg.Person.ID,
			Nama:              userReg.Person.Nama,
			NamaLain:          userReg.Person.NamaLain,
			Gender:            userReg.Person.Gender,
			TempatLahir:       userReg.Person.TempatLahir,
			TanggalLahir:      userReg.Person.TanggalLahir.Format("2006-01-02"),
			FaseHidup:         userReg.Person.FaseHidup,
			StatusPerkawinan:  userReg.Person.StatusPerkawinan,
			NamaPasangan:      userReg.Person.NamaPasangan,
			PasanganID:        userReg.Person.PasanganID,
			TanggalPerkawinan: userReg.Person.TanggalPerkawinan.Format("2006-01-02"),
			Alamat:            userReg.Person.Alamat,
			NomorTelepon:      userReg.Person.NomorTelepon,
			Email:             userReg.Person.Email,
			Ayah:              userReg.Person.Ayah,
			Ibu:               userReg.Person.Ibu,
			Kerinduan:         userReg.Person.Kerinduan,
			KomitmenBerjemaat: userReg.Person.KomitmenBerjemaat,
			Status:            userReg.Person.Status,
			KodeJemaat:        userReg.Person.KodeJemaat,
			ChurchID:          userReg.Person.ChurchID,
			Church:            userReg.Person.Church.Name,
			KabupatenID:       userReg.Person.KabupatenID,
			Kabupaten:         userReg.Person.Kabupaten.Name,
			CreatedAt:         userReg.Person.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         userReg.Person.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		CreatedAt: userReg.CreatedAt,
		UpdatedAt: userReg.UpdatedAt,
	}, nil
}

func (s *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:                        user.ID,
			Email:                     user.Email,
			ImageUrl:                  user.ImageUrl,
			IsActive:                  user.IsActive,
			HasChangedDefaultPassword: user.HasChangedDefaultPassword,
			LastLoginAt:               user.LastLoginAt,
			PersonID:                  user.PersonID,
			Person: dto.PersonResponse{
				ID:                user.Person.ID,
				Nama:              user.Person.Nama,
				NamaLain:          user.Person.NamaLain,
				Gender:            user.Person.Gender,
				TempatLahir:       user.Person.TempatLahir,
				TanggalLahir:      user.Person.TanggalLahir.Format("2006-01-02"),
				FaseHidup:         user.Person.FaseHidup,
				StatusPerkawinan:  user.Person.StatusPerkawinan,
				NamaPasangan:      user.Person.NamaPasangan,
				PasanganID:        user.Person.PasanganID,
				TanggalPerkawinan: user.Person.TanggalPerkawinan.Format("2006-01-02"),
				Alamat:            user.Person.Alamat,
				NomorTelepon:      user.Person.NomorTelepon,
				Email:             user.Person.Email,
				Ayah:              user.Person.Ayah,
				Ibu:               user.Person.Ibu,
				Kerinduan:         user.Person.Kerinduan,
				KomitmenBerjemaat: user.Person.KomitmenBerjemaat,
				Status:            user.Person.Status,
				KodeJemaat:        user.Person.KodeJemaat,
				ChurchID:          user.Person.ChurchID,
				Church:            user.Person.Church.Name,
				KabupatenID:       user.Person.KabupatenID,
				Kabupaten:         user.Person.Kabupaten.Name,
				CreatedAt:         user.Person.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:         user.Person.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:                        user.ID,
		Email:                     user.Email,
		ImageUrl:                  user.ImageUrl,
		IsActive:                  user.IsActive,
		HasChangedDefaultPassword: user.HasChangedDefaultPassword,
		LastLoginAt:               user.LastLoginAt,
		PersonID:                  user.PersonID,
		Person: dto.PersonResponse{
			ID:                user.Person.ID,
			Nama:              user.Person.Nama,
			NamaLain:          user.Person.NamaLain,
			Gender:            user.Person.Gender,
			TempatLahir:       user.Person.TempatLahir,
			TanggalLahir:      user.Person.TanggalLahir.Format("2006-01-02"),
			FaseHidup:         user.Person.FaseHidup,
			StatusPerkawinan:  user.Person.StatusPerkawinan,
			NamaPasangan:      user.Person.NamaPasangan,
			PasanganID:        user.Person.PasanganID,
			TanggalPerkawinan: user.Person.TanggalPerkawinan.Format("2006-01-02"),
			Alamat:            user.Person.Alamat,
			NomorTelepon:      user.Person.NomorTelepon,
			Email:             user.Person.Email,
			Ayah:              user.Person.Ayah,
			Ibu:               user.Person.Ibu,
			Kerinduan:         user.Person.Kerinduan,
			KomitmenBerjemaat: user.Person.KomitmenBerjemaat,
			Status:            user.Person.Status,
			KodeJemaat:        user.Person.KodeJemaat,
			ChurchID:          user.Person.ChurchID,
			Church:            user.Person.Church.Name,
			KabupatenID:       user.Person.KabupatenID,
			Kabupaten:         user.Person.Kabupaten.Name,
			CreatedAt:         user.Person.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         user.Person.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetByEmail
	}

	return dto.UserResponse{
		ID:                        user.ID,
		Email:                     user.Email,
		ImageUrl:                  user.ImageUrl,
		IsActive:                  user.IsActive,
		HasChangedDefaultPassword: user.HasChangedDefaultPassword,
		LastLoginAt:               user.LastLoginAt,
		PersonID:                  user.PersonID,
		Person: dto.PersonResponse{
			ID:                user.Person.ID,
			Nama:              user.Person.Nama,
			NamaLain:          user.Person.NamaLain,
			Gender:            user.Person.Gender,
			TempatLahir:       user.Person.TempatLahir,
			TanggalLahir:      user.Person.TanggalLahir.Format("2006-01-02"),
			FaseHidup:         user.Person.FaseHidup,
			StatusPerkawinan:  user.Person.StatusPerkawinan,
			NamaPasangan:      user.Person.NamaPasangan,
			PasanganID:        user.Person.PasanganID,
			TanggalPerkawinan: user.Person.TanggalPerkawinan.Format("2006-01-02"),
			Alamat:            user.Person.Alamat,
			NomorTelepon:      user.Person.NomorTelepon,
			Email:             user.Person.Email,
			Ayah:              user.Person.Ayah,
			Ibu:               user.Person.Ibu,
			Kerinduan:         user.Person.Kerinduan,
			KomitmenBerjemaat: user.Person.KomitmenBerjemaat,
			Status:            user.Person.Status,
			KodeJemaat:        user.Person.KodeJemaat,
			ChurchID:          user.Person.ChurchID,
			Church:            user.Person.Church.Name,
			KabupatenID:       user.Person.KabupatenID,
			Kabupaten:         user.Person.Kabupaten.Name,
			CreatedAt:         user.Person.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         user.Person.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUserByPersonID(ctx context.Context, personID uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepo.GetByPersonID(ctx, personID)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	return dto.UserResponse{
		ID:                        user.ID,
		Email:                     user.Email,
		ImageUrl:                  user.ImageUrl,
		IsActive:                  user.IsActive,
		HasChangedDefaultPassword: user.HasChangedDefaultPassword,
		LastLoginAt:               user.LastLoginAt,
		PersonID:                  user.PersonID,
		Person: dto.PersonResponse{
			ID:                user.Person.ID,
			Nama:              user.Person.Nama,
			NamaLain:          user.Person.NamaLain,
			Gender:            user.Person.Gender,
			TempatLahir:       user.Person.TempatLahir,
			TanggalLahir:      user.Person.TanggalLahir.Format("2006-01-02"),
			FaseHidup:         user.Person.FaseHidup,
			StatusPerkawinan:  user.Person.StatusPerkawinan,
			NamaPasangan:      user.Person.NamaPasangan,
			PasanganID:        user.Person.PasanganID,
			TanggalPerkawinan: user.Person.TanggalPerkawinan.Format("2006-01-02"),
			Alamat:            user.Person.Alamat,
			NomorTelepon:      user.Person.NomorTelepon,
			Email:             user.Person.Email,
			Ayah:              user.Person.Ayah,
			Ibu:               user.Person.Ibu,
			Kerinduan:         user.Person.Kerinduan,
			KomitmenBerjemaat: user.Person.KomitmenBerjemaat,
			Status:            user.Person.Status,
			KodeJemaat:        user.Person.KodeJemaat,
			ChurchID:          user.Person.ChurchID,
			Church:            user.Person.Church.Name,
			KabupatenID:       user.Person.KabupatenID,
			Kabupaten:         user.Person.Kabupaten.Name,
			CreatedAt:         user.Person.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         user.Person.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	// Cek apakah email sudah terdaftar (jika email diubah)
	if req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser != nil {
			return dto.UserUpdateResponse{}, dto.ErrEmailAlreadyExists
		}
	}

	// Hash password jika diubah
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.UserUpdateResponse{}, dto.ErrUpdateUser
		}
		user.Password = string(hashedPassword)
	}

	user.Email = req.Email
	user.ImageUrl = req.ImageUrl

	err = s.userRepo.Update(ctx, &user)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		ImageUrl: user.ImageUrl,
		IsActive: user.IsActive,
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}
	err = s.userRepo.Delete(ctx, user.ID)
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (s *userService) Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	// Debug logging for login attempts
	fmt.Printf("[DEBUG] Login attempt for email: %s\n", req.Email)

	check, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil || !flag {
		fmt.Printf("[DEBUG] Email not found: %s\n", req.Email)
		return dto.UserLoginResponse{}, dto.ErrEmailNotFound
	}

	fmt.Printf("[DEBUG] User found - ID: %s, IsActive: %t\n", check.ID, check.IsActive)

	// Check if user is active
	if !check.IsActive {
		fmt.Printf("[DEBUG] User account inactive for email: %s\n", req.Email)
		return dto.UserLoginResponse{}, dto.ErrUserInactive
	}

	// Check if user has active pelayanan assignments
	hasActivePelayanan, err := s.userRepo.HasActivePelayanan(ctx, check.PersonID)
	if err != nil {
		fmt.Printf("[DEBUG] Error checking pelayanan for PersonID: %s, error: %v\n", check.PersonID, err)
		return dto.UserLoginResponse{}, dto.ErrGetPelayanan
	}

	fmt.Printf("[DEBUG] HasActivePelayanan: %t for PersonID: %s\n", hasActivePelayanan, check.PersonID)

	if !hasActivePelayanan {
		// Auto-deactivate user if no pelayanan assignments
		fmt.Printf("[DEBUG] No active pelayanan for PersonID: %s, deactivating user\n", check.PersonID)
		s.userRepo.UpdateActivationStatus(ctx, check.ID, false)
		return dto.UserLoginResponse{}, dto.ErrUserNoPelayanan
	}

	checkPassword, err := helpers.CheckPassword(check.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		fmt.Printf("[DEBUG] Password check failed for email: %s, err: %v\n", req.Email, err)
		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	fmt.Printf("[DEBUG] Password check successful for email: %s\n", req.Email)

	pelayanan, err := s.personRepo.GetPelayananChurchByID(ctx, check.PersonID)
	if err != nil {
		return dto.UserLoginResponse{}, dto.ErrGetPelayanan
	}

	var pelayananResponses []dto.PersonHasPelayananResponse
	for _, p := range pelayanan {
		pelayananResponses = append(pelayananResponses, dto.PersonHasPelayananResponse{
			PelayananID: p.PelayananID,
			Pelayanan:   p.Pelayanan.Pelayanan,
			ChurchID:    p.ChurchID,
			ChurchName:  p.Church.Name,
			IsPic:       p.Pelayanan.IsPic,
		})
	}

	// Detect first-time login and password setup requirements
	isFirstTimeLogin := check.LastLoginAt == nil
	requiresPasswordSetup := !check.HasChangedDefaultPassword

	// Generate password hint if needed
	var defaultPasswordHint string
	if requiresPasswordSetup {
		defaultPasswordHint = "Your password is your birth date in DD/MM/YYYY format"
	}

	// Update last login timestamp
	now := time.Now()
	s.userRepo.Update(ctx, &entity.User{
		ID:          check.ID,
		LastLoginAt: &now,
	})

	token := s.jwtService.GenerateToken(check.ID.String(), check.Email, 24*3)

	return dto.UserLoginResponse{
		Token:                 token,
		Pelayanan:             pelayananResponses,
		Nama:                  check.Person.Nama,
		ImageUrl:              check.ImageUrl,
		IsFirstTimeLogin:      isFirstTimeLogin,
		RequiresPasswordSetup: requiresPasswordSetup,
		DefaultPasswordHint:   defaultPasswordHint,
		ExpiredAt:             time.Now().Add(time.Hour * 24 * 3),
	}, nil
}

func (s *userService) UploadProfileImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.documentService.UploadImage(file)
}

// CreateUserFromPerson creates a new user automatically from person data when assigned to pelayanan
func (s *userService) CreateUserFromPerson(ctx context.Context, person *entity.Person) (*entity.User, error) {
	// Validate person has email
	if person.Email == "" {
		return nil, fmt.Errorf("person must have email to create user account")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, person.Email)
	if err == nil && existingUser != nil {
		return existingUser, nil // User already exists, return existing user
	}

	// Generate password from birth date
	password := helpers.GeneratePasswordFromBirthDate(person.TanggalLahir)

	// Create new user
	user := &entity.User{
		ID:       uuid.New(),
		Email:    person.Email,
		Password: password, // Will be hashed by BeforeCreate hook
		ImageUrl: "",
		IsActive: true,
		PersonID: person.ID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUserActivationStatus updates user activation based on pelayanan assignments
func (s *userService) UpdateUserActivationStatus(ctx context.Context, personID uuid.UUID) error {
	// Get user by person ID
	user, err := s.userRepo.GetByPersonID(ctx, personID)
	if err != nil {
		// If user doesn't exist, skip activation update
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has active pelayanan assignments
	hasActivePelayanan, err := s.userRepo.HasActivePelayanan(ctx, personID)
	if err != nil {
		return fmt.Errorf("failed to check pelayanan assignments: %w", err)
	}

	// Update activation status based on pelayanan assignments
	newActiveStatus := hasActivePelayanan
	if user.IsActive != newActiveStatus {
		if err := s.userRepo.UpdateActivationStatus(ctx, user.ID, newActiveStatus); err != nil {
			return fmt.Errorf("failed to update activation status: %w", err)
		}
	}

	return nil
}

// ToggleUserActivationStatus toggles user activation status by person ID
// Handles three scenarios: 1) No user exists - creates account, 2) User exists but inactive - activates, 3) User active - deactivates
func (s *userService) ToggleUserActivationStatus(ctx context.Context, personID uuid.UUID, isActive bool) error {
	// Get user by person ID
	user, err := s.userRepo.GetByPersonID(ctx, personID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create new account if activating
			if isActive {
				// Get person data to create user account
				person, personErr := s.personRepo.GetByID(ctx, personID)
				if personErr != nil {
					return fmt.Errorf("failed to get person data: %w", personErr)
				}

				// Create user account from person data
				_, createErr := s.CreateUserFromPerson(ctx, person)
				if createErr != nil {
					return fmt.Errorf("failed to create user account: %w", createErr)
				}
				return nil
			} else {
				// Cannot deactivate non-existent account
				return fmt.Errorf("cannot deactivate non-existent user account for person ID: %s", personID)
			}
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// User exists, update activation status
	if err := s.userRepo.UpdateActivationStatus(ctx, user.ID, isActive); err != nil {
		return fmt.Errorf("failed to update activation status: %w", err)
	}

	return nil
}

// SetupPassword handles the password setup flow for first-time users
func (s *userService) SetupPassword(ctx context.Context, userID uuid.UUID, request dto.PasswordSetupRequest) error {
	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	switch request.Action {
	case "change":
		if request.NewPassword == "" {
			return fmt.Errorf("new password is required when action is 'change'")
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Update user with new password and mark as changed
		user.Password = string(hashedPassword)
		user.HasChangedDefaultPassword = true

	case "keep":
		// User chose to keep birth date password, just mark as acknowledged
		user.HasChangedDefaultPassword = true

	default:
		return fmt.Errorf("invalid action: %s", request.Action)
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
