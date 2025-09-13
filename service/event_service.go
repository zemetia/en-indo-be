package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type EventService interface {
	CreateEvent(req *dto.CreateEventRequest) (*dto.EventResponse, error)
	GetEvent(id uuid.UUID) (*dto.EventResponse, error)
	UpdateEvent(id uuid.UUID, req *dto.UpdateEventRequest) (*dto.EventResponse, error)
	DeleteEvent(id uuid.UUID) error
	ListEvents(req *dto.EventFilterRequest) (*dto.EventListResponse, error)
	
	// Event creation with PIC assignment
	CreateEventWithPICs(req *dto.CreateEventRequest, createdBy uuid.UUID) (*dto.EventResponse, error)

	// Recurring event operations
	UpdateRecurringEvent(id uuid.UUID, req *dto.UpdateRecurringEventRequest) error
	UpdateSingleOccurrence(id uuid.UUID, req *dto.UpdateOccurrenceRequest) error
	UpdateFutureOccurrences(id uuid.UUID, req *dto.UpdateFutureOccurrencesRequest) error
	DeleteOccurrence(id uuid.UUID, req *dto.DeleteOccurrenceRequest) error
	GetEventOccurrences(id uuid.UUID, req *dto.GetEventOccurrencesRequest) ([]dto.EventOccurrenceResponse, error)
	GetOccurrencesInRange(req *dto.GetEventOccurrencesRequest) ([]dto.EventOccurrenceResponse, error)

	// Validation and utility methods
	ValidateRecurrenceRule(rule *dto.CreateRecurrenceRuleRequest) error
	GetNextOccurrence(id uuid.UUID, after time.Time) (*time.Time, error)
}

type eventService struct {
	eventRepo           repository.EventRepository
	eventPICRepo        repository.EventPICRepository
	recurrenceGenerator *RecurrenceGenerator
}

func NewEventService(eventRepo repository.EventRepository, eventPICRepo repository.EventPICRepository) EventService {
	return &eventService{
		eventRepo:           eventRepo,
		eventPICRepo:        eventPICRepo,
		recurrenceGenerator: NewRecurrenceGenerator(),
	}
}

func (s *eventService) CreateEvent(req *dto.CreateEventRequest) (*dto.EventResponse, error) {
	// Parse dates
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, fmt.Errorf("invalid event date format: %w", err)
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %w", err)
	}

	// Combine date and time
	startDateTime := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	endDateTime := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	// Create event entity
	event := &entity.Event{
		Title:                 req.Title,
		BannerImage:           req.BannerImage,
		Description:           req.Description,
		Capacity:              req.Capacity,
		Type:                  req.Type,
		EventDate:             eventDate,
		EventLocation:         req.EventLocation,
		StartDatetime:         startDateTime,
		EndDatetime:           endDateTime,
		AllDay:                req.AllDay,
		Timezone:              req.Timezone,
		IsPublic:              req.IsPublic,
		DiscipleshipJourneyID: req.DiscipleshipJourneyID,
	}

	// Handle expected participant counts
	if req.ExpectedParticipants != nil {
		event.ExpectedParticipants = *req.ExpectedParticipants
	}
	if req.ExpectedAdults != nil {
		event.ExpectedAdults = *req.ExpectedAdults
	}
	if req.ExpectedYouth != nil {
		event.ExpectedYouth = *req.ExpectedYouth
	}
	if req.ExpectedKids != nil {
		event.ExpectedKids = *req.ExpectedKids
	}

	// Handle recurrence rule
	if req.RecurrenceRule != nil {
		rule, err := s.createRecurrenceRuleEntity(req.RecurrenceRule)
		if err != nil {
			return nil, fmt.Errorf("invalid recurrence rule: %w", err)
		}
		event.RecurrenceRule = rule
	}

	// Set default capacity if not provided
	if event.Capacity == 0 {
		event.Capacity = 99999
	}

	// Create event
	if err := s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return s.entityToResponse(event), nil
}

