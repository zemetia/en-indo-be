package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
)

func TestRecurrenceGenerator_BasicFunctionality(t *testing.T) {
	generator := service.NewRecurrenceGenerator()

	// Test weekly recurrence
	t.Run("Weekly recurrence - every Monday", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC) // Monday
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			ID:            uuid.New(),
			Title:         "Weekly Meeting",
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "WEEKLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should have 5 Mondays in January 2024 (1, 8, 15, 22, 29)
		assert.Len(t, occurrences, 5)
		
		// Verify all are Mondays
		for _, occ := range occurrences {
			assert.Equal(t, time.Monday, occ.Weekday())
		}
	})

	t.Run("Bi-weekly recurrence - every other Wednesday", func(t *testing.T) {
		startDate := time.Date(2024, 1, 3, 14, 0, 0, 0, time.UTC) // Wednesday
		endDate := time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			ID:            uuid.New(),
			Title:         "Bi-weekly Team Meeting",
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(2 * time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "WEEKLY",
			Interval:  2, // Every 2 weeks
			ByWeekday: []string{"WE"},
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should be Jan 3, 17, 31, Feb 14, 28
		assert.Len(t, occurrences, 5)
		
		// Verify 2-week intervals
		for i := 1; i < len(occurrences); i++ {
			diff := occurrences[i].Sub(occurrences[i-1])
			assert.Equal(t, 14*24*time.Hour, diff)
		}
	})

	t.Run("Monthly recurrence - 15th of every month", func(t *testing.T) {
		startDate := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			ID:            uuid.New(),
			Title:         "Monthly Review",
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency:  "MONTHLY",
			Interval:   1,
			ByMonthDay: []int64{15},
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should have 6 occurrences (Jan-Jun)
		assert.Len(t, occurrences, 6)
		
		// Verify all are on the 15th
		for _, occ := range occurrences {
			assert.Equal(t, 15, occ.Day())
		}
	})

	t.Run("Monthly recurrence - 2nd Monday of every month", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 4, 30, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			ID:            uuid.New(),
			Title:         "Board Meeting",
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(2 * time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "MONTHLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
			BySetPos:  []int64{2}, // 2nd occurrence
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should have 4 occurrences (Jan-Apr)
		assert.Len(t, occurrences, 4)
		
		// Verify all are Mondays and in the second week
		for _, occ := range occurrences {
			assert.Equal(t, time.Monday, occ.Weekday())
			// Second Monday should be between 8th and 14th
			assert.True(t, occ.Day() >= 8 && occ.Day() <= 14)
		}
	})
}

func TestRecurrenceGenerator_EdgeCases(t *testing.T) {
	generator := service.NewRecurrenceGenerator()

	t.Run("Count limitation", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		count := 5
		rule := &entity.RecurrenceRule{
			Frequency: "WEEKLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
			Count:     &count,
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should be limited to 5 occurrences
		assert.Len(t, occurrences, 5)
	})

	t.Run("Until date limitation", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		untilDate := time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "WEEKLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
			Until:     &untilDate,
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should only include occurrences until Jan 15 (Jan 1, 8)
		assert.Len(t, occurrences, 2)
		
		for _, occ := range occurrences {
			assert.True(t, occ.Before(untilDate) || occ.Equal(untilDate))
		}
	})

	t.Run("Exception handling - skipped occurrence", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "WEEKLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
		}
		
		// Skip January 8th occurrence
		skipDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
		exceptions := []entity.RecurrenceException{
			{
				ExceptionDate: skipDate,
				IsSkipped:     true,
			},
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, exceptions)
		require.NoError(t, err)
		
		// Should have 4 occurrences instead of 5 (skipping Jan 8)
		assert.Len(t, occurrences, 4)
		
		// Verify Jan 8 is not included
		for _, occ := range occurrences {
			assert.NotEqual(t, 8, occ.Day())
		}
	})
}

