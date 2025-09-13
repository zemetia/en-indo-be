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
	GetByPICLifegroupChurches(ctx context.Context, personID uuid.UUID) ([]dto.SimplePersonResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.PersonRequest) (*dto.PersonResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type personService struct {
	personRepository    repository.PersonRepository
	churchRepository    repository.ChurchRepository
	kabupatenRepository repository.KabupatenRepository
	lifeGroupRepository repository.LifeGroupRepository
	pelayananService    PelayananService
}

func NewPersonService(personRepository repository.PersonRepository, churchRepository repository.ChurchRepository, kabupatenRepository repository.KabupatenRepository, lifeGroupRepository repository.LifeGroupRepository, pelayananService PelayananService) PersonService {
	return &personService{
		personRepository:    personRepository,
		churchRepository:    churchRepository,
		kabupatenRepository: kabupatenRepository,
		lifeGroupRepository: lifeGroupRepository,
		pelayananService:    pelayananService,
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

	return s.toResponse(ctx, person), nil
}

func (s *personService) GetByChurchID(ctx context.Context, churchID uuid.UUID) ([]dto.PersonResponse, error) {
	persons, err := s.personRepository.GetByChurchID(ctx, churchID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PersonResponse
	for _, person := range persons {
		responses = append(responses, *s.toResponse(ctx, &person))
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
		responses = append(responses, *s.toResponse(ctx, &person))
	}

	return responses, nil
}

func (s *personService) GetByUserID(ctx context.Context, userID uuid.UUID) (*dto.PersonResponse, error) {
	person, err := s.personRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(ctx, person), nil
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

func (s *personService) GetByPICLifegroupChurches(ctx context.Context, personID uuid.UUID) ([]dto.SimplePersonResponse, error) {
	// First, get the pelayanan assignments of the requesting person to find which churches they are PIC Lifegroup for
	assignments, err := s.pelayananService.GetMyPelayanan(ctx, personID)
	if err != nil {
		return nil, err
	}
	
	// Find churches where the person is PIC Lifegroup 
	// Check by pelayanan name instead of hardcoded ID
	var churchIDs []uuid.UUID
	
	for _, assignment := range assignments {
		if assignment.Pelayanan == "PIC Lifegroup" && assignment.IsPic {
			churchIDs = append(churchIDs, assignment.ChurchID)
		}
	}
	
	if len(churchIDs) == 0 {
		return []dto.SimplePersonResponse{}, nil
	}
	
	// Get all persons from those churches
	var allPersons []dto.SimplePersonResponse
	for _, churchID := range churchIDs {
		persons, err := s.personRepository.GetByChurchID(ctx, churchID)
		if err != nil {
			continue // Skip this church if error, continue with others
		}
		
		// Convert to SimplePersonResponse
		for _, person := range persons {
			allPersons = append(allPersons, dto.SimplePersonResponse{
				ID:           person.ID,
				Nama:         person.Nama,
				Gender:       person.Gender,
				Alamat:       person.Alamat,
				Church:       person.Church.Name,
				TanggalLahir: person.TanggalLahir.Format("2006-01-02"),
				Email:        person.Email,
				NomorTelepon: person.NomorTelepon,
				IsAktif:      person.IsAktif,
			})
		}
	}
	
	return allPersons, nil
}

func (s *personService) toResponse(ctx context.Context, person *entity.Person) *dto.PersonResponse {
	if person == nil {
		return nil
	}

	// Format dates safely
	tanggalLahir := ""
	if !person.TanggalLahir.IsZero() {
		tanggalLahir = person.TanggalLahir.Format("2006-01-02")
	}

	tanggalPerkawinan := ""
	if !person.TanggalPerkawinan.IsZero() {
		tanggalPerkawinan = person.TanggalPerkawinan.Format("2006-01-02")
	}

	// Format timestamps safely
	var createdAt, updatedAt string
	if !person.CreatedAt.IsZero() {
		createdAt = person.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !person.UpdatedAt.IsZero() {
		updatedAt = person.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	// Get lifegroups safely
	var lifeGroups []dto.LifeGroupSimpleResponse
	for _, lg := range person.LifeGroups {
		lifeGroups = append(lifeGroups, dto.LifeGroupSimpleResponse{
			ID:   lg.ID,
			Name: lg.Name,
		})
	}

	// Get church name safely
	var churchName string
	if person.ChurchID != uuid.Nil && person.Church.Name != "" {
		churchName = person.Church.Name
	}

	// Get kabupaten name safely
	var kabupatenName string
	if person.KabupatenID > 0 && person.Kabupaten.Name != "" {
		kabupatenName = person.Kabupaten.Name
	}

	// Get pelayanan safely
	var pelayananResponses []dto.PersonHasPelayananResponse
	if s.pelayananService != nil {
		if pelayanan, err := s.pelayananService.GetMyPelayanan(ctx, person.ID); err == nil {
			pelayananResponses = pelayanan
		}
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
		Church:            churchName,
		KabupatenID:       person.KabupatenID,
		Kabupaten:         kabupatenName,
		LifeGroups:        lifeGroups,
		Pelayanan:         pelayananResponses,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
