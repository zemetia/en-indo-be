package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type EventPICService interface {
	// Basic PIC operations
	CreateEventPIC(eventID uuid.UUID, req *dto.CreateEventPICRequest, createdBy uuid.UUID) (*dto.EventPICResponse, error)
	GetEventPIC(id uuid.UUID) (*dto.EventPICResponse, error)
	UpdateEventPIC(id uuid.UUID, req *dto.UpdateEventPICRequest, updatedBy uuid.UUID) (*dto.EventPICResponse, error)
	DeleteEventPIC(id uuid.UUID, deletedBy uuid.UUID, reason string) error
	
	// Event-specific PIC operations
	GetPICsByEventID(eventID uuid.UUID) ([]dto.EventPICResponse, error)
	GetPrimaryPICByEventID(eventID uuid.UUID) (*dto.EventPICResponse, error)
	GetActivePICsByEventID(eventID uuid.UUID) ([]dto.EventPICResponse, error)
	
	// Person-specific PIC operations
	GetPICsByPersonID(personID uuid.UUID) ([]dto.EventPICResponse, error)
	GetActivePICsByPersonID(personID uuid.UUID) ([]dto.EventPICResponse, error)
	
	// Advanced operations
	ListEventPICs(req *dto.EventPICFilterRequest) (*dto.EventPICListResponse, error)
	AssignMultiplePICs(eventID uuid.UUID, req *dto.BulkAssignEventPICRequest, createdBy uuid.UUID) error
	TransferPICRole(eventID uuid.UUID, req *dto.TransferEventPICRequest, changedBy uuid.UUID) error
	
	// Validation operations
	ValidateEventPICPermissions(eventID, personID uuid.UUID, action string) (bool, error)
	CheckPICConflicts(eventID, personID uuid.UUID, role string, isPrimary bool) error
	
	// Role management
	CreateEventPICRole(req *dto.CreateEventPICRoleRequest) (*dto.EventPICRoleResponse, error)
	GetEventPICRole(id uuid.UUID) (*dto.EventPICRoleResponse, error)
	UpdateEventPICRole(id uuid.UUID, req *dto.UpdateEventPICRoleRequest) (*dto.EventPICRoleResponse, error)
	DeleteEventPICRole(id uuid.UUID) error
	ListEventPICRoles(page, limit int, search string) (*dto.EventPICRoleListResponse, error)
	
	// History and audit
	GetEventPICHistory(eventID uuid.UUID) ([]dto.EventPICHistoryResponse, error)
	GetPersonPICHistory(personID uuid.UUID) ([]dto.EventPICHistoryResponse, error)
	
	// Utility methods
	GetExpiringPICs(days int) ([]dto.EventPICResponse, error)
	NotifyPICsForEvent(eventID uuid.UUID, message string) error
}

type eventPICService struct {
	eventPICRepo repository.EventPICRepository
	eventRepo    repository.EventRepository
}

func NewEventPICService(eventPICRepo repository.EventPICRepository, eventRepo repository.EventRepository) EventPICService {
	return &eventPICService{
		eventPICRepo: eventPICRepo,
		eventRepo:    eventRepo,
	}
}

// Basic PIC operations
func (s *eventPICService) CreateEventPIC(eventID uuid.UUID, req *dto.CreateEventPICRequest, createdBy uuid.UUID) (*dto.EventPICResponse, error) {
	// Validate event exists
	_, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}
	
	// Check for conflicts
	if err := s.CheckPICConflicts(eventID, req.PersonID, req.Role, req.IsPrimary); err != nil {
		return nil, err
	}
	
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	
	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		endDate = &parsed
	}
	
	// If this is a primary PIC, ensure no other primary exists
	if req.IsPrimary {
		primaryCount, err := s.eventPICRepo.CountPrimaryPICsForEvent(eventID)
		if err != nil {
			return nil, fmt.Errorf("failed to check primary PIC count: %w", err)
		}
		if primaryCount > 0 {
			return nil, fmt.Errorf("event already has a primary PIC")
		}
	}
	
	// Create EventPIC entity
	eventPIC := &entity.EventPIC{
		EventID:           eventID,
		PersonID:          req.PersonID,
		Role:              req.Role,
		Description:       req.Description,
		IsActive:          true,
		IsPrimary:         req.IsPrimary,
		StartDate:         startDate,
		EndDate:           endDate,
		CanEdit:           req.CanEdit,
		CanDelete:         req.CanDelete,
		CanAssignPIC:      req.CanAssignPIC,
		NotifyOnChanges:   req.NotifyOnChanges,
		NotifyOnReminders: req.NotifyOnReminders,
	}
	
	if err := s.eventPICRepo.Create(eventPIC); err != nil {
		return nil, fmt.Errorf("failed to create event PIC: %w", err)
	}
	
	// Create history record
	history := &entity.EventPICHistory{
		EventID:    eventID,
		PersonID:   req.PersonID,
		Action:     entity.EventPICActionAssigned,
		NewRole:    req.Role,
		ChangedBy:  createdBy,
		Reason:     "PIC assigned",
		ActionDate: time.Now(),
	}
	s.eventPICRepo.CreateHistory(history)
	
	return s.GetEventPIC(eventPIC.ID)
}