func (s *eventService) CreateEventWithPICs(req *dto.CreateEventRequest, createdBy uuid.UUID) (*dto.EventResponse, error) {
	// Create the event first
	eventResponse, err := s.CreateEvent(req)
	if err != nil {
		return nil, err
	}
	
	// Assign PICs if provided
	if len(req.EventPICs) > 0 {
		var eventPICs []entity.EventPIC
		
		for _, picReq := range req.EventPICs {
			// Parse dates
			startDate, err := time.Parse("2006-01-02", picReq.StartDate)
			if err != nil {
				return nil, fmt.Errorf("invalid start date format for PIC: %w", err)
			}
			
			var endDate *time.Time
			if picReq.EndDate != nil {
				parsed, err := time.Parse("2006-01-02", *picReq.EndDate)
				if err != nil {
					return nil, fmt.Errorf("invalid end date format for PIC: %w", err)
				}
				endDate = &parsed
			}
			
			eventPIC := entity.EventPIC{
				EventID:           eventResponse.ID,
				PersonID:          picReq.PersonID,
				Role:              picReq.Role,
				Description:       picReq.Description,
				IsActive:          true,
				IsPrimary:         picReq.IsPrimary,
				StartDate:         startDate,
				EndDate:           endDate,
				CanEdit:           picReq.CanEdit,
				CanDelete:         picReq.CanDelete,
				CanAssignPIC:      picReq.CanAssignPIC,
				NotifyOnChanges:   picReq.NotifyOnChanges,
				NotifyOnReminders: picReq.NotifyOnReminders,
			}
			
			eventPICs = append(eventPICs, eventPIC)
		}
		
		// Validate no multiple primary PICs
		primaryCount := 0
		for _, pic := range eventPICs {
			if pic.IsPrimary {
				primaryCount++
			}
		}
		if primaryCount > 1 {
			return nil, fmt.Errorf("cannot assign multiple primary PICs")
		}
		
		// Create PICs
		if err := s.eventPICRepo.CreateMultiple(eventPICs); err != nil {
			return nil, fmt.Errorf("failed to create event PICs: %w", err)
		}
		
		// Create history records
		for _, pic := range eventPICs {
			history := &entity.EventPICHistory{
				EventID:    eventResponse.ID,
				PersonID:   pic.PersonID,
				Action:     entity.EventPICActionAssigned,
				NewRole:    pic.Role,
				ChangedBy:  createdBy,
				Reason:     "Initial PIC assignment",
				ActionDate: time.Now(),
			}
			s.eventPICRepo.CreateHistory(history)
		}
	}
	
	// Reload event with PICs
	return s.GetEvent(eventResponse.ID)
}

func (s *eventService) GetEvent(id uuid.UUID) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return s.entityToResponse(event), nil
}

