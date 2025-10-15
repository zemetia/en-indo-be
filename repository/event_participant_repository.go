package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type EventParticipantRepository interface {
	Create(participant *entity.EventParticipant) error
	BulkCreate(participants []entity.EventParticipant) error
	Update(participant *entity.EventParticipant) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entity.EventParticipant, error)
	GetByEventAndParticipant(eventID uuid.UUID, occurrenceDate time.Time, participantType entity.ParticipantType, participantID uuid.UUID) (*entity.EventParticipant, error)
	GetByEventAndOccurrence(eventID uuid.UUID, occurrenceDate time.Time, filters EventParticipantFilters) ([]entity.EventParticipant, int64, error)
	List(filters EventParticipantFilters) ([]entity.EventParticipant, int64, error)
	GetAttendanceStats(eventID uuid.UUID, occurrenceDate time.Time) (*AttendanceStats, error)
	LoadParticipantsData(participants []entity.EventParticipant) (map[uuid.UUID]*entity.Person, map[uuid.UUID]*entity.Visitor, error)
}

type EventParticipantFilters struct {
	EventID            *uuid.UUID
	OccurrenceDate     *time.Time
	ParticipantType    *entity.ParticipantType
	RegistrationStatus *entity.RegistrationStatus
	AttendanceStatus   *entity.AttendanceStatus
	Limit              int
	Offset             int
}

type AttendanceStats struct {
	TotalRegistered int
	TotalPresent    int
	TotalAbsent     int
	TotalExcused    int
	TotalCancelled  int
	TotalNoShow     int
}

type eventParticipantRepository struct {
	db *gorm.DB
}

func NewEventParticipantRepository(db *gorm.DB) EventParticipantRepository {
	return &eventParticipantRepository{db: db}
}

func (r *eventParticipantRepository) Create(participant *entity.EventParticipant) error {
	return r.db.Create(participant).Error
}

func (r *eventParticipantRepository) BulkCreate(participants []entity.EventParticipant) error {
	if len(participants) == 0 {
		return nil
	}
	return r.db.Create(&participants).Error
}

func (r *eventParticipantRepository) Update(participant *entity.EventParticipant) error {
	return r.db.Save(participant).Error
}

func (r *eventParticipantRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.EventParticipant{}, id).Error
}

func (r *eventParticipantRepository) GetByID(id uuid.UUID) (*entity.EventParticipant, error) {
	var participant entity.EventParticipant
	err := r.db.
		Preload("Event").
		Preload("RecordedByPerson").
		First(&participant, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *eventParticipantRepository) GetByEventAndParticipant(
	eventID uuid.UUID,
	occurrenceDate time.Time,
	participantType entity.ParticipantType,
	participantID uuid.UUID,
) (*entity.EventParticipant, error) {
	var participant entity.EventParticipant

	// Use DATE() function to compare only the date portion, ignoring time/timezone
	err := r.db.
		Preload("Event").
		Preload("RecordedByPerson").
		Where("event_id = ? AND DATE(occurrence_date) = DATE(?) AND participant_type = ? AND participant_id = ?",
			eventID, occurrenceDate, participantType, participantID).
		First(&participant).Error

	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *eventParticipantRepository) GetByEventAndOccurrence(
	eventID uuid.UUID,
	occurrenceDate time.Time,
	filters EventParticipantFilters,
) ([]entity.EventParticipant, int64, error) {
	var participants []entity.EventParticipant
	var total int64

	// Use DATE() function to compare only the date portion, ignoring time/timezone
	query := r.db.Model(&entity.EventParticipant{}).
		Where("event_id = ? AND DATE(occurrence_date) = DATE(?)", eventID, occurrenceDate)

	// Apply filters
	query = r.applyFilters(query, filters)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit).Offset(filters.Offset)
	}

	// Fetch with preloads
	err := query.
		Preload("Event").
		Preload("RecordedByPerson").
		Order("created_at DESC").
		Find(&participants).Error

	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

func (r *eventParticipantRepository) List(filters EventParticipantFilters) ([]entity.EventParticipant, int64, error) {
	var participants []entity.EventParticipant
	var total int64

	query := r.db.Model(&entity.EventParticipant{})

	// Apply filters
	query = r.applyFilters(query, filters)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit).Offset(filters.Offset)
	}

	// Fetch with preloads
	err := query.
		Preload("Event").
		Preload("RecordedByPerson").
		Order("occurrence_date DESC, created_at DESC").
		Find(&participants).Error

	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

