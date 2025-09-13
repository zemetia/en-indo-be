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

func TestEventPICService_BasicOperations(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventPICRepo := repository.NewEventPICRepository(db)
	eventService := service.NewEventService(eventRepo, eventPICRepo)
	eventPICService := service.NewEventPICService(eventPICRepo, eventRepo)

	// Create a test event first
	eventReq := &dto.CreateEventRequest{
		Title:         "Test Event for PIC",
		Description:   "Test event to test PIC functionality",
		EventDate:     "2024-12-01",
		StartTime:     "10:00",
		EndTime:       "12:00",
		EventLocation: "Test Location",
		Type:          "event",
		Timezone:      "Asia/Jakarta",
		IsPublic:      true,
	}

	testEvent, err := eventService.CreateEvent(eventReq)
	require.NoError(t, err)
	defer eventService.DeleteEvent(testEvent.ID)

	// Create a test person (mock)
	testPersonID := uuid.New()
	createdBy := uuid.New()

	t.Run("Create Event PIC", func(t *testing.T) {
		picReq := &dto.CreateEventPICRequest{
			PersonID:          testPersonID,
			Role:              "Primary PIC",
			Description:       "Main organizer for the event",
			IsPrimary:         true,
			StartDate:         "2024-11-01",
			EndDate:           nil,
			CanEdit:           true,
			CanDelete:         true,
			CanAssignPIC:      true,
			NotifyOnChanges:   true,
			NotifyOnReminders: true,
		}

		pic, err := eventPICService.CreateEventPIC(testEvent.ID, picReq, createdBy)
		require.NoError(t, err)
		require.NotNil(t, pic)

		assert.Equal(t, testEvent.ID, pic.EventID)
		assert.Equal(t, testPersonID, pic.PersonID)
		assert.Equal(t, "Primary PIC", pic.Role)
		assert.True(t, pic.IsPrimary)
		assert.True(t, pic.IsActive)
		assert.True(t, pic.CanEdit)
		assert.True(t, pic.CanDelete)
		assert.True(t, pic.CanAssignPIC)

		// Clean up
		eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
	})

	t.Run("Prevent Multiple Primary PICs", func(t *testing.T) {
		// Create first primary PIC
		picReq1 := &dto.CreateEventPICRequest{
			PersonID:    testPersonID,
			Role:        "Primary PIC",
			IsPrimary:   true,
			StartDate:   "2024-11-01",
			CanEdit:     true,
			CanDelete:   true,
			CanAssignPIC: true,
		}

		pic1, err := eventPICService.CreateEventPIC(testEvent.ID, picReq1, createdBy)
		require.NoError(t, err)

		// Try to create second primary PIC - should fail
		picReq2 := &dto.CreateEventPICRequest{
			PersonID:    uuid.New(),
			Role:        "Another Primary PIC",
			IsPrimary:   true,
			StartDate:   "2024-11-01",
			CanEdit:     true,
			CanDelete:   false,
			CanAssignPIC: false,
		}

		_, err = eventPICService.CreateEventPIC(testEvent.ID, picReq2, createdBy)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event already has a primary PIC")

		// Clean up
		eventPICService.DeleteEventPIC(pic1.ID, createdBy, "Test cleanup")
	})

	t.Run("Update Event PIC", func(t *testing.T) {
		// Create PIC
		picReq := &dto.CreateEventPICRequest{
			PersonID:    testPersonID,
			Role:        "Test PIC",
			IsPrimary:   false,
			StartDate:   "2024-11-01",
			CanEdit:     false,
			CanDelete:   false,
			CanAssignPIC: false,
		}

		pic, err := eventPICService.CreateEventPIC(testEvent.ID, picReq, createdBy)
		require.NoError(t, err)

		// Update PIC
		updateReq := &dto.UpdateEventPICRequest{
			Role:         stringPtr("Updated PIC Role"),
			Description:  stringPtr("Updated description"),
			CanEdit:      boolPtr(true),
			CanDelete:    boolPtr(true),
			CanAssignPIC: boolPtr(true),
		}

		updatedPIC, err := eventPICService.UpdateEventPIC(pic.ID, updateReq, createdBy)
		require.NoError(t, err)

		assert.Equal(t, "Updated PIC Role", updatedPIC.Role)
		assert.Equal(t, "Updated description", updatedPIC.Description)
		assert.True(t, updatedPIC.CanEdit)
		assert.True(t, updatedPIC.CanDelete)
		assert.True(t, updatedPIC.CanAssignPIC)

		// Clean up
		eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
	})
}