func (s *eventService) UpdateEvent(id uuid.UUID, req *dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Update fields if provided
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.BannerImage != nil {
		event.BannerImage = *req.BannerImage
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Capacity != nil {
		event.Capacity = *req.Capacity
	}
	if req.Type != nil {
		event.Type = *req.Type
	}
	if req.EventLocation != nil {
		event.EventLocation = *req.EventLocation
	}
	if req.IsPublic != nil {
		event.IsPublic = *req.IsPublic
	}
	if req.DiscipleshipJourneyID != nil {
		event.DiscipleshipJourneyID = req.DiscipleshipJourneyID
	}

	// Handle date/time updates
	if req.EventDate != nil || req.StartTime != nil || req.EndTime != nil {
		eventDateStr := event.EventDate.Format("2006-01-02")
		startTimeStr := event.StartDatetime.Format("15:04")
		endTimeStr := event.EndDatetime.Format("15:04")

		if req.EventDate != nil {
			eventDateStr = *req.EventDate
		}
		if req.StartTime != nil {
			startTimeStr = *req.StartTime
		}
		if req.EndTime != nil {
			endTimeStr = *req.EndTime
		}

		eventDate, err := time.Parse("2006-01-02", eventDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid event date format: %w", err)
		}

		startTime, err := time.Parse("15:04", startTimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %w", err)
		}

		endTime, err := time.Parse("15:04", endTimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}

		event.EventDate = eventDate
		event.StartDatetime = time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
		event.EndDatetime = time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
	}

	if req.AllDay != nil {
		event.AllDay = *req.AllDay
	}
	if req.Timezone != nil {
		event.Timezone = *req.Timezone
	}

	// Handle expected participant count updates
	if req.ExpectedParticipants != nil {
		event.ExpectedParticipants = *req.ExpectedParticipants
	}
	if req.ExpectedAdults != nil {
		event.ExpectedAdults = *req.ExpectedAdults
	}
	if req.ExpectedYouth != nil {
		event.ExpectedYouth = *req.ExpectedYouth
	}
	if req.ExpectedKids != nil {
		event.ExpectedKids = *req.ExpectedKids
	}

	// Update event
	if err := s.eventRepo.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return s.entityToResponse(event), nil
}

func (s *eventService) DeleteEvent(id uuid.UUID) error {
	if err := s.eventRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (s *eventService) ListEvents(req *dto.EventFilterRequest) (*dto.EventListResponse, error) {
	filters := repository.EventFilters{
		Type:     req.Type,
		IsPublic: req.IsPublic,
		Search:   req.Search,
		Limit:    req.Limit,
		Offset:   (req.Page - 1) * req.Limit,
	}

	// Set defaults
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}

	events, total, err := s.eventRepo.List(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = *s.entityToResponse(&event)
	}

	return &dto.EventListResponse{
		Events:     responses,
		TotalCount: int(total),
		Page:       req.Page,
		Limit:      filters.Limit,
	}, nil
}

func (s *eventService) UpdateRecurringEvent(id uuid.UUID, req *dto.UpdateRecurringEventRequest) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event.RecurrenceRule == nil {
		return fmt.Errorf("event is not recurring")
	}

	switch req.UpdateType {
	case dto.UpdateThisEvent:
		return s.updateSingleOccurrenceFromRecurring(event, req)
	case dto.UpdateAllEvents:
		return s.updateEntireSeries(event, req)
	case dto.UpdateFutureEvents:
		return s.updateThisAndFutureOccurrences(event, req)
	default:
		return fmt.Errorf("invalid update type: %s", req.UpdateType)
	}
}

func (s *eventService) DeleteOccurrence(id uuid.UUID, req *dto.DeleteOccurrenceRequest) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event.RecurrenceRule == nil {
		return fmt.Errorf("event is not recurring")
	}

	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return fmt.Errorf("invalid occurrence date format: %w", err)
	}

	switch req.DeleteType {
	case "single":
		return s.skipSingleOccurrence(event, occurrenceDate)
	case "future":
		return s.deleteFutureOccurrences(event, occurrenceDate)
	default:
		return fmt.Errorf("invalid delete type: %s", req.DeleteType)
	}
}

func (s *eventService) GetEventOccurrences(id uuid.UUID, req *dto.GetEventOccurrencesRequest) ([]dto.EventOccurrenceResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	return s.generateOccurrences(event, startDate, endDate)
}

func (s *eventService) GetOccurrencesInRange(req *dto.GetEventOccurrencesRequest) ([]dto.EventOccurrenceResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Get all events that could have occurrences in this range
	events, err := s.eventRepo.GetEventsWithRecurrenceInRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurring events: %w", err)
	}

	var allOccurrences []dto.EventOccurrenceResponse

	for _, event := range events {
		occurrences, err := s.generateOccurrences(&event, startDate, endDate)
		if err != nil {
			continue // Skip events with generation errors
		}
		allOccurrences = append(allOccurrences, occurrences...)
	}

	return allOccurrences, nil
}

// Helper methods