func TestRecurrenceGenerator_Validation(t *testing.T) {
	generator := service.NewRecurrenceGenerator()

	t.Run("Valid rules", func(t *testing.T) {
		validRules := []*entity.RecurrenceRule{
			{Frequency: "DAILY", Interval: 1},
			{Frequency: "WEEKLY", Interval: 2, ByWeekday: []string{"MO", "WE", "FR"}},
			{Frequency: "MONTHLY", Interval: 1, ByMonthDay: []int64{1, 15}},
			{Frequency: "YEARLY", Interval: 1, ByMonth: []int64{3, 6, 9, 12}},
		}

		for _, rule := range validRules {
			err := generator.ValidateRecurrenceRule(rule)
			assert.NoError(t, err, "Rule should be valid: %+v", rule)
		}
	})

	t.Run("Invalid rules", func(t *testing.T) {
		invalidRules := []*entity.RecurrenceRule{
			{Frequency: "INVALID", Interval: 1},
			{Frequency: "WEEKLY", Interval: 0},
			{Frequency: "WEEKLY", Interval: 1, ByWeekday: []string{"XX"}},
			{Frequency: "MONTHLY", Interval: 1, ByMonthDay: []int64{32}},
			{Frequency: "YEARLY", Interval: 1, ByMonth: []int64{13}},
		}

		for _, rule := range invalidRules {
			err := generator.ValidateRecurrenceRule(rule)
			assert.Error(t, err, "Rule should be invalid: %+v", rule)
		}
	})
}

func TestEventService_RecurrenceOperations(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)

	t.Run("Create recurring event", func(t *testing.T) {
		req := &dto.CreateEventRequest{
			Title:         "Daily Standup",
			Description:   "Daily team standup meeting",
			EventDate:     "2024-01-01",
			StartTime:     "09:00",
			EndTime:       "09:30",
			EventLocation: "Conference Room A",
			Type:          "event",
			Timezone:      "Asia/Jakarta",
			IsPublic:      true,
			RecurrenceRule: &dto.CreateRecurrenceRuleRequest{
				Frequency: "WEEKLY",
				Interval:  1,
				ByWeekday: []string{"MO", "TU", "WE", "TH", "FR"},
			},
		}

		response, err := eventService.CreateEvent(req)
		require.NoError(t, err)
		require.NotNil(t, response)
		
		assert.Equal(t, req.Title, response.Title)
		assert.NotNil(t, response.RecurrenceRule)
		assert.Equal(t, "WEEKLY", response.RecurrenceRule.Frequency)
		assert.Equal(t, 1, response.RecurrenceRule.Interval)
		assert.Contains(t, response.RecurrenceRule.ByWeekday, "MO")

		// Clean up
		eventService.DeleteEvent(response.ID)
	})

	t.Run("Validate recurrence rule", func(t *testing.T) {
		validRule := &dto.CreateRecurrenceRuleRequest{
			Frequency: "MONTHLY",
			Interval:  2,
			BySetPos:  []int64{1, -1}, // First and last
			ByWeekday: []string{"FR"},
		}

		err := eventService.ValidateRecurrenceRule(validRule)
		assert.NoError(t, err)

		invalidRule := &dto.CreateRecurrenceRuleRequest{
			Frequency: "MONTHLY",
			Interval:  0, // Invalid interval
		}

		err = eventService.ValidateRecurrenceRule(invalidRule)
		assert.Error(t, err)
	})
}