func (r *eventParticipantRepository) GetAttendanceStats(eventID uuid.UUID, occurrenceDate time.Time) (*AttendanceStats, error) {
	var stats AttendanceStats

	// Count by attendance status - using DATE() for proper comparison
	var attendanceResults []struct {
		AttendanceStatus string
		Count            int
	}
	if err := r.db.Model(&entity.EventParticipant{}).
		Select("attendance_status, COUNT(*) as count").
		Where("event_id = ? AND DATE(occurrence_date) = DATE(?)", eventID, occurrenceDate).
		Group("attendance_status").
		Scan(&attendanceResults).Error; err != nil {
		return nil, err
	}

	for _, result := range attendanceResults {
		switch entity.AttendanceStatus(result.AttendanceStatus) {
		case entity.AttendanceStatusPresent:
			stats.TotalPresent = result.Count
		case entity.AttendanceStatusAbsent:
			stats.TotalAbsent = result.Count
		case entity.AttendanceStatusExcused:
			stats.TotalExcused = result.Count
		}
	}

	// Count by registration status - using DATE() for proper comparison
	var registrationResults []struct {
		RegistrationStatus string
		Count              int
	}
	if err := r.db.Model(&entity.EventParticipant{}).
		Select("registration_status, COUNT(*) as count").
		Where("event_id = ? AND DATE(occurrence_date) = DATE(?)", eventID, occurrenceDate).
		Group("registration_status").
		Scan(&registrationResults).Error; err != nil {
		return nil, err
	}

	for _, result := range registrationResults {
		switch entity.RegistrationStatus(result.RegistrationStatus) {
		case entity.RegistrationStatusCancelled:
			stats.TotalCancelled = result.Count
		case entity.RegistrationStatusNoShow:
			stats.TotalNoShow = result.Count
		}
	}

	// Calculate total registered (excluding cancelled) - using DATE() for proper comparison
	var totalRegistered int64
	if err := r.db.Model(&entity.EventParticipant{}).
		Where("event_id = ? AND DATE(occurrence_date) = DATE(?) AND registration_status != ?",
			eventID, occurrenceDate, entity.RegistrationStatusCancelled).
		Count(&totalRegistered).Error; err != nil {
		return nil, err
	}
	stats.TotalRegistered = int(totalRegistered)

	return &stats, nil
}

// LoadParticipantsData batch loads all people and visitors for a list of participants
// Returns two maps: personID -> Person and visitorID -> Visitor
// This is much more efficient than loading participants one by one
func (r *eventParticipantRepository) LoadParticipantsData(participants []entity.EventParticipant) (map[uuid.UUID]*entity.Person, map[uuid.UUID]*entity.Visitor, error) {
	// Collect all person and visitor IDs
	personIDs := make([]uuid.UUID, 0)
	visitorIDs := make([]uuid.UUID, 0)

	for _, p := range participants {
		if p.ParticipantType == entity.ParticipantTypePerson {
			personIDs = append(personIDs, p.ParticipantID)
		} else if p.ParticipantType == entity.ParticipantTypeVisitor {
			visitorIDs = append(visitorIDs, p.ParticipantID)
		}
	}

	// Batch load all people
	peopleMap := make(map[uuid.UUID]*entity.Person)
	if len(personIDs) > 0 {
		var people []entity.Person
		if err := r.db.Where("id IN ?", personIDs).Find(&people).Error; err != nil {
			return nil, nil, err
		}
		for i := range people {
			peopleMap[people[i].ID] = &people[i]
		}
	}

	// Batch load all visitors
	visitorsMap := make(map[uuid.UUID]*entity.Visitor)
	if len(visitorIDs) > 0 {
		var visitors []entity.Visitor
		if err := r.db.Where("id IN ?", visitorIDs).Find(&visitors).Error; err != nil {
			return nil, nil, err
		}
		for i := range visitors {
			visitorsMap[visitors[i].ID] = &visitors[i]
		}
	}

	return peopleMap, visitorsMap, nil
}

// Helper method to apply filters to query
func (r *eventParticipantRepository) applyFilters(query *gorm.DB, filters EventParticipantFilters) *gorm.DB {
	if filters.EventID != nil {
		query = query.Where("event_id = ?", *filters.EventID)
	}
	if filters.OccurrenceDate != nil {
		// Use DATE() function to compare only the date portion
		query = query.Where("DATE(occurrence_date) = DATE(?)", *filters.OccurrenceDate)
	}
	if filters.ParticipantType != nil {
		query = query.Where("participant_type = ?", *filters.ParticipantType)
	}
	if filters.RegistrationStatus != nil {
		query = query.Where("registration_status = ?", *filters.RegistrationStatus)
	}
	if filters.AttendanceStatus != nil {
		query = query.Where("attendance_status = ?", *filters.AttendanceStatus)
	}
	return query
}