func (s *eventService) createRecurrenceRuleEntity(req *dto.CreateRecurrenceRuleRequest) (*entity.RecurrenceRule, error) {
	// Convert slice fields to JSON strings
	byWeekdayJSON, err := s.sliceToJSON(req.ByWeekday)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ByWeekday to JSON: %w", err)
	}

	byMonthDayJSON, err := s.sliceToJSON(req.ByMonthDay)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ByMonthDay to JSON: %w", err)
	}

	byMonthJSON, err := s.sliceToJSON(req.ByMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ByMonth to JSON: %w", err)
	}

	bySetPosJSON, err := s.sliceToJSON(req.BySetPos)
	if err != nil {
		return nil, fmt.Errorf("failed to convert BySetPos to JSON: %w", err)
	}

	byYearDayJSON, err := s.sliceToJSON(req.ByYearDay)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ByYearDay to JSON: %w", err)
	}

	rule := &entity.RecurrenceRule{
		Frequency:  req.Frequency,
		Interval:   req.Interval,
		ByWeekday:  byWeekdayJSON,
		ByMonthDay: byMonthDayJSON,
		ByMonth:    byMonthJSON,
		BySetPos:   bySetPosJSON,
		WeekStart:  req.WeekStart,
		ByYearDay:  byYearDayJSON,
		Count:      req.Count,
	}

	if rule.Interval == 0 {
		rule.Interval = 1
	}

	if rule.WeekStart == "" {
		rule.WeekStart = "MO" // Default to Monday
	}

	if req.Until != nil {
		until, err := time.Parse("2006-01-02", *req.Until)
		if err != nil {
			return nil, fmt.Errorf("invalid until date format: %w", err)
		}
		rule.Until = &until
	}

	// Validate the rule
	if err := s.recurrenceGenerator.ValidateRecurrenceRule(rule); err != nil {
		return nil, fmt.Errorf("invalid recurrence rule: %w", err)
	}

	return rule, nil
}

func (s *eventService) entityToResponse(event *entity.Event) *dto.EventResponse {
	response := &dto.EventResponse{
		ID:                    event.ID,
		Title:                 event.Title,
		BannerImage:           event.BannerImage,
		Description:           event.Description,
		Capacity:              event.Capacity,
		Type:                  event.Type,
		EventDate:             event.EventDate,
		EventLocation:         event.EventLocation,
		StartDatetime:         event.StartDatetime,
		EndDatetime:           event.EndDatetime,
		AllDay:                event.AllDay,
		Timezone:              event.Timezone,
		IsPublic:              event.IsPublic,
		DiscipleshipJourneyID: event.DiscipleshipJourneyID,
		ExpectedParticipants:  event.ExpectedParticipants,
		ExpectedAdults:        event.ExpectedAdults,
		ExpectedYouth:         event.ExpectedYouth,
		ExpectedKids:          event.ExpectedKids,
		CreatedAt:             event.CreatedAt,
		UpdatedAt:             event.UpdatedAt,
	}
	
	// Convert EventPICs
	if len(event.EventPICs) > 0 {
		eventPICs := make([]dto.EventPICResponse, 0, len(event.EventPICs))
		var primaryPIC *dto.EventPICResponse
		
		for _, pic := range event.EventPICs {
			picResponse := dto.EventPICResponse{
				ID:          pic.ID,
				EventID:     pic.EventID,
				PersonID:    pic.PersonID,
				Person: dto.PersonSummary{
					ID:           pic.Person.ID,
					Nama:         pic.Person.Nama,
					Email:        pic.Person.Email,
					NomorTelepon: pic.Person.NomorTelepon,
					ChurchID:     pic.Person.ChurchID,
				},
				Role:              pic.Role,
				Description:       pic.Description,
				IsActive:          pic.IsActive,
				IsPrimary:         pic.IsPrimary,
				StartDate:         pic.StartDate,
				EndDate:           pic.EndDate,
				CanEdit:           pic.CanEdit,
				CanDelete:         pic.CanDelete,
				CanAssignPIC:      pic.CanAssignPIC,
				NotifyOnChanges:   pic.NotifyOnChanges,
				NotifyOnReminders: pic.NotifyOnReminders,
				CreatedAt:         pic.CreatedAt,
				UpdatedAt:         pic.UpdatedAt,
			}
			
			eventPICs = append(eventPICs, picResponse)
			
			if pic.IsPrimary && pic.IsActive {
				primaryPIC = &picResponse
			}
		}
		
		response.EventPICs = eventPICs
		response.PrimaryPIC = primaryPIC
	}

	if event.RecurrenceRule != nil {
		// Convert JSON strings to slices
		byWeekday, err := s.jsonToStringSlice(event.RecurrenceRule.ByWeekday)
		if err != nil {
			byWeekday = []string{} // Default to empty slice on error
		}

		byMonthDay, err := s.jsonToInt64Slice(event.RecurrenceRule.ByMonthDay)
		if err != nil {
			byMonthDay = []int64{} // Default to empty slice on error
		}

		byMonth, err := s.jsonToInt64Slice(event.RecurrenceRule.ByMonth)
		if err != nil {
			byMonth = []int64{} // Default to empty slice on error
		}

		bySetPos, err := s.jsonToInt64Slice(event.RecurrenceRule.BySetPos)
		if err != nil {
			bySetPos = []int64{} // Default to empty slice on error
		}

		byYearDay, err := s.jsonToInt64Slice(event.RecurrenceRule.ByYearDay)
		if err != nil {
			byYearDay = []int64{} // Default to empty slice on error
		}

		response.RecurrenceRule = &dto.RecurrenceRuleResponse{
			ID:         event.RecurrenceRule.ID,
			Frequency:  event.RecurrenceRule.Frequency,
			Interval:   event.RecurrenceRule.Interval,
			ByWeekday:  byWeekday,
			ByMonthDay: byMonthDay,
			ByMonth:    byMonth,
			BySetPos:   bySetPos,
			WeekStart:  event.RecurrenceRule.WeekStart,
			ByYearDay:  byYearDay,
			Count:      event.RecurrenceRule.Count,
			Until:      event.RecurrenceRule.Until,
		}
	}

	if len(event.Lagu) > 0 {
		laguResponses := make([]dto.LaguResponse, len(event.Lagu))
		for i, lagu := range event.Lagu {
			laguResponses[i] = dto.LaguResponse{
				ID:    lagu.ID,
				Title: lagu.Judul,
			}
		}
		response.Lagu = laguResponses
	}

	return response
}

