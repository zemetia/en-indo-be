package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type EventParticipantService interface {
	RegisterParticipant(req *dto.CreateEventParticipantRequest) (*dto.EventParticipantResponse, error)
	BulkRegisterParticipants(req *dto.BulkRegisterParticipantsRequest) ([]dto.EventParticipantResponse, error)
	UpdateParticipant(id uuid.UUID, req *dto.UpdateEventParticipantRequest) (*dto.EventParticipantResponse, error)
	DeleteParticipant(id uuid.UUID) error
	CheckInParticipant(id uuid.UUID, req *dto.CheckInParticipantRequest) (*dto.EventParticipantResponse, error)
	BulkCheckIn(req *dto.BulkCheckInRequest) ([]dto.EventParticipantResponse, error)
	GetParticipant(id uuid.UUID) (*dto.EventParticipantResponse, error)
	GetEventParticipants(eventID uuid.UUID, occurrenceDate time.Time, filterReq *dto.EventParticipantFilterRequest) (*dto.EventParticipantListResponse, error)
	GetAttendanceReport(eventID uuid.UUID, occurrenceDate time.Time) (*dto.AttendanceReportResponse, error)
	ListParticipants(filterReq *dto.EventParticipantFilterRequest) (*dto.EventParticipantListResponse, error)
	ScanQRCheckIn(personID uuid.UUID, req *dto.QRScanCheckInRequest) (*dto.QRScanCheckInResponse, error)
}

type eventParticipantService struct {
	participantRepo repository.EventParticipantRepository
	eventRepo       repository.EventRepository
}

func NewEventParticipantService(
	participantRepo repository.EventParticipantRepository,
	eventRepo repository.EventRepository,
) EventParticipantService {
	return &eventParticipantService{
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
	}
}

func (s *eventParticipantService) RegisterParticipant(req *dto.CreateEventParticipantRequest) (*dto.EventParticipantResponse, error) {
	// Validate event exists
	_, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// Validate participant type
	participantType := entity.ParticipantType(req.ParticipantType)
	if participantType != entity.ParticipantTypePerson && participantType != entity.ParticipantTypeVisitor {
		return nil, fmt.Errorf("invalid participant type: %s", req.ParticipantType)
	}

	// Check for duplicate registration
	existing, _ := s.participantRepo.GetByEventAndParticipant(
		req.EventID,
		occurrenceDate,
		participantType,
		req.ParticipantID,
	)
	if existing != nil {
		return nil, fmt.Errorf("participant already registered for this event occurrence")
	}

	// Create participant
	participant := &entity.EventParticipant{
		EventID:            req.EventID,
		OccurrenceDate:     occurrenceDate,
		ParticipantType:    participantType,
		ParticipantID:      req.ParticipantID,
		RegistrationStatus: entity.RegistrationStatusRegistered,
		AttendanceStatus:   entity.AttendanceStatusAbsent,
		Notes:              req.Notes,
	}

	if err := s.participantRepo.Create(participant); err != nil {
		// Check for duplicate key error from database
		errMsg := err.Error()
		if strings.Contains(errMsg, "Duplicate entry") || strings.Contains(errMsg, "1062") {
			return nil, fmt.Errorf("participant already registered for this event occurrence")
		}
		return nil, fmt.Errorf("failed to register participant: %w", err)
	}

	// Reload to get associations
	return s.GetParticipant(participant.ID)
}

