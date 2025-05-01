package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type PersonService interface {
	Create(ctx context.Context, req *dto.PersonRequest) (*dto.PersonResponse, error)
	GetAll(ctx context.Context) ([]dto.SimplePersonResponse, error)
	Search(ctx context.Context, search *dto.PersonSearchDto) ([]dto.SimplePersonResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PersonResponse, error)
	GetByChurchID(ctx context.Context, churchID uuid.UUID) ([]dto.PersonResponse, error)
	GetByKabupatenID(ctx context.Context, kabupatenID uuid.UUID) ([]dto.PersonResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.PersonResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.PersonRequest) (*dto.PersonResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type personService struct {
	personRepository    repository.PersonRepository
	churchRepository    repository.ChurchRepository
	kabupatenRepository repository.KabupatenRepository
	lifeGroupRepository repository.LifeGroupRepository
}

func NewPersonService(personRepository repository.PersonRepository, churchRepository repository.ChurchRepository, kabupatenRepository repository.KabupatenRepository, lifeGroupRepository repository.LifeGroupRepository) PersonService {
	return &personService{
		personRepository:    personRepository,
		churchRepository:    churchRepository,
		kabupatenRepository: kabupatenRepository,
		lifeGroupRepository: lifeGroupRepository,
	}
}

func (s *personService) Create(ctx context.Context, req *dto.PersonRequest) (*dto.PersonResponse, error) {
	person := &entity.Person{
		Nama:              req.Nama,
		NamaLain:          req.NamaLain,
		Gender:            req.Gender,
		TempatLahir:       req.TempatLahir,
		TanggalLahir:      req.TanggalLahir,
		FaseHidup:         req.FaseHidup,
		StatusPerkawinan:  req.StatusPerkawinan,
		NamaPasangan:      req.NamaPasangan,
		PasanganID:        req.PasanganID,
		TanggalPerkawinan: req.TanggalPerkawinan,
		Alamat:            req.Alamat,
		NomorTelepon:      req.NomorTelepon,
		Email:             req.Email,
		Ayah:              req.Ayah,
		Ibu:               req.Ibu,
		Kerinduan:         req.Kerinduan,
		KomitmenBerjemaat: req.KomitmenBerjemaat,
		Status:            req.Status,
		KodeJemaat:        req.KodeJemaat,
		ChurchID:          req.ChurchID,
		KabupatenID:       req.KabupatenID,
	}

	if err := s.personRepository.Create(ctx, person); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, person.ID)
}

func (s *personService) GetAll(ctx context.Context) ([]dto.SimplePersonResponse, error) {
	persons, err := s.personRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.SimplePersonResponse
	for _, person := range persons {
		responses = append(responses, dto.SimplePersonResponse{
			ID:           person.ID,
			Nama:         person.Nama,
			Gender:       person.Gender,
			Alamat:       person.Alamat,
			Church:       person.Church.Name,
			TanggalLahir: person.TanggalLahir.Format("2006-01-02"),
			Email:        person.Email,
			NomorTelepon: person.NomorTelepon,
			IsAktif:      person.IsAktif,
			// responses = append(responses, *s.toResponse(&person))
		})
	}

	return responses, nil
}

func (s *personService) Search(ctx context.Context, search *dto.PersonSearchDto) ([]dto.SimplePersonResponse, error) {
	persons, err := s.personRepository.Search(ctx, search)
	if err != nil {
		return nil, err
	}

	var responses []dto.SimplePersonResponse
	for _, person := range persons {
		responses = append(responses, dto.SimplePersonResponse{
			ID:           person.ID,
			Nama:         person.Nama,
			Gender:       person.Gender,
			Alamat:       person.Alamat,
			Church:       person.Church.Name,
			TanggalLahir: person.TanggalLahir.Format("2006-01-02"),
			Email:        person.Email,
			NomorTelepon: person.NomorTelepon,
			IsAktif:      person.IsAktif,
			// responses = append(responses, *s.toResponse(&person))
		})
	}

	return responses, nil
}

func (s *personService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PersonResponse, error) {
	person, err := s.personRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(person), nil
}

func (s *personService) GetByChurchID(ctx context.Context, churchID uuid.UUID) ([]dto.PersonResponse, error) {
	persons, err := s.personRepository.GetByChurchID(ctx, churchID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, person := range persons {
		responses = append(responses, *s.toResponse(&person))
	}

	return responses, nil
}

func (s *personService) GetByKabupatenID(ctx context.Context, kabupatenID uuid.UUID) ([]dto.PersonResponse, error) {
	persons, err := s.personRepository.GetByKabupatenID(ctx, kabupatenID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, person := range persons {
		responses = append(responses, *s.toResponse(&person))
	}

	return responses, nil
}

func (s *personService) GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.PersonResponse, error) {
	person, err := s.personRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(person), nil
}

func (s *personService) Update(ctx context.Context, id uuid.UUID, req *dto.PersonRequest) (*dto.PersonResponse, error) {
	person, err := s.personRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	person.Nama = req.Nama
	person.NamaLain = req.NamaLain
	person.Gender = req.Gender
	person.TempatLahir = req.TempatLahir
	person.TanggalLahir = req.TanggalLahir
	person.FaseHidup = req.FaseHidup
	person.StatusPerkawinan = req.StatusPerkawinan
	person.NamaPasangan = req.NamaPasangan
	person.PasanganID = req.PasanganID
	person.TanggalPerkawinan = req.TanggalPerkawinan
	person.Alamat = req.Alamat
	person.NomorTelepon = req.NomorTelepon
	person.Email = req.Email
	person.Ayah = req.Ayah
	person.Ibu = req.Ibu
	person.Kerinduan = req.Kerinduan
	person.KomitmenBerjemaat = req.KomitmenBerjemaat
	person.Status = req.Status
	person.KodeJemaat = req.KodeJemaat
	person.ChurchID = req.ChurchID
	person.KabupatenID = req.KabupatenID

	if err := s.personRepository.Update(ctx, person); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *personService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.personRepository.Delete(ctx, id)
}

func (s *personService) toResponse(person *entity.Person) *dto.PersonResponse {
	// Format dates
	tanggalLahir := ""
	if !person.TanggalLahir.IsZero() {
		tanggalLahir = person.TanggalLahir.Format("2006-01-02")
	}

	tanggalPerkawinan := ""
	if !person.TanggalPerkawinan.IsZero() {
		tanggalPerkawinan = person.TanggalPerkawinan.Format("2006-01-02")
	}

	// Get lifegroups
	var lifeGroups []dto.LifeGroupSimpleResponse
	for _, lg := range person.LifeGroups {
		lifeGroups = append(lifeGroups, dto.LifeGroupSimpleResponse{
			ID:   lg.ID,
			Name: lg.Name,
		})
	}

	return &dto.PersonResponse{
		ID:                person.ID,
		Nama:              person.Nama,
		NamaLain:          person.NamaLain,
		Gender:            person.Gender,
		TempatLahir:       person.TempatLahir,
		TanggalLahir:      tanggalLahir,
		FaseHidup:         person.FaseHidup,
		StatusPerkawinan:  person.StatusPerkawinan,
		NamaPasangan:      person.NamaPasangan,
		PasanganID:        person.PasanganID,
		TanggalPerkawinan: tanggalPerkawinan,
		Alamat:            person.Alamat,
		NomorTelepon:      person.NomorTelepon,
		Email:             person.Email,
		Ayah:              person.Ayah,
		Ibu:               person.Ibu,
		Kerinduan:         person.Kerinduan,
		KomitmenBerjemaat: person.KomitmenBerjemaat,
		Status:            person.Status,
		KodeJemaat:        person.KodeJemaat,
		ChurchID:          person.ChurchID,
		Church:            person.Church.Name,
		KabupatenID:       person.KabupatenID,
		Kabupaten:         person.Kabupaten.Name,
		LifeGroups:        lifeGroups,
		CreatedAt:         person.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         person.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
