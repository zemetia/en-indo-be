package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/zemetia/en-indo-be/entity"
)

// RecurrenceGenerator handles generation of recurring event occurrences
type RecurrenceGenerator struct{}

// NewRecurrenceGenerator creates a new RecurrenceGenerator instance
func NewRecurrenceGenerator() *RecurrenceGenerator {
	return &RecurrenceGenerator{}
}

// GenerateOccurrences generates event occurrences for a given date range
func (rg *RecurrenceGenerator) GenerateOccurrences(
	event *entity.Event,
	rule *entity.RecurrenceRule,
	startDate, endDate time.Time,
	exceptions []entity.RecurrenceException,
) ([]time.Time, error) {
	if rule == nil {
		return []time.Time{event.StartDatetime}, nil
	}

	var occurrences []time.Time
	current := event.StartDatetime

	// Adjust current to start boundary
	if current.Before(startDate) {
		current = startDate
		// Adjust to the next valid occurrence based on the rule
		current = rg.adjustToNextValidOccurrence(current, event.StartDatetime, rule)
	}

	// Generate occurrences based on frequency
	switch strings.ToUpper(rule.Frequency) {
	case "DAILY":
		occurrences = rg.generateDailyOccurrences(current, endDate, rule, event.StartDatetime)
	case "WEEKLY":
		occurrences = rg.generateWeeklyOccurrences(current, endDate, rule, event.StartDatetime)
	case "MONTHLY":
		occurrences = rg.generateMonthlyOccurrences(current, endDate, rule, event.StartDatetime)
	case "YEARLY":
		occurrences = rg.generateYearlyOccurrences(current, endDate, rule, event.StartDatetime)
	default:
		return nil, fmt.Errorf("unsupported frequency: %s", rule.Frequency)
	}

	// Apply count limitation if specified
	if rule.Count != nil && len(occurrences) > *rule.Count {
		occurrences = occurrences[:*rule.Count]
	}

	// Apply until date limitation if specified
	if rule.Until != nil {
		var filteredOccurrences []time.Time
		for _, occ := range occurrences {
			if occ.After(*rule.Until) {
				break
			}
			filteredOccurrences = append(filteredOccurrences, occ)
		}
		occurrences = filteredOccurrences
	}

	// Apply exceptions
	occurrences = rg.applyExceptions(occurrences, exceptions)

	// Filter by date range
	var result []time.Time
	for _, occ := range occurrences {
		if (occ.Equal(startDate) || occ.After(startDate)) && (occ.Equal(endDate) || occ.Before(endDate)) {
			result = append(result, occ)
		}
	}

	return result, nil
}

// generateDailyOccurrences generates daily recurring occurrences
func (rg *RecurrenceGenerator) generateDailyOccurrences(start, end time.Time, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time
	current := start
	interval := rule.Interval
	if interval <= 0 {
		interval = 1
	}

	for current.Before(end) || current.Equal(end) {
		occurrences = append(occurrences, current)
		current = current.AddDate(0, 0, interval)
	}

	return occurrences
}

// generateWeeklyOccurrences generates weekly recurring occurrences
func (rg *RecurrenceGenerator) generateWeeklyOccurrences(start, end time.Time, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time
	interval := rule.Interval
	if interval <= 0 {
		interval = 1
	}

	weekdays := rg.jsonToStringSlice(rule.ByWeekday)
	if len(weekdays) == 0 {
		// Default to the original event's weekday
		weekdays = []string{rg.getWeekdayAbbreviation(originalStart.Weekday())}
	}

	// Convert weekday abbreviations to time.Weekday
	targetWeekdays := make([]time.Weekday, 0, len(weekdays))
	for _, wd := range weekdays {
		if weekday, ok := rg.parseWeekday(wd); ok {
			targetWeekdays = append(targetWeekdays, weekday)
		}
	}

	if len(targetWeekdays) == 0 {
		return occurrences
	}

	// Start from the beginning of the week containing start date
	weekStart := rg.getWeekStart(start, rule.WeekStart)
	current := weekStart

	for current.Before(end) || current.Equal(end) {
		// Check each day in the current week
		for i := 0; i < 7; i++ {
			currentDay := current.AddDate(0, 0, i)
			if currentDay.After(end) {
				break
			}

			// Check if this weekday is in our target list
			for _, targetWeekday := range targetWeekdays {
				if currentDay.Weekday() == targetWeekday {
					// Adjust time to match original event time
					occurrence := time.Date(
						currentDay.Year(), currentDay.Month(), currentDay.Day(),
						originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
						originalStart.Nanosecond(), originalStart.Location(),
					)

					if (occurrence.Equal(start) || occurrence.After(start)) &&
						(occurrence.Equal(end) || occurrence.Before(end)) {
						occurrences = append(occurrences, occurrence)
					}
					break
				}
			}
		}
		// Move to next interval week
		current = current.AddDate(0, 0, 7*interval)
	}

	return occurrences
}