func (s *eventService) generateOccurrences(event *entity.Event, startDate, endDate time.Time) ([]dto.EventOccurrenceResponse, error) {
	var occurrenceResponses []dto.EventOccurrenceResponse

	// Get exceptions for this event
	exceptions, err := s.eventRepo.GetRecurrenceExceptions(event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurrence exceptions: %w", err)
	}

	if event.RecurrenceRule == nil {
		// Single event - check if it's in range
		if (event.EventDate.Equal(startDate) || event.EventDate.After(startDate)) &&
			(event.EventDate.Equal(endDate) || event.EventDate.Before(endDate)) {
			occurrenceResponses = append(occurrenceResponses, dto.EventOccurrenceResponse{
				EventID:        event.ID,
				OccurrenceDate: event.EventDate,
				StartDatetime:  event.StartDatetime,
				EndDatetime:    event.EndDatetime,
				IsException:    false,
				IsSkipped:      false,
				OriginalEvent:  s.entityToResponse(event),
			})
		}
		return occurrenceResponses, nil
	}

	// Use the RecurrenceGenerator for recurring events
	occurrenceDates, err := s.recurrenceGenerator.GenerateOccurrences(event, event.RecurrenceRule, startDate, endDate, exceptions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate occurrences: %w", err)
	}

	// Create exception map for quick lookups
	exceptionMap := make(map[string]*entity.RecurrenceException)
	for i := range exceptions {
		dateKey := exceptions[i].ExceptionDate.Format("2006-01-02")
		exceptionMap[dateKey] = &exceptions[i]
	}

	// Convert generated dates to response objects
	for _, occurrenceDate := range occurrenceDates {
		dateKey := occurrenceDate.Format("2006-01-02")
		exception := exceptionMap[dateKey]

		// Calculate start and end times for this occurrence
		duration := event.EndDatetime.Sub(event.StartDatetime)
		startTime := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(),
			event.StartDatetime.Hour(), event.StartDatetime.Minute(), event.StartDatetime.Second(),
			event.StartDatetime.Nanosecond(), event.StartDatetime.Location())
		endTime := startTime.Add(duration)

		// Use exception override times if available
		if exception != nil {
			if exception.OverrideStart != nil {
				startTime = *exception.OverrideStart
			}
			if exception.OverrideEnd != nil {
				endTime = *exception.OverrideEnd
			}
		}

		occurrenceResponses = append(occurrenceResponses, dto.EventOccurrenceResponse{
			EventID:        event.ID,
			OccurrenceDate: occurrenceDate,
			StartDatetime:  startTime,
			EndDatetime:    endTime,
			IsException:    exception != nil,
			IsSkipped:      false,
			ExceptionNotes: func() string {
				if exception != nil {
					return exception.Notes
				}
				return ""
			}(),
			OriginalEvent: s.entityToResponse(event),
		})
	}

	return occurrenceResponses, nil
}