func (s *eventParticipantService) BulkRegisterParticipants(req *dto.BulkRegisterParticipantsRequest) ([]dto.EventParticipantResponse, error) {
	// Validate event exists
	if _, err := s.eventRepo.GetByID(req.EventID); err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	var participants []entity.EventParticipant

	// Add persons
	for _, personID := range req.PersonIDs {
		// Check for duplicate
		existing, _ := s.participantRepo.GetByEventAndParticipant(
			req.EventID,
			occurrenceDate,
			entity.ParticipantTypePerson,
			personID,
		)
		if existing != nil {
			continue // Skip already registered
		}

		participants = append(participants, entity.EventParticipant{
			EventID:            req.EventID,
			OccurrenceDate:     occurrenceDate,
			ParticipantType:    entity.ParticipantTypePerson,
			ParticipantID:      personID,
			RegistrationStatus: entity.RegistrationStatusRegistered,
			AttendanceStatus:   entity.AttendanceStatusAbsent,
		})
	}

	// Add visitors
	for _, visitorID := range req.VisitorIDs {
		// Check for duplicate
		existing, _ := s.participantRepo.GetByEventAndParticipant(
			req.EventID,
			occurrenceDate,
			entity.ParticipantTypeVisitor,
			visitorID,
		)
		if existing != nil {
			continue // Skip already registered
		}

		participants = append(participants, entity.EventParticipant{
			EventID:            req.EventID,
			OccurrenceDate:     occurrenceDate,
			ParticipantType:    entity.ParticipantTypeVisitor,
			ParticipantID:      visitorID,
			RegistrationStatus: entity.RegistrationStatusRegistered,
			AttendanceStatus:   entity.AttendanceStatusAbsent,
		})
	}

	if len(participants) == 0 {
		return []dto.EventParticipantResponse{}, nil
	}

	// Bulk create - database will handle foreign key validation
	if err := s.participantRepo.BulkCreate(participants); err != nil {
		return nil, fmt.Errorf("failed to bulk register participants: %w", err)
	}

	// Get all registered participants
	filters := &dto.EventParticipantFilterRequest{
		EventID:        &req.EventID,
		OccurrenceDate: &req.OccurrenceDate,
	}
	result, err := s.GetEventParticipants(req.EventID, occurrenceDate, filters)
	if err != nil {
		return []dto.EventParticipantResponse{}, nil
	}

	return result.Participants, nil
}

func (s *eventParticipantService) UpdateParticipant(id uuid.UUID, req *dto.UpdateEventParticipantRequest) (*dto.EventParticipantResponse, error) {
	participant, err := s.participantRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}

	// Update fields
	if req.RegistrationStatus != nil {
		participant.RegistrationStatus = entity.RegistrationStatus(*req.RegistrationStatus)
	}
	if req.AttendanceStatus != nil {
		participant.AttendanceStatus = entity.AttendanceStatus(*req.AttendanceStatus)
	}
	if req.Notes != nil {
		participant.Notes = *req.Notes
	}

	if err := s.participantRepo.Update(participant); err != nil {
		return nil, fmt.Errorf("failed to update participant: %w", err)
	}

	return s.GetParticipant(id)
}

func (s *eventParticipantService) DeleteParticipant(id uuid.UUID) error {
	if _, err := s.participantRepo.GetByID(id); err != nil {
		return fmt.Errorf("participant not found: %w", err)
	}

	return s.participantRepo.Delete(id)
}

func (s *eventParticipantService) CheckInParticipant(id uuid.UUID, req *dto.CheckInParticipantRequest) (*dto.EventParticipantResponse, error) {
	participant, err := s.participantRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}

	// Mark as present
	checkInMethod := entity.CheckInMethod(req.CheckInMethod)
	participant.MarkAsPresent(checkInMethod, req.RecordedBy)

	// Store QR scan timestamp if provided (for qr_scan method)
	if req.QRScanTimestamp != nil {
		participant.QRScanTimestamp = req.QRScanTimestamp
	}

	if req.Notes != "" {
		participant.Notes = req.Notes
	}

	if err := s.participantRepo.Update(participant); err != nil {
		return nil, fmt.Errorf("failed to check in participant: %w", err)
	}

	return s.GetParticipant(id)
}

func (s *eventParticipantService) BulkCheckIn(req *dto.BulkCheckInRequest) ([]dto.EventParticipantResponse, error) {
	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	var responses []dto.EventParticipantResponse

	for _, participantID := range req.ParticipantIDs {
		participant, err := s.participantRepo.GetByID(participantID)
		if err != nil {
			continue // Skip not found participants
		}

		// Verify participant belongs to the event and occurrence
		if participant.EventID != req.EventID || !participant.OccurrenceDate.Equal(occurrenceDate) {
			continue
		}

		// Mark as present
		checkInMethod := entity.CheckInMethod(req.CheckInMethod)
		participant.MarkAsPresent(checkInMethod, req.RecordedBy)

		if err := s.participantRepo.Update(participant); err != nil {
			continue
		}

		response, _ := s.GetParticipant(participantID)
		if response != nil {
			responses = append(responses, *response)
		}
	}

	return responses, nil
}