// generateMonthlyOccurrences generates monthly recurring occurrences
func (rg *RecurrenceGenerator) generateMonthlyOccurrences(start, end time.Time, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time
	interval := rule.Interval
	if interval <= 0 {
		interval = 1
	}

	current := time.Date(originalStart.Year(), originalStart.Month(), 1,
		originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
		originalStart.Nanosecond(), originalStart.Location())

	// Advance to start month if needed
	for current.Before(start) {
		current = current.AddDate(0, interval, 0)
	}

	for current.Before(end) || current.Equal(end) {
		monthOccurrences := rg.generateMonthOccurrences(current, rule, originalStart)
		for _, occ := range monthOccurrences {
			if (occ.Equal(start) || occ.After(start)) && (occ.Equal(end) || occ.Before(end)) {
				occurrences = append(occurrences, occ)
			}
		}
		current = current.AddDate(0, interval, 0)
	}

	return occurrences
}

// generateMonthOccurrences generates occurrences for a specific month
func (rg *RecurrenceGenerator) generateMonthOccurrences(monthStart time.Time, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time

	// Handle ByMonthDay (specific days of month)
	if len(rule.ByMonthDay) > 0 {
		for _, day := range rule.ByMonthDay {
			dayInt := int(day)
			if dayInt < 0 {
				// Negative values count from end of month
				daysInMonth := rg.daysInMonth(monthStart.Year(), monthStart.Month())
				dayInt = daysInMonth + dayInt + 1
			}

			if dayInt > 0 && dayInt <= rg.daysInMonth(monthStart.Year(), monthStart.Month()) {
				occurrence := time.Date(
					monthStart.Year(), monthStart.Month(), dayInt,
					originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
					originalStart.Nanosecond(), originalStart.Location(),
				)
				occurrences = append(occurrences, occurrence)
			}
		}
		return occurrences
	}

	// Handle ByWeekday with BySetPos (e.g., 2nd Monday, last Friday)
	weekdays := rg.jsonToStringSlice(rule.ByWeekday)
	if len(weekdays) > 0 {
		weekdayOccurrences := rg.getWeekdayOccurrencesInMonth(monthStart, weekdays, originalStart)

		bySetPos := rg.jsonToInt64Slice(rule.BySetPos)
		if len(bySetPos) > 0 {
			// Apply BySetPos filtering
			var filtered []time.Time
			for _, pos := range bySetPos {
				posInt := int(pos)
				if posInt > 0 && posInt <= len(weekdayOccurrences) {
					filtered = append(filtered, weekdayOccurrences[posInt-1])
				} else if posInt < 0 && abs(posInt) <= len(weekdayOccurrences) {
					filtered = append(filtered, weekdayOccurrences[len(weekdayOccurrences)+posInt])
				}
			}
			return filtered
		}
		return weekdayOccurrences
	}

	// Default: use original day of month
	dayOfMonth := originalStart.Day()
	if dayOfMonth <= rg.daysInMonth(monthStart.Year(), monthStart.Month()) {
		occurrence := time.Date(
			monthStart.Year(), monthStart.Month(), dayOfMonth,
			originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
			originalStart.Nanosecond(), originalStart.Location(),
		)
		occurrences = append(occurrences, occurrence)
	}

	return occurrences
}