// New methods for handling single and future occurrence updates

func (s *eventService) UpdateSingleOccurrence(id uuid.UUID, req *dto.UpdateOccurrenceRequest) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event.RecurrenceRule == nil {
		return fmt.Errorf("event is not recurring")
	}

	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return fmt.Errorf("invalid occurrence date format: %w", err)
	}

	return s.createOrUpdateException(event, occurrenceDate, req.StartTime, req.EndTime, &req.Event, entity.ModificationTypeSingle)
}

func (s *eventService) UpdateFutureOccurrences(id uuid.UUID, req *dto.UpdateFutureOccurrencesRequest) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event.RecurrenceRule == nil {
		return fmt.Errorf("event is not recurring")
	}

	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		return fmt.Errorf("invalid from date format: %w", err)
	}

	// End the current series at the from date
	if err := s.eventRepo.SetRecurrenceUntilDate(event.ID, fromDate.AddDate(0, 0, -1)); err != nil {
		return fmt.Errorf("failed to end current series: %w", err)
	}

	// Create a new event for the updated series starting from the from date
	return s.createNewSeriesFromDate(event, fromDate, req.StartTime, req.EndTime, &req.Event, req.RecurrenceRule)
}

func (s *eventService) ValidateRecurrenceRule(rule *dto.CreateRecurrenceRuleRequest) error {
	// Convert DTO to entity for validation
	byWeekdayJSON, _ := s.sliceToJSON(rule.ByWeekday)
	byMonthDayJSON, _ := s.sliceToJSON(rule.ByMonthDay)
	byMonthJSON, _ := s.sliceToJSON(rule.ByMonth)
	bySetPosJSON, _ := s.sliceToJSON(rule.BySetPos)
	byYearDayJSON, _ := s.sliceToJSON(rule.ByYearDay)

	entityRule := &entity.RecurrenceRule{
		Frequency:  rule.Frequency,
		Interval:   rule.Interval,
		ByWeekday:  byWeekdayJSON,
		ByMonthDay: byMonthDayJSON,
		ByMonth:    byMonthJSON,
		BySetPos:   bySetPosJSON,
		WeekStart:  rule.WeekStart,
		ByYearDay:  byYearDayJSON,
	}

	if entityRule.Interval == 0 {
		entityRule.Interval = 1
	}
	if entityRule.WeekStart == "" {
		entityRule.WeekStart = "MO"
	}

	return s.recurrenceGenerator.ValidateRecurrenceRule(entityRule)
}

func (s *eventService) GetNextOccurrence(id uuid.UUID, after time.Time) (*time.Time, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return s.recurrenceGenerator.GetNextOccurrence(event, event.RecurrenceRule, after)
}

// Recurring event update helper methods

func (s *eventService) updateSingleOccurrenceFromRecurring(event *entity.Event, req *dto.UpdateRecurringEventRequest) error {
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return fmt.Errorf("invalid occurrence date format: %w", err)
	}

	return s.createOrUpdateException(event, occurrenceDate, req.StartTime, req.EndTime, &req.Event, entity.ModificationTypeSingle)
}

