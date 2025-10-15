package entity

import "github.com/google/uuid"

// EventDepartment represents the many-to-many relationship between events and departments
// This allows PICs to specify which departments are needed for an event
type EventDepartment struct {
	ID           uuid.UUID  `gorm:"type:char(36);primary_key" json:"id"`
	EventID      uuid.UUID  `gorm:"type:char(36);not null;index" json:"event_id"`
	DepartmentID uuid.UUID  `gorm:"type:char(36);not null;index" json:"department_id"`
	Event        Event      `gorm:"foreignKey:EventID"`
	Department   Department `gorm:"foreignKey:DepartmentID"`

	Timestamp
}