// generateYearlyOccurrences generates yearly recurring occurrences
func (rg *RecurrenceGenerator) generateYearlyOccurrences(start, end time.Time, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time
	interval := rule.Interval
	if interval <= 0 {
		interval = 1
	}

	currentYear := originalStart.Year()
	for currentYear <= end.Year() {
		yearStart := time.Date(currentYear, time.January, 1,
			originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
			originalStart.Nanosecond(), originalStart.Location())

		if yearStart.After(end) {
			break
		}

		yearOccurrences := rg.generateYearOccurrences(currentYear, rule, originalStart)
		for _, occ := range yearOccurrences {
			if (occ.Equal(start) || occ.After(start)) && (occ.Equal(end) || occ.Before(end)) {
				occurrences = append(occurrences, occ)
			}
		}

		currentYear += interval
	}

	return occurrences
}

// generateYearOccurrences generates occurrences for a specific year
func (rg *RecurrenceGenerator) generateYearOccurrences(year int, rule *entity.RecurrenceRule, originalStart time.Time) []time.Time {
	var occurrences []time.Time

	months := rg.jsonToInt64Slice(rule.ByMonth)
	if len(months) == 0 {
		// Default to original month
		months = []int64{int64(originalStart.Month())}
	}

	for _, month := range months {
		monthInt := int(month)
		if monthInt < 1 || monthInt > 12 {
			continue
		}

		monthStart := time.Date(year, time.Month(monthInt), 1,
			originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
			originalStart.Nanosecond(), originalStart.Location())

		monthOccurrences := rg.generateMonthOccurrences(monthStart, rule, originalStart)
		occurrences = append(occurrences, monthOccurrences...)
	}

	return occurrences
}

// Helper functions

func (rg *RecurrenceGenerator) adjustToNextValidOccurrence(current, originalStart time.Time, rule *entity.RecurrenceRule) time.Time {
	// This would contain logic to adjust current time to next valid occurrence
	// For now, return current as is
	return current
}

func (rg *RecurrenceGenerator) applyExceptions(occurrences []time.Time, exceptions []entity.RecurrenceException) []time.Time {
	if len(exceptions) == 0 {
		return occurrences
	}

	exceptionDates := make(map[string]bool)
	for _, exception := range exceptions {
		if exception.IsSkipped {
			dateKey := exception.ExceptionDate.Format("2006-01-02")
			exceptionDates[dateKey] = true
		}
	}

	var filtered []time.Time
	for _, occ := range occurrences {
		dateKey := occ.Format("2006-01-02")
		if !exceptionDates[dateKey] {
			filtered = append(filtered, occ)
		}
	}

	return filtered
}

func (rg *RecurrenceGenerator) getWeekStart(date time.Time, weekStart string) time.Time {
	// Default to Monday if not specified
	if weekStart == "" {
		weekStart = "MO"
	}

	targetWeekday, ok := rg.parseWeekday(weekStart)
	if !ok {
		targetWeekday = time.Monday
	}

	// Calculate days to subtract to get to week start
	daysFromWeekStart := (int(date.Weekday()) - int(targetWeekday) + 7) % 7
	return date.AddDate(0, 0, -daysFromWeekStart)
}

func (rg *RecurrenceGenerator) getWeekdayAbbreviation(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "SU",
		time.Monday:    "MO",
		time.Tuesday:   "TU",
		time.Wednesday: "WE",
		time.Thursday:  "TH",
		time.Friday:    "FR",
		time.Saturday:  "SA",
	}
	return weekdays[weekday]
}

func (rg *RecurrenceGenerator) parseWeekday(abbr string) (time.Weekday, bool) {
	weekdays := map[string]time.Weekday{
		"SU": time.Sunday,
		"MO": time.Monday,
		"TU": time.Tuesday,
		"WE": time.Wednesday,
		"TH": time.Thursday,
		"FR": time.Friday,
		"SA": time.Saturday,
	}
	weekday, ok := weekdays[strings.ToUpper(abbr)]
	return weekday, ok
}