func (s *eventPICService) GetEventPIC(id uuid.UUID) (*dto.EventPICResponse, error) {
	eventPIC, err := s.eventPICRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event PIC: %w", err)
	}
	
	return s.entityToResponse(eventPIC), nil
}

func (s *eventPICService) UpdateEventPIC(id uuid.UUID, req *dto.UpdateEventPICRequest, updatedBy uuid.UUID) (*dto.EventPICResponse, error) {
	eventPIC, err := s.eventPICRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event PIC: %w", err)
	}
	
	oldRole := eventPIC.Role
	historyAction := entity.EventPICActionPermissionsChanged
	
	// Update fields if provided
	if req.Role != nil {
		eventPIC.Role = *req.Role
		historyAction = entity.EventPICActionRoleChanged
	}
	if req.Description != nil {
		eventPIC.Description = *req.Description
	}
	if req.IsActive != nil {
		eventPIC.IsActive = *req.IsActive
	}
	if req.IsPrimary != nil {
		// Check for primary PIC conflicts
		if *req.IsPrimary && !eventPIC.IsPrimary {
			primaryCount, err := s.eventPICRepo.CountPrimaryPICsForEvent(eventPIC.EventID)
			if err != nil {
				return nil, fmt.Errorf("failed to check primary PIC count: %w", err)
			}
			if primaryCount > 0 {
				return nil, fmt.Errorf("event already has a primary PIC")
			}
		}
		eventPIC.IsPrimary = *req.IsPrimary
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		eventPIC.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		eventPIC.EndDate = &endDate
	}
	if req.CanEdit != nil {
		eventPIC.CanEdit = *req.CanEdit
	}
	if req.CanDelete != nil {
		eventPIC.CanDelete = *req.CanDelete
	}
	if req.CanAssignPIC != nil {
		eventPIC.CanAssignPIC = *req.CanAssignPIC
	}
	if req.NotifyOnChanges != nil {
		eventPIC.NotifyOnChanges = *req.NotifyOnChanges
	}
	if req.NotifyOnReminders != nil {
		eventPIC.NotifyOnReminders = *req.NotifyOnReminders
	}
	
	if err := s.eventPICRepo.Update(eventPIC); err != nil {
		return nil, fmt.Errorf("failed to update event PIC: %w", err)
	}
	
	// Create history record
	history := &entity.EventPICHistory{
		EventID:    eventPIC.EventID,
		PersonID:   eventPIC.PersonID,
		Action:     historyAction,
		OldRole:    oldRole,
		NewRole:    eventPIC.Role,
		ChangedBy:  updatedBy,
		Reason:     "PIC updated",
		ActionDate: time.Now(),
	}
	s.eventPICRepo.CreateHistory(history)
	
	return s.GetEventPIC(eventPIC.ID)
}

func (s *eventPICService) DeleteEventPIC(id uuid.UUID, deletedBy uuid.UUID, reason string) error {
	eventPIC, err := s.eventPICRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event PIC: %w", err)
	}
	
	// Create history record before deletion
	history := &entity.EventPICHistory{
		EventID:    eventPIC.EventID,
		PersonID:   eventPIC.PersonID,
		Action:     entity.EventPICActionRemoved,
		OldRole:    eventPIC.Role,
		ChangedBy:  deletedBy,
		Reason:     reason,
		ActionDate: time.Now(),
	}
	s.eventPICRepo.CreateHistory(history)
	
	if err := s.eventPICRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete event PIC: %w", err)
	}
	
	return nil
}

// Event-specific PIC operations
func (s *eventPICService) GetPICsByEventID(eventID uuid.UUID) ([]dto.EventPICResponse, error) {
	pics, err := s.eventPICRepo.GetPICsByEventID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PICs for event: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return responses, nil
}

