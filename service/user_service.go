package service

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"strings"
	"time"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/helpers"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserService interface {
		Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
		GetByEmail(ctx context.Context, email string) (dto.UserResponse, error)
		SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
		VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
		Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
		Delete(ctx context.Context, userId string) error
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		UploadProfileImage(ctx context.Context, file *multipart.FileHeader) (string, error)
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

const (
	LOCAL_URL          = "http://localhost:3000"
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

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

	imageUrl, err := s.UploadProfileImage(ctx, req.Image)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUploadProfileImage
	}

	user := entity.User{
		ID:         uuid.New(),
		Email:      req.Email,
		Password:   string(hashedPassword),
		ImageUrl:   imageUrl,
		IsVerified: false,
		PersonID:   req.PersonID,
		Person:     *person,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	draftEmail, err := makeVerificationEmail(userReg.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	err = utils.SendMail(userReg.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:         userReg.ID,
		Email:      userReg.Email,
		ImageUrl:   userReg.ImageUrl,
		IsVerified: userReg.IsVerified,
		PersonID:   userReg.PersonID,
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

func makeVerificationEmail(receiverEmail string) (map[string]string, error) {
	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := receiverEmail + "_" + expired
	token, err := utils.AESEncrypt(plainText)
	if err != nil {
		return nil, err
	}

	verifyLink := LOCAL_URL + "/" + VERIFY_EMAIL_ROUTE + "?token=" + token

	readHtml, err := os.ReadFile("utils/email-template/base_mail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email  string
		Verify string
	}{
		Email:  receiverEmail,
		Verify: verifyLink,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Cakno - Go Gin Template",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.ErrEmailNotFound
	}

	draftEmail, err := makeVerificationEmail(user.Email)
	if err != nil {
		return err
	}

	err = utils.SendMail(user.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if expiredTime.Sub(now) < 0 {
		return dto.VerifyEmailResponse{
			Email:      email,
			IsVerified: false,
		}, dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	if user.IsVerified {
		return dto.VerifyEmailResponse{}, dto.ErrAccountAlreadyVerified
	}

	err = s.userRepo.Update(ctx, &entity.User{
		ID:         user.ID,
		IsVerified: true,
	})
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
	}

	return dto.VerifyEmailResponse{
		Email:      email,
		IsVerified: true,
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
			ID:         user.ID,
			Email:      user.Email,
			ImageUrl:   user.ImageUrl,
			IsVerified: user.IsVerified,
			PersonID:   user.PersonID,
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
		ID:         user.ID,
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		IsVerified: user.IsVerified,
		PersonID:   user.PersonID,
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
		ID:         user.ID,
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		IsVerified: user.IsVerified,
		PersonID:   user.PersonID,
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
		ID:         user.ID.String(),
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		IsVerified: user.IsVerified,
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
	check, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrEmailNotFound
	}

	if !check.IsVerified {
		return dto.UserLoginResponse{}, dto.ErrAccountNotVerified
	}

	checkPassword, err := helpers.CheckPassword(check.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

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
			IsPic:       p.IsPic,
		})
	}

	token := s.jwtService.GenerateToken(check.ID.String(), check.Email, 24*3)

	return dto.UserLoginResponse{
		Token:      token,
		Pelayanan:  pelayananResponses,
		Nama:       check.Person.Nama,
		ImageUrl:   check.ImageUrl,
		IsVerified: check.IsVerified,
		ExpiredAt:  time.Now().Add(time.Hour * 24 * 3),
	}, nil
}

func (s *userService) UploadProfileImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	return s.documentService.UploadImage(file)
}