func (rg *RecurrenceGenerator) daysInMonth(year int, month time.Month) int {
	// Get the first day of next month, then subtract one day to get last day of current month
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

func (rg *RecurrenceGenerator) getWeekdayOccurrencesInMonth(monthStart time.Time, weekdays []string, originalStart time.Time) []time.Time {
	var occurrences []time.Time
	daysInMonth := rg.daysInMonth(monthStart.Year(), monthStart.Month())

	for _, weekdayStr := range weekdays {
		weekday, ok := rg.parseWeekday(weekdayStr)
		if !ok {
			continue
		}

		// Find all occurrences of this weekday in the month
		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(
				monthStart.Year(), monthStart.Month(), day,
				originalStart.Hour(), originalStart.Minute(), originalStart.Second(),
				originalStart.Nanosecond(), originalStart.Location(),
			)
			if date.Weekday() == weekday {
				occurrences = append(occurrences, date)
			}
		}
	}

	// Sort occurrences by date
	sort.Slice(occurrences, func(i, j int) bool {
		return occurrences[i].Before(occurrences[j])
	})

	return occurrences
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ValidateRecurrenceRule validates a recurrence rule for correctness
func (rg *RecurrenceGenerator) ValidateRecurrenceRule(rule *entity.RecurrenceRule) error {
	if rule == nil {
		return fmt.Errorf("recurrence rule cannot be nil")
	}

	// Validate frequency
	validFrequencies := map[string]bool{
		"DAILY":   true,
		"WEEKLY":  true,
		"MONTHLY": true,
		"YEARLY":  true,
	}
	if !validFrequencies[strings.ToUpper(rule.Frequency)] {
		return fmt.Errorf("invalid frequency: %s", rule.Frequency)
	}

	// Validate interval
	if rule.Interval < 1 {
		return fmt.Errorf("interval must be greater than 0")
	}

	// Validate weekdays
	weekdays := rg.jsonToStringSlice(rule.ByWeekday)
	for _, weekday := range weekdays {
		if _, ok := rg.parseWeekday(weekday); !ok {
			return fmt.Errorf("invalid weekday: %s", weekday)
		}
	}

	// Validate month days
	for _, day := range rule.ByMonthDay {
		if day < -31 || day > 31 || day == 0 {
			return fmt.Errorf("invalid month day: %d", day)
		}
	}

	// Validate months
	for _, month := range rule.ByMonth {
		if month < 1 || month > 12 {
			return fmt.Errorf("invalid month: %d", month)
		}
	}

	// Validate set positions
	for _, pos := range rule.BySetPos {
		if pos < -53 || pos > 53 || pos == 0 {
			return fmt.Errorf("invalid set position: %d", pos)
		}
	}

	// Validate year days
	for _, yearDay := range rule.ByYearDay {
		if yearDay < -366 || yearDay > 366 || yearDay == 0 {
			return fmt.Errorf("invalid year day: %d", yearDay)
		}
	}

	return nil
}

// GetNextOccurrence gets the next occurrence after a given date
func (rg *RecurrenceGenerator) GetNextOccurrence(event *entity.Event, rule *entity.RecurrenceRule, after time.Time) (*time.Time, error) {
	if rule == nil {
		if event.StartDatetime.After(after) {
			return &event.StartDatetime, nil
		}
		return nil, nil
	}

	// Generate occurrences for a reasonable range (next year)
	endDate := after.AddDate(1, 0, 0)
	occurrences, err := rg.GenerateOccurrences(event, rule, after.AddDate(0, 0, 1), endDate, nil)
	if err != nil {
		return nil, err
	}

	if len(occurrences) > 0 {
		return &occurrences[0], nil
	}

	return nil, nil
}

// Helper functions to convert JSON strings to slices
func (rg *RecurrenceGenerator) jsonToStringSlice(jsonStr string) []string {
	if jsonStr == "" {
		return []string{}
	}

	var result []string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return []string{}
	}

	return result
}

func (rg *RecurrenceGenerator) jsonToInt64Slice(jsonStr string) []int64 {
	if jsonStr == "" {
		return []int64{}
	}

	var result []int64
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return []int64{}
	}

	return result
}