func (s *eventPICService) GetPrimaryPICByEventID(eventID uuid.UUID) (*dto.EventPICResponse, error) {
	pic, err := s.eventPICRepo.GetPrimaryPICByEventID(eventID)
	if err != nil {
		return nil, err
	}
	
	return s.entityToResponse(pic), nil
}

func (s *eventPICService) GetActivePICsByEventID(eventID uuid.UUID) ([]dto.EventPICResponse, error) {
	pics, err := s.eventPICRepo.GetActivePICsByEventID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active PICs for event: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return responses, nil
}

// Person-specific PIC operations
func (s *eventPICService) GetPICsByPersonID(personID uuid.UUID) ([]dto.EventPICResponse, error) {
	pics, err := s.eventPICRepo.GetPICsByPersonID(personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PICs for person: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return responses, nil
}

func (s *eventPICService) GetActivePICsByPersonID(personID uuid.UUID) ([]dto.EventPICResponse, error) {
	pics, err := s.eventPICRepo.GetActivePICsByPersonID(personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active PICs for person: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return responses, nil
}

// Advanced operations
func (s *eventPICService) ListEventPICs(req *dto.EventPICFilterRequest) (*dto.EventPICListResponse, error) {
	filters := repository.EventPICFilters{
		EventID:   req.EventID,
		PersonID:  req.PersonID,
		Role:      req.Role,
		IsActive:  req.IsActive,
		IsPrimary: req.IsPrimary,
		Search:    req.Search,
		Limit:     req.Limit,
		Offset:    (req.Page - 1) * req.Limit,
	}
	
	// Set defaults
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}
	
	pics, total, err := s.eventPICRepo.List(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list event PICs: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return &dto.EventPICListResponse{
		PICs:       responses,
		TotalCount: int(total),
		Page:       req.Page,
		Limit:      filters.Limit,
	}, nil
}

func (s *eventPICService) AssignMultiplePICs(eventID uuid.UUID, req *dto.BulkAssignEventPICRequest, createdBy uuid.UUID) error {
	// Validate event exists
	_, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	
	var eventPICs []entity.EventPIC
	
	for _, picReq := range req.PICs {
		// Validate each PIC assignment
		if err := s.CheckPICConflicts(eventID, picReq.PersonID, picReq.Role, picReq.IsPrimary); err != nil {
			return fmt.Errorf("PIC conflict for person %s: %w", picReq.PersonID, err)
		}
		
		// Parse dates
		startDate, err := time.Parse("2006-01-02", picReq.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format for person %s: %w", picReq.PersonID, err)
		}
		
		var endDate *time.Time
		if picReq.EndDate != nil {
			parsed, err := time.Parse("2006-01-02", *picReq.EndDate)
			if err != nil {
				return fmt.Errorf("invalid end date format for person %s: %w", picReq.PersonID, err)
			}
			endDate = &parsed
		}
		
		eventPIC := entity.EventPIC{
			EventID:           eventID,
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
	
	// Check for primary PIC conflicts across the batch
	primaryCount := 0
	for _, pic := range eventPICs {
		if pic.IsPrimary {
			primaryCount++
		}
	}
	if primaryCount > 1 {
		return fmt.Errorf("cannot assign multiple primary PICs")
	}
	
	// Create all PICs
	if err := s.eventPICRepo.CreateMultiple(eventPICs); err != nil {
		return fmt.Errorf("failed to create multiple event PICs: %w", err)
	}
	
	// Create history records
	for _, pic := range eventPICs {
		history := &entity.EventPICHistory{
			EventID:    eventID,
			PersonID:   pic.PersonID,
			Action:     entity.EventPICActionAssigned,
			NewRole:    pic.Role,
			ChangedBy:  createdBy,
			Reason:     "Bulk PIC assignment",
			ActionDate: time.Now(),
		}
		s.eventPICRepo.CreateHistory(history)
	}
	
	return nil
}

func (s *eventPICService) TransferPICRole(eventID uuid.UUID, req *dto.TransferEventPICRequest, changedBy uuid.UUID) error {
	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return fmt.Errorf("invalid effective date format: %w", err)
	}
	
	if err := s.eventPICRepo.TransferPICRole(req.FromPersonID, req.ToPersonID, eventID, req.TransferType); err != nil {
		return fmt.Errorf("failed to transfer PIC role: %w", err)
	}
	
	// Create history records
	history := &entity.EventPICHistory{
		EventID:    eventID,
		PersonID:   req.FromPersonID,
		Action:     "transferred",
		ChangedBy:  changedBy,
		Reason:     req.Reason,
		ActionDate: effectiveDate,
	}
	s.eventPICRepo.CreateHistory(history)
	
	return nil
}

// Validation operations
func (s *eventPICService) ValidateEventPICPermissions(eventID, personID uuid.UUID, action string) (bool, error) {
	pics, err := s.eventPICRepo.GetPICsWithPermission(eventID, action)
	if err != nil {
		return false, err
	}
	
	for _, pic := range pics {
		if pic.PersonID == personID {
			return true, nil
		}
	}
	
	return false, nil
}

func (s *eventPICService) CheckPICConflicts(eventID, personID uuid.UUID, role string, isPrimary bool) error {
	// Check if person already has PIC role for this event
	hasPIC, err := s.eventPICRepo.HasPersonPICRoleForEvent(eventID, personID)
	if err != nil {
		return err
	}
	if hasPIC {
		return fmt.Errorf("person already assigned as PIC for this event")
	}
	
	// Check primary PIC conflicts
	if isPrimary {
		primaryCount, err := s.eventPICRepo.CountPrimaryPICsForEvent(eventID)
		if err != nil {
			return err
		}
		if primaryCount > 0 {
			return fmt.Errorf("event already has a primary PIC")
		}
	}
	
	return nil
}

// Role management
func (s *eventPICService) CreateEventPICRole(req *dto.CreateEventPICRoleRequest) (*dto.EventPICRoleResponse, error) {
	// Check if role name already exists
	_, err := s.eventPICRepo.GetRoleByName(req.Name)
	if err == nil {
		return nil, fmt.Errorf("role name already exists")
	}
	
	role := &entity.EventPICRole{
		Name:                req.Name,
		Description:         req.Description,
		IsActive:            true,
		DefaultCanEdit:      req.DefaultCanEdit,
		DefaultCanDelete:    req.DefaultCanDelete,
		DefaultCanAssignPIC: req.DefaultCanAssignPIC,
	}
	
	if err := s.eventPICRepo.CreateRole(role); err != nil {
		return nil, fmt.Errorf("failed to create event PIC role: %w", err)
	}
	
	return s.roleEntityToResponse(role), nil
}

func (s *eventPICService) GetEventPICRole(id uuid.UUID) (*dto.EventPICRoleResponse, error) {
	role, err := s.eventPICRepo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event PIC role: %w", err)
	}
	
	return s.roleEntityToResponse(role), nil
}

func (s *eventPICService) UpdateEventPICRole(id uuid.UUID, req *dto.UpdateEventPICRoleRequest) (*dto.EventPICRoleResponse, error) {
	role, err := s.eventPICRepo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event PIC role: %w", err)
	}
	
	// Update fields if provided
	if req.Name != nil {
		// Check for name conflicts
		existingRole, err := s.eventPICRepo.GetRoleByName(*req.Name)
		if err == nil && existingRole.ID != id {
			return nil, fmt.Errorf("role name already exists")
		}
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}
	if req.DefaultCanEdit != nil {
		role.DefaultCanEdit = *req.DefaultCanEdit
	}
	if req.DefaultCanDelete != nil {
		role.DefaultCanDelete = *req.DefaultCanDelete
	}
	if req.DefaultCanAssignPIC != nil {
		role.DefaultCanAssignPIC = *req.DefaultCanAssignPIC
	}
	
	if err := s.eventPICRepo.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("failed to update event PIC role: %w", err)
	}
	
	return s.roleEntityToResponse(role), nil
}

