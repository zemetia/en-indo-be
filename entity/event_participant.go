package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ParticipantType represents the type of participant
type ParticipantType string

const (
	ParticipantTypePerson  ParticipantType = "person"
	ParticipantTypeVisitor ParticipantType = "visitor"
)

// RegistrationStatus represents registration state
type RegistrationStatus string

const (
	RegistrationStatusRegistered RegistrationStatus = "registered"
	RegistrationStatusAttended   RegistrationStatus = "attended"
	RegistrationStatusCancelled  RegistrationStatus = "cancelled"
	RegistrationStatusNoShow     RegistrationStatus = "no_show"
)

// AttendanceStatus represents actual attendance
type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
	AttendanceStatusExcused AttendanceStatus = "excused"
)

// CheckInMethod represents how attendance was recorded
type CheckInMethod string

const (
	CheckInMethodQRScan       CheckInMethod = "qr_scan"
	CheckInMethodManual       CheckInMethod = "manual"
	CheckInMethodNFC          CheckInMethod = "nfc"
	CheckInMethodSelfCheckIn  CheckInMethod = "self_checkin"
)

// EventParticipant represents a participant in an event (either Person or Visitor)
type EventParticipant struct {
	ID             uuid.UUID          `gorm:"type:char(36);primary_key"`
	EventID        uuid.UUID          `gorm:"type:char(36);not null;index:idx_event_participant"`
	Event          Event              `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OccurrenceDate time.Time          `gorm:"type:date;not null;index:idx_event_participant"`

	// Polymorphic participant (Person or Visitor)
	ParticipantType ParticipantType    `gorm:"type:varchar(20);not null;index:idx_event_participant"`
	ParticipantID   uuid.UUID          `gorm:"type:char(36);not null;index:idx_event_participant"`

	// Relationships (only one will be loaded based on ParticipantType)
	// IMPORTANT: Using gorm:"-" to prevent GORM from creating database-level foreign keys
	// Polymorphic relationships CANNOT have FK constraints at the database level because
	// participant_id can reference either people OR visitors table (not both simultaneously)
	// We load these manually in the repository based on participant_type
	Person  *Person  `gorm:"-"`
	Visitor *Visitor `gorm:"-"`

	// Registration and attendance tracking
	RegistrationStatus RegistrationStatus `gorm:"type:varchar(20);not null;default:'registered'"`
	AttendanceStatus   AttendanceStatus   `gorm:"type:varchar(20);not null;default:'absent'"`

	// Check-in/out tracking
	CheckInTime      *time.Time     `gorm:"type:timestamp;null"`
	CheckOutTime     *time.Time     `gorm:"type:timestamp;null"`
	CheckInMethod    *CheckInMethod `gorm:"type:varchar(20);null"`
	QRScanTimestamp  *time.Time     `gorm:"column:qr_scan_timestamp;type:timestamp;null"` // Timestamp from QR code when scanned

	// Additional information
	Notes       string     `gorm:"type:text"`
	RecordedBy  *uuid.UUID `gorm:"type:char(36);null"` // Person who recorded attendance
	RecordedByPerson *Person `gorm:"foreignKey:RecordedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Timestamp
}

// TableName specifies the table name for EventParticipant
func (EventParticipant) TableName() string {
	return "event_participants"
}

// BeforeCreate hook to generate UUID
func (ep *EventParticipant) BeforeCreate(tx *gorm.DB) error {
	if ep.ID == uuid.Nil {
		ep.ID = uuid.New()
	}

	// Set default values
	if ep.RegistrationStatus == "" {
		ep.RegistrationStatus = RegistrationStatusRegistered
	}
	if ep.AttendanceStatus == "" {
		ep.AttendanceStatus = AttendanceStatusAbsent
	}

	return nil
}

// GetPerson returns the Person if participant type is person
func (ep *EventParticipant) GetPerson() *Person {
	if ep.ParticipantType == ParticipantTypePerson {
		return ep.Person
	}
	return nil
}

// GetVisitor returns the Visitor if participant type is visitor
func (ep *EventParticipant) GetVisitor() *Visitor {
	if ep.ParticipantType == ParticipantTypeVisitor {
		return ep.Visitor
	}
	return nil
}

// IsCheckedIn returns true if participant has checked in
func (ep *EventParticipant) IsCheckedIn() bool {
	return ep.CheckInTime != nil
}

// IsCheckedOut returns true if participant has checked out
func (ep *EventParticipant) IsCheckedOut() bool {
	return ep.CheckOutTime != nil
}

// MarkAsPresent marks the participant as present
func (ep *EventParticipant) MarkAsPresent(method CheckInMethod, recordedBy *uuid.UUID) {
	now := time.Now()
	ep.AttendanceStatus = AttendanceStatusPresent
	ep.CheckInTime = &now
	ep.CheckInMethod = &method
	if recordedBy != nil {
		ep.RecordedBy = recordedBy
	}
}

// MarkAsAbsent marks the participant as absent
func (ep *EventParticipant) MarkAsAbsent() {
	ep.AttendanceStatus = AttendanceStatusAbsent
	ep.CheckInTime = nil
	ep.CheckOutTime = nil
	ep.CheckInMethod = nil
}