func (s *eventParticipantService) GetParticipant(id uuid.UUID) (*dto.EventParticipantResponse, error) {
	participant, err := s.participantRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}

	// Load participant data
	peopleMap, visitorsMap, err := s.participantRepo.LoadParticipantsData([]entity.EventParticipant{*participant})
	if err != nil {
		return nil, fmt.Errorf("failed to load participant data: %w", err)
	}

	// Attach the participant details
	if participant.ParticipantType == entity.ParticipantTypePerson {
		participant.Person = peopleMap[participant.ParticipantID]
	} else if participant.ParticipantType == entity.ParticipantTypeVisitor {
		participant.Visitor = visitorsMap[participant.ParticipantID]
	}

	return s.entityToResponse(participant), nil
}

func (s *eventParticipantService) GetEventParticipants(
	eventID uuid.UUID,
	occurrenceDate time.Time,
	filterReq *dto.EventParticipantFilterRequest,
) (*dto.EventParticipantListResponse, error) {
	filters := s.buildFilters(filterReq)

	participants, total, err := s.participantRepo.GetByEventAndOccurrence(eventID, occurrenceDate, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	// Batch load all participant data
	peopleMap, visitorsMap, err := s.participantRepo.LoadParticipantsData(participants)
	if err != nil {
		return nil, fmt.Errorf("failed to load participants data: %w", err)
	}

	// Attach participant details to each participant
	for i := range participants {
		if participants[i].ParticipantType == entity.ParticipantTypePerson {
			participants[i].Person = peopleMap[participants[i].ParticipantID]
		} else if participants[i].ParticipantType == entity.ParticipantTypeVisitor {
			participants[i].Visitor = visitorsMap[participants[i].ParticipantID]
		}
	}

	// Get stats
	stats, err := s.participantRepo.GetAttendanceStats(eventID, occurrenceDate)
	if err != nil {
		stats = &repository.AttendanceStats{}
	}

	responses := make([]dto.EventParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = *s.entityToResponse(&p)
	}

	// Calculate attendance rate
	attendanceRate := 0.0
	if stats.TotalRegistered > 0 {
		attendanceRate = float64(stats.TotalPresent) / float64(stats.TotalRegistered) * 100
	}

	return &dto.EventParticipantListResponse{
		Participants: responses,
		TotalCount:   int(total),
		Page:         filterReq.Page,
		Limit:        filterReq.Limit,
		Stats: dto.AttendanceStats{
			TotalRegistered: stats.TotalRegistered,
			TotalPresent:    stats.TotalPresent,
			TotalAbsent:     stats.TotalAbsent,
			TotalExcused:    stats.TotalExcused,
			TotalCancelled:  stats.TotalCancelled,
			TotalNoShow:     stats.TotalNoShow,
			AttendanceRate:  attendanceRate,
		},
	}, nil
}

func (s *eventParticipantService) GetAttendanceReport(eventID uuid.UUID, occurrenceDate time.Time) (*dto.AttendanceReportResponse, error) {
	// Get event
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Get participants
	filters := repository.EventParticipantFilters{}
	participants, _, err := s.participantRepo.GetByEventAndOccurrence(eventID, occurrenceDate, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	// Batch load all participant data
	peopleMap, visitorsMap, err := s.participantRepo.LoadParticipantsData(participants)
	if err != nil {
		return nil, fmt.Errorf("failed to load participants data: %w", err)
	}

	// Attach participant details to each participant
	for i := range participants {
		if participants[i].ParticipantType == entity.ParticipantTypePerson {
			participants[i].Person = peopleMap[participants[i].ParticipantID]
		} else if participants[i].ParticipantType == entity.ParticipantTypeVisitor {
			participants[i].Visitor = visitorsMap[participants[i].ParticipantID]
		}
	}

	// Get stats
	stats, err := s.participantRepo.GetAttendanceStats(eventID, occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	responses := make([]dto.EventParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = *s.entityToResponse(&p)
	}

	// Calculate attendance rate
	attendanceRate := 0.0
	if stats.TotalRegistered > 0 {
		attendanceRate = float64(stats.TotalPresent) / float64(stats.TotalRegistered) * 100
	}

	return &dto.AttendanceReportResponse{
		EventID:        eventID,
		EventTitle:     event.Title,
		OccurrenceDate: occurrenceDate,
		Stats: dto.AttendanceStats{
			TotalRegistered: stats.TotalRegistered,
			TotalPresent:    stats.TotalPresent,
			TotalAbsent:     stats.TotalAbsent,
			TotalExcused:    stats.TotalExcused,
			TotalCancelled:  stats.TotalCancelled,
			TotalNoShow:     stats.TotalNoShow,
			AttendanceRate:  attendanceRate,
		},
		Participants: responses,
		GeneratedAt:  time.Now(),
	}, nil
}

func (s *eventParticipantService) ListParticipants(filterReq *dto.EventParticipantFilterRequest) (*dto.EventParticipantListResponse, error) {
	filters := s.buildFilters(filterReq)

	participants, total, err := s.participantRepo.List(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}

	// Batch load all participant data
	peopleMap, visitorsMap, err := s.participantRepo.LoadParticipantsData(participants)
	if err != nil {
		return nil, fmt.Errorf("failed to load participants data: %w", err)
	}

	// Attach participant details to each participant
	for i := range participants {
		if participants[i].ParticipantType == entity.ParticipantTypePerson {
			participants[i].Person = peopleMap[participants[i].ParticipantID]
		} else if participants[i].ParticipantType == entity.ParticipantTypeVisitor {
			participants[i].Visitor = visitorsMap[participants[i].ParticipantID]
		}
	}

	responses := make([]dto.EventParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = *s.entityToResponse(&p)
	}

	return &dto.EventParticipantListResponse{
		Participants: responses,
		TotalCount:   int(total),
		Page:         filterReq.Page,
		Limit:        filterReq.Limit,
		Stats:        dto.AttendanceStats{}, // No stats for general list
	}, nil
}

func (s *eventParticipantService) ScanQRCheckIn(personID uuid.UUID, req *dto.QRScanCheckInRequest) (*dto.QRScanCheckInResponse, error) {
	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// Get event details
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Find participant (person type) by event, occurrence date, and person ID
	participant, err := s.participantRepo.GetByEventAndParticipant(
		req.EventID,
		occurrenceDate,
		entity.ParticipantTypePerson,
		personID,
	)

	// If participant not found, return "not_participant" status
	if err != nil {
		return &dto.QRScanCheckInResponse{
			Status:     "not_participant",
			Message:    "Kamu tidak berpartisipasi pada acara ini",
			EventTitle: event.Title,
		}, nil
	}

	// If already checked in, return "already_checked_in" status
	if participant.CheckInTime != nil {
		// Load participant data for response
		peopleMap, _, err := s.participantRepo.LoadParticipantsData([]entity.EventParticipant{*participant})
		if err == nil && participant.ParticipantType == entity.ParticipantTypePerson {
			participant.Person = peopleMap[participant.ParticipantID]
		}

		participantResponse := s.entityToResponse(participant)

		return &dto.QRScanCheckInResponse{
			Status:      "already_checked_in",
			Message:     "Kamu telah presensi",
			EventTitle:  event.Title,
			CheckInTime: participant.CheckInTime,
			Participant: participantResponse,
		}, nil
	}

	// Perform check-in
	checkInMethod := entity.CheckInMethodQRScan
	now := time.Now()
	participant.MarkAsPresent(checkInMethod, nil) // No recorded_by for self check-in
	participant.QRScanTimestamp = &now

	if err := s.participantRepo.Update(participant); err != nil {
		return nil, fmt.Errorf("failed to check in participant: %w", err)
	}

	// Reload participant to get complete data
	participantResponse, err := s.GetParticipant(participant.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant details: %w", err)
	}

	// Get participant name for success message
	participantName := ""
	if participantResponse.ParticipantDetails.Person != nil {
		participantName = participantResponse.ParticipantDetails.Person.Nama
	}

	return &dto.QRScanCheckInResponse{
		Status:      "success",
		Message:     fmt.Sprintf("Berhasil presensi pada acara %s pada %s terimakasih %s", event.Title, participant.CheckInTime.Format("15:04"), participantName),
		EventTitle:  event.Title,
		CheckInTime: participant.CheckInTime,
		Participant: participantResponse,
	}, nil
}

// Helper methods

func (s *eventParticipantService) entityToResponse(participant *entity.EventParticipant) *dto.EventParticipantResponse {
	response := &dto.EventParticipantResponse{
		ID:                 participant.ID,
		EventID:            participant.EventID,
		OccurrenceDate:     participant.OccurrenceDate,
		ParticipantType:    string(participant.ParticipantType),
		ParticipantID:      participant.ParticipantID,
		RegistrationStatus: string(participant.RegistrationStatus),
		AttendanceStatus:   string(participant.AttendanceStatus),
		CheckInTime:        participant.CheckInTime,
		CheckOutTime:       participant.CheckOutTime,
		QRScanTimestamp:    participant.QRScanTimestamp,
		Notes:              participant.Notes,
		RecordedBy:         participant.RecordedBy,
		CreatedAt:          participant.CreatedAt,
		UpdatedAt:          participant.UpdatedAt,
	}

	if participant.CheckInMethod != nil {
		method := string(*participant.CheckInMethod)
		response.CheckInMethod = &method
	}

	if participant.Event.ID != uuid.Nil {
		response.EventTitle = participant.Event.Title
	}

	// Build participant details
	participantDetails := dto.ParticipantDetails{
		Type: string(participant.ParticipantType),
	}

	if participant.ParticipantType == entity.ParticipantTypePerson && participant.Person != nil {
		participantDetails.Person = &dto.PersonSummary{
			ID:           participant.Person.ID,
			Nama:         participant.Person.Nama,
			Email:        participant.Person.Email,
			NomorTelepon: participant.Person.NomorTelepon,
			ChurchID:     participant.Person.ChurchID,
		}
	} else if participant.ParticipantType == entity.ParticipantTypeVisitor && participant.Visitor != nil {
		participantDetails.Visitor = &dto.VisitorSummary{
			ID:          participant.Visitor.ID,
			Name:        participant.Visitor.Name,
			PhoneNumber: participant.Visitor.PhoneNumber,
			IGUsername:  participant.Visitor.IGUsername,
		}
	}

	response.ParticipantDetails = participantDetails

	// Add recorded by details
	if participant.RecordedByPerson != nil && participant.RecordedByPerson.ID != uuid.Nil {
		response.RecordedByDetails = &dto.PersonSummary{
			ID:           participant.RecordedByPerson.ID,
			Nama:         participant.RecordedByPerson.Nama,
			Email:        participant.RecordedByPerson.Email,
			NomorTelepon: participant.RecordedByPerson.NomorTelepon,
			ChurchID:     participant.RecordedByPerson.ChurchID,
		}
	}

	return response
}

func (s *eventParticipantService) buildFilters(filterReq *dto.EventParticipantFilterRequest) repository.EventParticipantFilters {
	filters := repository.EventParticipantFilters{}

	if filterReq.EventID != nil {
		filters.EventID = filterReq.EventID
	}
	if filterReq.OccurrenceDate != nil {
		occurrenceDate, err := time.Parse("2006-01-02", *filterReq.OccurrenceDate)
		if err == nil {
			filters.OccurrenceDate = &occurrenceDate
		}
	}
	if filterReq.ParticipantType != nil {
		participantType := entity.ParticipantType(*filterReq.ParticipantType)
		filters.ParticipantType = &participantType
	}
	if filterReq.RegistrationStatus != nil {
		registrationStatus := entity.RegistrationStatus(*filterReq.RegistrationStatus)
		filters.RegistrationStatus = &registrationStatus
	}
	if filterReq.AttendanceStatus != nil {
		attendanceStatus := entity.AttendanceStatus(*filterReq.AttendanceStatus)
		filters.AttendanceStatus = &attendanceStatus
	}

	// Set defaults
	if filterReq.Limit == 0 {
		filters.Limit = 50
	} else {
		filters.Limit = filterReq.Limit
	}

	if filterReq.Page > 0 {
		filters.Offset = (filterReq.Page - 1) * filters.Limit
	}

	return filters
}