func TestEventPICService_AdvancedOperations(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventPICRepo := repository.NewEventPICRepository(db)
	eventService := service.NewEventService(eventRepo, eventPICRepo)
	eventPICService := service.NewEventPICService(eventPICRepo, eventRepo)

	// Create a test event
	eventReq := &dto.CreateEventRequest{
		Title:         "Advanced PIC Test Event",
		Description:   "Test event for advanced PIC functionality",
		EventDate:     "2024-12-15",
		StartTime:     "14:00",
		EndTime:       "16:00",
		EventLocation: "Advanced Test Location",
		Type:          "event",
		Timezone:      "Asia/Jakarta",
		IsPublic:      false,
	}

	testEvent, err := eventService.CreateEvent(eventReq)
	require.NoError(t, err)
	defer eventService.DeleteEvent(testEvent.ID)

	createdBy := uuid.New()

	t.Run("Bulk Assign PICs", func(t *testing.T) {
		bulkReq := &dto.BulkAssignEventPICRequest{
			PICs: []dto.CreateEventPICRequest{
				{
					PersonID:    uuid.New(),
					Role:        "Primary PIC",
					IsPrimary:   true,
					StartDate:   "2024-11-01",
					CanEdit:     true,
					CanDelete:   true,
					CanAssignPIC: true,
				},
				{
					PersonID:    uuid.New(),
					Role:        "Technical PIC",
					IsPrimary:   false,
					StartDate:   "2024-11-01",
					CanEdit:     true,
					CanDelete:   false,
					CanAssignPIC: false,
				},
				{
					PersonID:    uuid.New(),
					Role:        "Logistics PIC",
					IsPrimary:   false,
					StartDate:   "2024-11-01",
					CanEdit:     false,
					CanDelete:   false,
					CanAssignPIC: false,
				},
			},
		}

		err := eventPICService.AssignMultiplePICs(testEvent.ID, bulkReq, createdBy)
		require.NoError(t, err)

		// Verify PICs were created
		pics, err := eventPICService.GetPICsByEventID(testEvent.ID)
		require.NoError(t, err)
		assert.Len(t, pics, 3)

		// Check primary PIC exists
		primaryPIC, err := eventPICService.GetPrimaryPICByEventID(testEvent.ID)
		require.NoError(t, err)
		assert.Equal(t, "Primary PIC", primaryPIC.Role)
		assert.True(t, primaryPIC.IsPrimary)

		// Clean up
		for _, pic := range pics {
			eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
		}
	})

	t.Run("Transfer PIC Role", func(t *testing.T) {
		fromPersonID := uuid.New()
		toPersonID := uuid.New()

		// Create initial PIC
		picReq := &dto.CreateEventPICRequest{
			PersonID:    fromPersonID,
			Role:        "Primary PIC",
			IsPrimary:   true,
			StartDate:   "2024-11-01",
			CanEdit:     true,
			CanDelete:   true,
			CanAssignPIC: true,
		}

		pic, err := eventPICService.CreateEventPIC(testEvent.ID, picReq, createdBy)
		require.NoError(t, err)

		// Transfer PIC role
		transferReq := &dto.TransferEventPICRequest{
			FromPersonID:  fromPersonID,
			ToPersonID:    toPersonID,
			TransferType:  "replace",
			Reason:        "Person no longer available",
			EffectiveDate: "2024-12-01",
		}

		err = eventPICService.TransferPICRole(testEvent.ID, transferReq, createdBy)
		require.NoError(t, err)

		// Verify transfer worked
		pics, err := eventPICService.GetActivePICsByEventID(testEvent.ID)
		require.NoError(t, err)
		
		// Should have new PIC for toPersonID and old PIC should be inactive
		activePIC := pics[0]
		assert.Equal(t, toPersonID, activePIC.PersonID)
		assert.True(t, activePIC.IsPrimary)

		// Clean up
		allPICs, _ := eventPICService.GetPICsByEventID(testEvent.ID)
		for _, pic := range allPICs {
			eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
		}
	})

	t.Run("Get Expiring PICs", func(t *testing.T) {
		// Create PIC with end date in 15 days
		endDate := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
		picReq := &dto.CreateEventPICRequest{
			PersonID:    uuid.New(),
			Role:        "Temporary PIC",
			IsPrimary:   false,
			StartDate:   "2024-11-01",
			EndDate:     &endDate,
			CanEdit:     true,
			CanDelete:   false,
			CanAssignPIC: false,
		}

		pic, err := eventPICService.CreateEventPIC(testEvent.ID, picReq, createdBy)
		require.NoError(t, err)

		// Get expiring PICs within 30 days
		expiringPICs, err := eventPICService.GetExpiringPICs(30)
		require.NoError(t, err)

		// Should find our PIC
		found := false
		for _, expiringPIC := range expiringPICs {
			if expiringPIC.ID == pic.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find expiring PIC")

		// Clean up
		eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
	})
}

func TestEventPICRoles(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventPICRepo := repository.NewEventPICRepository(db)
	eventPICService := service.NewEventPICService(eventPICRepo, eventRepo)

	t.Run("Create and Manage PIC Roles", func(t *testing.T) {
		// Create new role
		roleReq := &dto.CreateEventPICRoleRequest{
			Name:                "Test Custom Role",
			Description:         "Custom role for testing",
			DefaultCanEdit:      true,
			DefaultCanDelete:    false,
			DefaultCanAssignPIC: false,
		}

		role, err := eventPICService.CreateEventPICRole(roleReq)
		require.NoError(t, err)
		assert.Equal(t, "Test Custom Role", role.Name)
		assert.True(t, role.DefaultCanEdit)
		assert.False(t, role.DefaultCanDelete)

		// Update role
		updateReq := &dto.UpdateEventPICRoleRequest{
			Description:         stringPtr("Updated description"),
			DefaultCanDelete:    boolPtr(true),
			DefaultCanAssignPIC: boolPtr(true),
		}

		updatedRole, err := eventPICService.UpdateEventPICRole(role.ID, updateReq)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", updatedRole.Description)
		assert.True(t, updatedRole.DefaultCanDelete)
		assert.True(t, updatedRole.DefaultCanAssignPIC)

		// List roles
		rolesList, err := eventPICService.ListEventPICRoles(1, 10, "")
		require.NoError(t, err)
		assert.Greater(t, rolesList.TotalCount, 0)

		// Search roles
		searchResults, err := eventPICService.ListEventPICRoles(1, 10, "Test Custom")
		require.NoError(t, err)
		found := false
		for _, r := range searchResults.Roles {
			if r.ID == role.ID {
				found = true
				break
			}
		}
		assert.True(t, found)

		// Delete role
		err = eventPICService.DeleteEventPICRole(role.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = eventPICService.GetEventPICRole(role.ID)
		assert.Error(t, err)
	})
}

func TestEventPICPermissions(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventPICRepo := repository.NewEventPICRepository(db)
	eventService := service.NewEventService(eventRepo, eventPICRepo)
	eventPICService := service.NewEventPICService(eventPICRepo, eventRepo)

	// Create test event
	eventReq := &dto.CreateEventRequest{
		Title:         "Permission Test Event",
		Description:   "Event for testing PIC permissions",
		EventDate:     "2024-12-20",
		StartTime:     "09:00",
		EndTime:       "11:00",
		EventLocation: "Permission Test Location",
		Type:          "event",
		Timezone:      "Asia/Jakarta",
		IsPublic:      true,
	}

	testEvent, err := eventService.CreateEvent(eventReq)
	require.NoError(t, err)
	defer eventService.DeleteEvent(testEvent.ID)

	personID := uuid.New()
	createdBy := uuid.New()

	t.Run("Validate PIC Permissions", func(t *testing.T) {
		// Create PIC with specific permissions
		picReq := &dto.CreateEventPICRequest{
			PersonID:    personID,
			Role:        "Limited PIC",
			IsPrimary:   false,
			StartDate:   "2024-11-01",
			CanEdit:     true,
			CanDelete:   false,
			CanAssignPIC: false,
		}

		pic, err := eventPICService.CreateEventPIC(testEvent.ID, picReq, createdBy)
		require.NoError(t, err)

		// Test edit permission
		hasEditPermission, err := eventPICService.ValidateEventPICPermissions(testEvent.ID, personID, "edit")
		require.NoError(t, err)
		assert.True(t, hasEditPermission)

		// Test delete permission (should be false)
		hasDeletePermission, err := eventPICService.ValidateEventPICPermissions(testEvent.ID, personID, "delete")
		require.NoError(t, err)
		assert.False(t, hasDeletePermission)

		// Test assign_pic permission (should be false)
		hasAssignPermission, err := eventPICService.ValidateEventPICPermissions(testEvent.ID, personID, "assign_pic")
		require.NoError(t, err)
		assert.False(t, hasAssignPermission)

		// Test non-existing person
		randomPersonID := uuid.New()
		hasPermission, err := eventPICService.ValidateEventPICPermissions(testEvent.ID, randomPersonID, "edit")
		require.NoError(t, err)
		assert.False(t, hasPermission)

		// Clean up
		eventPICService.DeleteEventPIC(pic.ID, createdBy, "Test cleanup")
	})
}

func TestEventWithPICCreation(t *testing.T) {
	db := SetUpDatabaseConnection()
	eventRepo := repository.NewEventRepository(db)
	eventPICRepo := repository.NewEventPICRepository(db)
	eventService := service.NewEventService(eventRepo, eventPICRepo)

	createdBy := uuid.New()
	personID1 := uuid.New()
	personID2 := uuid.New()

	t.Run("Create Event with PICs", func(t *testing.T) {
		eventReq := &dto.CreateEventRequest{
			Title:         "Event with PICs",
			Description:   "Test event created with PICs",
			EventDate:     "2024-12-25",
			StartTime:     "15:00",
			EndTime:       "17:00",
			EventLocation: "Test Location with PICs",
			Type:          "event",
			Timezone:      "Asia/Jakarta",
			IsPublic:      true,
			EventPICs: []dto.CreateEventPICRequest{
				{
					PersonID:    personID1,
					Role:        "Primary PIC",
					IsPrimary:   true,
					StartDate:   "2024-11-01",
					CanEdit:     true,
					CanDelete:   true,
					CanAssignPIC: true,
				},
				{
					PersonID:    personID2,
					Role:        "Technical Support",
					IsPrimary:   false,
					StartDate:   "2024-11-01",
					CanEdit:     false,
					CanDelete:   false,
					CanAssignPIC: false,
				},
			},
		}

		event, err := eventService.CreateEventWithPICs(eventReq, createdBy)
		require.NoError(t, err)
		require.NotNil(t, event)

		// Verify event has PICs
		assert.Len(t, event.EventPICs, 2)
		assert.NotNil(t, event.PrimaryPIC)
		assert.Equal(t, "Primary PIC", event.PrimaryPIC.Role)
		assert.Equal(t, personID1, event.PrimaryPIC.PersonID)

		// Clean up
		eventService.DeleteEvent(event.ID)
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}