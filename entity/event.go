package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Event struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Title       string    `gorm:"type:varchar(255);not null"`
	BannerImage string    `gorm:"type:varchar(255);null"`
	Description string    `gorm:"type:text"`
	Capacity    int       `gorm:"default:99999"`
	Type        string    `gorm:"type:varchar(255);not null"` // event, ibadah, spiritual journey

	EventDate        time.Time       `gorm:"type:datetime;not null"`
	EventLocation    string          `gorm:"type:varchar(255);not null"`
	StartDatetime    time.Time       `gorm:"not null"`
	EndDatetime      time.Time       `gorm:"not null"`
	AllDay           bool            `gorm:"default:false"`
	Timezone         string          `gorm:"type:varchar(64);not null"`
	RecurrenceRuleID *uuid.UUID      `gorm:"type:char(36);index"`
	RecurrenceRule   *RecurrenceRule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsPublic              bool                 `gorm:"default:false"`
	DiscipleshipJourneyID *uuid.UUID           `gorm:"type:char(36);index"`
	DiscipleshipJourney   *DiscipleshipJourney `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Lagu []Lagu `gorm:"many2many:event_lagu;"`

	Timestamp
}

type RecurrenceRule struct {
	ID         uuid.UUID      `gorm:"type:char(36);primary_key"`
	Frequency  string         `gorm:"type:varchar(16);not null"` // DAILY, WEEKLY, MONTHLY, YEARLY
	Interval   int            `gorm:"default:1;not null"`        // multiples of frequency
	ByWeekday  pq.StringArray `gorm:"type:text[]"`               // e.g. ["MO","TU"], only for WEEKLY
	ByMonthDay pq.Int64Array  `gorm:"type:integer[]"`            // e.g. [1,15], only for MONTHLY
	ByMonth    pq.Int64Array  `gorm:"type:integer[]"`            // e.g. [1,6,12], only for YEARLY
	BySetPos   pq.Int64Array  `gorm:"type:integer[]"`            // e.g. [1,-1] for first/last occurrence
	WeekStart  string         `gorm:"type:varchar(2);default:'MO'"` // week start day (MO, SU, etc.)
	ByYearDay  pq.Int64Array  `gorm:"type:integer[]"`            // day of year (1-366)
	Count      *int           `gorm:""`                          // optional: limit total occurrences
	Until      *time.Time     `gorm:""`                          // optional: end date for occurrences

	Timestamp
}

type RecurrenceException struct {
	ID                uuid.UUID  `gorm:"type:char(36);primary_key"`
	EventID           uuid.UUID  `gorm:"type:char(36);index;not null"`
	Event             Event      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:EventID"`
	ExceptionDate     time.Time  `gorm:"type:date;not null"`       // date of the instance to skip or override
	ModificationType  string     `gorm:"type:varchar(10);default:'single'"` // single, future, all
	IsSkipped         bool       `gorm:"default:false"`            // true if this occurrence is removed
	OverrideStart     *time.Time `gorm:""`                         // optional: new start time for this instance
	OverrideEnd       *time.Time `gorm:""`                         // optional: new end time for this instance
	OriginalStartTime *time.Time `gorm:""`                         // original start time before modification
	OriginalEndTime   *time.Time `gorm:""`                         // original end time before modification
	SplitFromDate     *time.Time `gorm:"type:date"`                // when "future" edits split the series
	Notes             string     `gorm:"type:text"`                // optional remarks

	Timestamp
}

// ModificationTypes for recurrence exceptions
const (
	ModificationTypeSingle string = "single"
	ModificationTypeFuture string = "future"
	ModificationTypeAll    string = "all"
)