func TestEventService_ThreeTierModifications(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)

	// Create a recurring event for testing
	createReq := &dto.CreateEventRequest{
		Title:         "Weekly Team Meeting",
		Description:   "Weekly team sync",
		EventDate:     "2024-01-01", // Monday
		StartTime:     "10:00",
		EndTime:       "11:00",
		EventLocation: "Meeting Room 1",
		Type:          "event",
		Timezone:      "Asia/Jakarta",
		IsPublic:      false,
		RecurrenceRule: &dto.CreateRecurrenceRuleRequest{
			Frequency: "WEEKLY",
			Interval:  1,
			ByWeekday: []string{"MO"},
		},
	}

	event, err := eventService.CreateEvent(createReq)
	require.NoError(t, err)
	defer eventService.DeleteEvent(event.ID)

	t.Run("Update single occurrence", func(t *testing.T) {
		req := &dto.UpdateOccurrenceRequest{
			OccurrenceDate: "2024-01-08", // Second Monday
			StartTime:      strPtr("14:00"),
			EndTime:        strPtr("15:00"),
			Event: dto.UpdateEventRequest{
				Title: strPtr("Special Team Meeting"),
			},
		}

		err := eventService.UpdateSingleOccurrence(event.ID, req)
		assert.NoError(t, err)

		// Get occurrences to verify the change
		occReq := &dto.GetEventOccurrencesRequest{
			StartDate: "2024-01-01",
			EndDate:   "2024-01-31",
		}
		
		occurrences, err := eventService.GetEventOccurrences(event.ID, occReq)
		require.NoError(t, err)
		
		// Find the modified occurrence
		var modifiedOccurrence *dto.EventOccurrenceResponse
		for i := range occurrences {
			if occurrences[i].OccurrenceDate.Format("2006-01-02") == "2024-01-08" {
				modifiedOccurrence = &occurrences[i]
				break
			}
		}
		
		require.NotNil(t, modifiedOccurrence)
		assert.True(t, modifiedOccurrence.IsException)
		assert.Equal(t, 14, modifiedOccurrence.StartDatetime.Hour())
		assert.Equal(t, 15, modifiedOccurrence.EndDatetime.Hour())
	})

	t.Run("Update future occurrences", func(t *testing.T) {
		req := &dto.UpdateFutureOccurrencesRequest{
			FromDate:  "2024-01-15", // Third Monday
			StartTime: strPtr("15:00"),
			EndTime:   strPtr("16:30"),
			Event: dto.UpdateEventRequest{
				Title:       strPtr("Updated Team Meeting"),
				Description: strPtr("Updated weekly team sync with new format"),
			},
		}

		err := eventService.UpdateFutureOccurrences(event.ID, req)
		assert.NoError(t, err)

		// Verify that a new series was created
		// The original series should be limited to before Jan 15
		occReq := &dto.GetEventOccurrencesRequest{
			StartDate: "2024-01-01",
			EndDate:   "2024-02-28",
		}
		
		allOccurrences, err := eventService.GetOccurrencesInRange(occReq)
		require.NoError(t, err)
		
		// Should have both original series (up to Jan 8) and new series (from Jan 15)
		var originalCount, newCount int
		for _, occ := range allOccurrences {
			if occ.OriginalEvent.Title == "Weekly Team Meeting" {
				originalCount++
			} else if occ.OriginalEvent.Title == "Updated Team Meeting" {
				newCount++
			}
		}
		
		assert.Greater(t, originalCount, 0, "Should have original occurrences")
		assert.Greater(t, newCount, 0, "Should have new series occurrences")
	})

	t.Run("Delete single occurrence", func(t *testing.T) {
		req := &dto.DeleteOccurrenceRequest{
			OccurrenceDate: "2024-01-22", // A Monday to delete
			DeleteType:     "single",
		}

		err := eventService.DeleteOccurrence(event.ID, req)
		assert.NoError(t, err)

		// Verify the occurrence is skipped
		occReq := &dto.GetEventOccurrencesRequest{
			StartDate: "2024-01-20",
			EndDate:   "2024-01-25",
		}
		
		occurrences, err := eventService.GetEventOccurrences(event.ID, occReq)
		require.NoError(t, err)
		
		// Should not find Jan 22 occurrence
		for _, occ := range occurrences {
			assert.NotEqual(t, "2024-01-22", occ.OccurrenceDate.Format("2006-01-02"))
		}
	})

	t.Run("Get next occurrence", func(t *testing.T) {
		after := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		
		nextOcc, err := eventService.GetNextOccurrence(event.ID, after)
		require.NoError(t, err)
		require.NotNil(t, nextOcc)
		
		assert.True(t, nextOcc.After(after))
		assert.Equal(t, time.Monday, nextOcc.Weekday())
	})
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}

func TestRecurrenceGenerator_ComplexPatterns(t *testing.T) {
	generator := service.NewRecurrenceGenerator()

	t.Run("Last Friday of every month", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency: "MONTHLY",
			Interval:  1,
			ByWeekday: []string{"FR"},
			BySetPos:  []int64{-1}, // Last occurrence
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should have 6 occurrences (Jan-Jun)
		assert.Len(t, occurrences, 6)
		
		// Verify all are Fridays and are the last Friday of their respective months
		for _, occ := range occurrences {
			assert.Equal(t, time.Friday, occ.Weekday())
			
			// Check it's the last Friday of the month
			nextWeek := occ.AddDate(0, 0, 7)
			assert.NotEqual(t, occ.Month(), nextWeek.Month(), 
				"Should be the last Friday of month for %s", occ.Format("2006-01-02"))
		}
	})

	t.Run("Every 3 months on the 15th", func(t *testing.T) {
		startDate := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		
		event := &entity.Event{
			StartDatetime: startDate,
			EndDatetime:   startDate.Add(time.Hour),
		}
		
		rule := &entity.RecurrenceRule{
			Frequency:  "MONTHLY",
			Interval:   3, // Every 3 months
			ByMonthDay: []int64{15},
		}
		
		occurrences, err := generator.GenerateOccurrences(event, rule, startDate, endDate, nil)
		require.NoError(t, err)
		
		// Should have 4 occurrences: Jan 15, Apr 15, Jul 15, Oct 15
		assert.Len(t, occurrences, 4)
		
		expectedMonths := []time.Month{time.January, time.April, time.July, time.October}
		for i, occ := range occurrences {
			assert.Equal(t, expectedMonths[i], occ.Month())
			assert.Equal(t, 15, occ.Day())
		}
	})
}