func (s *eventPICService) DeleteEventPICRole(id uuid.UUID) error {
	if err := s.eventPICRepo.DeleteRole(id); err != nil {
		return fmt.Errorf("failed to delete event PIC role: %w", err)
	}
	return nil
}

func (s *eventPICService) ListEventPICRoles(page, limit int, search string) (*dto.EventPICRoleListResponse, error) {
	if limit == 0 {
		limit = 20
	}
	if page == 0 {
		page = 1
	}
	
	filters := repository.EventPICRoleFilters{
		Search: search,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
	
	roles, total, err := s.eventPICRepo.ListRoles(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list event PIC roles: %w", err)
	}
	
	responses := make([]dto.EventPICRoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *s.roleEntityToResponse(&role)
	}
	
	return &dto.EventPICRoleListResponse{
		Roles:      responses,
		TotalCount: int(total),
		Page:       page,
		Limit:      limit,
	}, nil
}

// History and audit
func (s *eventPICService) GetEventPICHistory(eventID uuid.UUID) ([]dto.EventPICHistoryResponse, error) {
	history, err := s.eventPICRepo.GetHistoryByEventID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event PIC history: %w", err)
	}
	
	responses := make([]dto.EventPICHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = s.historyEntityToResponse(&h)
	}
	
	return responses, nil
}