func (s *eventService) createOrUpdateException(event *entity.Event, occurrenceDate time.Time, startTime, endTime *string, eventUpdates *dto.UpdateEventRequest, modificationType string) error {
	// Create or get existing exception
	exception, err := s.eventRepo.GetExceptionByEventAndDate(event.ID, occurrenceDate)
	if err != nil {
		// Create new exception
		exception = &entity.RecurrenceException{
			EventID:          event.ID,
			ExceptionDate:    occurrenceDate,
			ModificationType: modificationType,
			IsSkipped:        false,
		}

		// Store original times for reference
		originalStart := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(),
			event.StartDatetime.Hour(), event.StartDatetime.Minute(), event.StartDatetime.Second(),
			event.StartDatetime.Nanosecond(), event.StartDatetime.Location())
		originalEnd := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(),
			event.EndDatetime.Hour(), event.EndDatetime.Minute(), event.EndDatetime.Second(),
			event.EndDatetime.Nanosecond(), event.EndDatetime.Location())

		exception.OriginalStartTime = &originalStart
		exception.OriginalEndTime = &originalEnd
	}

	// Apply time overrides if provided
	if startTime != nil && endTime != nil {
		parsedStartTime, err := time.Parse("15:04", *startTime)
		if err != nil {
			return fmt.Errorf("invalid start time format: %w", err)
		}
		parsedEndTime, err := time.Parse("15:04", *endTime)
		if err != nil {
			return fmt.Errorf("invalid end time format: %w", err)
		}

		// Validate that start time is before end time
		if !parsedStartTime.Before(parsedEndTime) {
			return fmt.Errorf("start time must be before end time")
		}

		startDateTime := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(),
			parsedStartTime.Hour(), parsedStartTime.Minute(), 0, 0, time.UTC)
		endDateTime := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(),
			parsedEndTime.Hour(), parsedEndTime.Minute(), 0, 0, time.UTC)

		exception.OverrideStart = &startDateTime
		exception.OverrideEnd = &endDateTime
	}

	// Add notes if there are other changes
	if eventUpdates.Title != nil || eventUpdates.Description != nil || eventUpdates.EventLocation != nil {
		exception.Notes = fmt.Sprintf("Modified on %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	// Save the exception
	if exception.ID == uuid.Nil {
		return s.eventRepo.CreateRecurrenceException(exception)
	}
	return s.eventRepo.UpdateRecurrenceException(exception)
}

func (s *eventService) createNewSeriesFromDate(originalEvent *entity.Event, fromDate time.Time, startTime, endTime *string, eventUpdates *dto.UpdateEventRequest, recurrenceRule *dto.CreateRecurrenceRuleRequest) error {
	// Build create request for new series
	createReq := &dto.CreateEventRequest{
		Title:                 originalEvent.Title,
		BannerImage:           originalEvent.BannerImage,
		Description:           originalEvent.Description,
		Capacity:              originalEvent.Capacity,
		Type:                  originalEvent.Type,
		EventDate:             fromDate.Format("2006-01-02"),
		EventLocation:         originalEvent.EventLocation,
		StartTime:             originalEvent.StartDatetime.Format("15:04"),
		EndTime:               originalEvent.EndDatetime.Format("15:04"),
		AllDay:                originalEvent.AllDay,
		Timezone:              originalEvent.Timezone,
		IsPublic:              originalEvent.IsPublic,
		DiscipleshipJourneyID: originalEvent.DiscipleshipJourneyID,
	}

	// Apply updates from request
	if eventUpdates.Title != nil {
		createReq.Title = *eventUpdates.Title
	}
	if eventUpdates.BannerImage != nil {
		createReq.BannerImage = *eventUpdates.BannerImage
	}
	if eventUpdates.Description != nil {
		createReq.Description = *eventUpdates.Description
	}
	if eventUpdates.EventLocation != nil {
		createReq.EventLocation = *eventUpdates.EventLocation
	}
	if startTime != nil {
		createReq.StartTime = *startTime
	}
	if endTime != nil {
		createReq.EndTime = *endTime
	}

	// Set up recurrence rule - use provided rule or copy from original
	if recurrenceRule != nil {
		createReq.RecurrenceRule = recurrenceRule
	} else if originalEvent.RecurrenceRule != nil {
		// Convert JSON strings back to slices for DTO
		byWeekday, _ := s.jsonToStringSlice(originalEvent.RecurrenceRule.ByWeekday)
		byMonthDay, _ := s.jsonToInt64Slice(originalEvent.RecurrenceRule.ByMonthDay)
		byMonth, _ := s.jsonToInt64Slice(originalEvent.RecurrenceRule.ByMonth)
		bySetPos, _ := s.jsonToInt64Slice(originalEvent.RecurrenceRule.BySetPos)
		byYearDay, _ := s.jsonToInt64Slice(originalEvent.RecurrenceRule.ByYearDay)

		createReq.RecurrenceRule = &dto.CreateRecurrenceRuleRequest{
			Frequency:  originalEvent.RecurrenceRule.Frequency,
			Interval:   originalEvent.RecurrenceRule.Interval,
			ByWeekday:  byWeekday,
			ByMonthDay: byMonthDay,
			ByMonth:    byMonth,
			BySetPos:   bySetPos,
			WeekStart:  originalEvent.RecurrenceRule.WeekStart,
			ByYearDay:  byYearDay,
			Count:      originalEvent.RecurrenceRule.Count,
		}
		if originalEvent.RecurrenceRule.Until != nil {
			until := originalEvent.RecurrenceRule.Until.Format("2006-01-02")
			createReq.RecurrenceRule.Until = &until
		}
	}

	// Create the new event series
	_, err := s.CreateEvent(createReq)
	return err
}

func (s *eventService) updateEntireSeries(event *entity.Event, req *dto.UpdateRecurringEventRequest) error {
	// Update the main event record
	_, err := s.UpdateEvent(event.ID, &req.Event)
	return err
}

func (s *eventService) updateThisAndFutureOccurrences(event *entity.Event, req *dto.UpdateRecurringEventRequest) error {
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// End the current series at the occurrence date (before it)
	if err := s.eventRepo.SetRecurrenceUntilDate(event.ID, occurrenceDate.AddDate(0, 0, -1)); err != nil {
		return fmt.Errorf("failed to end current series: %w", err)
	}

	// Create a new event for the updated series starting from the occurrence date
	return s.createNewSeriesFromDate(event, occurrenceDate, req.StartTime, req.EndTime, &req.Event, nil)
}

func (s *eventService) skipSingleOccurrence(event *entity.Event, occurrenceDate time.Time) error {
	exception := &entity.RecurrenceException{
		EventID:          event.ID,
		ExceptionDate:    occurrenceDate,
		ModificationType: entity.ModificationTypeSingle,
		IsSkipped:        true,
		Notes:            "Occurrence deleted by user",
	}

	return s.eventRepo.CreateRecurrenceException(exception)
}

func (s *eventService) deleteFutureOccurrences(event *entity.Event, fromDate time.Time) error {
	return s.eventRepo.DeleteFutureOccurrences(event.ID, fromDate)
}

// Helper function to convert slice to JSON string
func (s *eventService) sliceToJSON(slice interface{}) (string, error) {
	if slice == nil {
		return "", nil
	}

	// Handle empty slices
	switch v := slice.(type) {
	case []string:
		if len(v) == 0 {
			return "", nil
		}
	case []int64:
		if len(v) == 0 {
			return "", nil
		}
	}

	jsonBytes, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// Helper function to convert JSON string to string slice
func (s *eventService) jsonToStringSlice(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}

	var result []string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// Helper function to convert JSON string to int64 slice
func (s *eventService) jsonToInt64Slice(jsonStr string) ([]int64, error) {
	if jsonStr == "" {
		return []int64{}, nil
	}

	var result []int64
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return []int64{}, err
	}

	return result, nil
}
