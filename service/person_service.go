package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type PersonService struct {
	personRepo *repository.PersonRepository
}

func NewPersonService(personRepo *repository.PersonRepository) *PersonService {
	return &PersonService{
		personRepo: personRepo,
	}
}

func (s *PersonService) Create(req *dto.PersonRequest) (*dto.PersonResponse, error) {
	person := &entity.Person{
		ID:                uuid.New(),
		Nama:              req.Nama,
		NamaLain:          req.NamaLain,
		Gender:            req.Gender,
		TempatLahir:       req.TempatLahir,
		TanggalLahir:      req.TanggalLahir,
		FaseHidup:         req.FaseHidup,
		StatusPerkawinan:  req.StatusPerkawinan,
		Pasangan:          req.Pasangan,
		TanggalPerkawinan: req.TanggalPerkawinan,
		NomorTelepon:      req.NomorTelepon,
		Email:             req.Email,
		Gereja:            req.Gereja,
		Ayah:              req.Ayah,
		Ibu:               req.Ibu,
		Kerinduan:         req.Kerinduan,
		KomitmenBerjemaat: req.KomitmenBerjemaat,
		DateAdded:         req.DateAdded,
		Status:            req.Status,
		TagListID:         req.TagListID,
		KodeJemaat:        req.KodeJemaat,
		ChurchID:          req.ChurchID,
	}

	if err := s.personRepo.Create(person); err != nil {
		return nil, err
	}

	return s.toResponse(person), nil
}

func (s *PersonService) GetAll() ([]dto.PersonResponse, error) {
	persons, err := s.personRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, person := range persons {
		responses = append(responses, *s.toResponse(&person))
	}

	return responses, nil
}

func (s *PersonService) GetByID(id uuid.UUID) (*dto.PersonResponse, error) {
	person, err := s.personRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(person), nil
}

func (s *PersonService) Update(id uuid.UUID, req *dto.PersonRequest) (*dto.PersonResponse, error) {
	person, err := s.personRepo.GetByID(id)
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
	person.Pasangan = req.Pasangan
	person.TanggalPerkawinan = req.TanggalPerkawinan
	person.NomorTelepon = req.NomorTelepon
	person.Email = req.Email
	person.Gereja = req.Gereja
	person.Ayah = req.Ayah
	person.Ibu = req.Ibu
	person.Kerinduan = req.Kerinduan
	person.KomitmenBerjemaat = req.KomitmenBerjemaat
	person.DateAdded = req.DateAdded
	person.Status = req.Status
	person.TagListID = req.TagListID
	person.KodeJemaat = req.KodeJemaat
	person.ChurchID = req.ChurchID

	if err := s.personRepo.Update(person); err != nil {
		return nil, err
	}

	return s.toResponse(person), nil
}

func (s *PersonService) Delete(id uuid.UUID) error {
	return s.personRepo.Delete(id)
}

func (s *PersonService) GetByChurchID(churchID uuid.UUID) ([]dto.PersonResponse, error) {
	persons, err := s.personRepo.GetByChurchID(churchID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, person := range persons {
		responses = append(responses, *s.toResponse(&person))
	}

	return responses, nil
}

func (s *PersonService) toResponse(person *entity.Person) *dto.PersonResponse {
	return &dto.PersonResponse{
		ID:                person.ID,
		PersonID:          person.PersonID,
		Nama:              person.Nama,
		NamaLain:          person.NamaLain,
		Gender:            person.Gender,
		TempatLahir:       person.TempatLahir,
		TanggalLahir:      person.TanggalLahir,
		FaseHidup:         person.FaseHidup,
		StatusPerkawinan:  person.StatusPerkawinan,
		Pasangan:          person.Pasangan,
		TanggalPerkawinan: person.TanggalPerkawinan,
		NomorTelepon:      person.NomorTelepon,
		Email:             person.Email,
		Gereja:            person.Gereja,
		Ayah:              person.Ayah,
		Ibu:               person.Ibu,
		Kerinduan:         person.Kerinduan,
		KomitmenBerjemaat: person.KomitmenBerjemaat,
		DateAdded:         person.DateAdded,
		Status:            person.Status,
		TagListID:         person.TagListID,
		KodeJemaat:        person.KodeJemaat,
		ChurchID:          person.ChurchID,
		Church:            person.Church,
		UserID:            person.UserID,
		User:              person.User,
		CreatedAt:         person.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         person.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