func (s *eventPICService) GetPersonPICHistory(personID uuid.UUID) ([]dto.EventPICHistoryResponse, error) {
	history, err := s.eventPICRepo.GetHistoryByPersonID(personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person PIC history: %w", err)
	}
	
	responses := make([]dto.EventPICHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = s.historyEntityToResponse(&h)
	}
	
	return responses, nil
}

// Utility methods
func (s *eventPICService) GetExpiringPICs(days int) ([]dto.EventPICResponse, error) {
	pics, err := s.eventPICRepo.GetExpiringPICs(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring PICs: %w", err)
	}
	
	responses := make([]dto.EventPICResponse, len(pics))
	for i, pic := range pics {
		responses[i] = *s.entityToResponse(&pic)
	}
	
	return responses, nil
}

func (s *eventPICService) NotifyPICsForEvent(eventID uuid.UUID, message string) error {
	pics, err := s.eventPICRepo.GetActivePICsByEventID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get active PICs: %w", err)
	}
	
	// TODO: Implement notification logic (email, push notification, etc.)
	for _, pic := range pics {
		if pic.NotifyOnChanges {
			// Send notification to pic.Person
			fmt.Printf("Notification to %s: %s\n", pic.Person.Email, message)
		}
	}
	
	return nil
}

// Helper methods
func (s *eventPICService) entityToResponse(eventPIC *entity.EventPIC) *dto.EventPICResponse {
	return &dto.EventPICResponse{
		ID:          eventPIC.ID,
		EventID:     eventPIC.EventID,
		PersonID:    eventPIC.PersonID,
		Person: dto.PersonSummary{
			ID:           eventPIC.Person.ID,
			Nama:         eventPIC.Person.Nama,
			Email:        eventPIC.Person.Email,
			NomorTelepon: eventPIC.Person.NomorTelepon,
			ChurchID:     eventPIC.Person.ChurchID,
		},
		Role:              eventPIC.Role,
		Description:       eventPIC.Description,
		IsActive:          eventPIC.IsActive,
		IsPrimary:         eventPIC.IsPrimary,
		StartDate:         eventPIC.StartDate,
		EndDate:           eventPIC.EndDate,
		CanEdit:           eventPIC.CanEdit,
		CanDelete:         eventPIC.CanDelete,
		CanAssignPIC:      eventPIC.CanAssignPIC,
		NotifyOnChanges:   eventPIC.NotifyOnChanges,
		NotifyOnReminders: eventPIC.NotifyOnReminders,
		CreatedAt:         eventPIC.CreatedAt,
		UpdatedAt:         eventPIC.UpdatedAt,
	}
}

func (s *eventPICService) roleEntityToResponse(role *entity.EventPICRole) *dto.EventPICRoleResponse {
	return &dto.EventPICRoleResponse{
		ID:                  role.ID,
		Name:                role.Name,
		Description:         role.Description,
		IsActive:            role.IsActive,
		DefaultCanEdit:      role.DefaultCanEdit,
		DefaultCanDelete:    role.DefaultCanDelete,
		DefaultCanAssignPIC: role.DefaultCanAssignPIC,
		CreatedAt:           role.CreatedAt,
		UpdatedAt:           role.UpdatedAt,
	}
}

func (s *eventPICService) historyEntityToResponse(history *entity.EventPICHistory) dto.EventPICHistoryResponse {
	return dto.EventPICHistoryResponse{
		ID:       history.ID,
		EventID:  history.EventID,
		PersonID: history.PersonID,
		Person: dto.PersonSummary{
			ID:           history.Person.ID,
			Nama:         history.Person.Nama,
			Email:        history.Person.Email,
			NomorTelepon: history.Person.NomorTelepon,
			ChurchID:     history.Person.ChurchID,
		},
		Action:     history.Action,
		OldRole:    history.OldRole,
		NewRole:    history.NewRole,
		ChangedBy:  history.ChangedBy,
		ChangedByPerson: dto.PersonSummary{
			ID:           history.ChangedByPerson.ID,
			Nama:         history.ChangedByPerson.Nama,
			Email:        history.ChangedByPerson.Email,
			NomorTelepon: history.ChangedByPerson.NomorTelepon,
			ChurchID:     history.ChangedByPerson.ChurchID,
		},
		Reason:     history.Reason,
		ActionDate: history.ActionDate,
		CreatedAt:  history.CreatedAt,
	}
}