package seeds

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func SeedEventPICRoles(db *gorm.DB) error {
	// Check if roles already exist
	var count int64
	if err := db.Model(&entity.EventPICRole{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		// Roles already seeded
		return nil
	}

	// Default Event PIC roles
	roles := []entity.EventPICRole{
		{
			ID:                  uuid.New(),
			Name:                entity.EventPICRolePrimary,
			Description:         "Primary Person in Charge - Overall responsibility for the event",
			DefaultCanEdit:      true,
			DefaultCanDelete:    true,
			DefaultCanAssignPIC: true,
			IsActive:            true,
		},
		{
			ID:                  uuid.New(),
			Name:                entity.EventPICRoleSecondary,
			Description:         "Secondary Person in Charge - Assists the primary PIC",
			DefaultCanEdit:      true,
			DefaultCanDelete:    false,
			DefaultCanAssignPIC: false,
			IsActive:            true,
		},
		{
			ID:                  uuid.New(),
			Name:                entity.EventPICRoleTechnical,
			Description:         "Technical Person in Charge - Handles technical aspects of the event",
			DefaultCanEdit:      true,
			DefaultCanDelete:    false,
			DefaultCanAssignPIC: false,
			IsActive:            true,
		},
		{
			ID:                  uuid.New(),
			Name:                entity.EventPICRoleLogistics,
			Description:         "Logistics Person in Charge - Handles event logistics and supplies",
			DefaultCanEdit:      true,
			DefaultCanDelete:    false,
			DefaultCanAssignPIC: false,
			IsActive:            true,
		},
		{
			ID:                  uuid.New(),
			Name:                entity.EventPICRoleRegistration,
			Description:         "Registration Person in Charge - Handles participant registration",
			DefaultCanEdit:      false,
			DefaultCanDelete:    false,
			DefaultCanAssignPIC: false,
			IsActive:            true,
		},
		{
			ID:              uuid.New(),
			Name:            "Worship Leader",
			Description:     "Worship Leader for the event - Leads worship and praise",
			DefaultCanEdit:  false,
			DefaultCanDelete: false,
			DefaultCanAssignPIC: false,
			IsActive:        true,
		},
		{
			ID:              uuid.New(),
			Name:            "Prayer Coordinator",
			Description:     "Prayer Coordinator - Organizes and leads prayer activities",
			DefaultCanEdit:  false,
			DefaultCanDelete: false,
			DefaultCanAssignPIC: false,
			IsActive:        true,
		},
		{
			ID:              uuid.New(),
			Name:            "Youth Coordinator",
			Description:     "Youth Coordinator - Handles youth-specific activities and programs",
			DefaultCanEdit:  false,
			DefaultCanDelete: false,
			DefaultCanAssignPIC: false,
			IsActive:        true,
		},
		{
			ID:              uuid.New(),
			Name:            "Children's Ministry Leader",
			Description:     "Children's Ministry Leader - Oversees children's programs during the event",
			DefaultCanEdit:  false,
			DefaultCanDelete: false,
			DefaultCanAssignPIC: false,
			IsActive:        true,
		},
		{
			ID:              uuid.New(),
			Name:            "Security Coordinator",
			Description:     "Security Coordinator - Ensures event safety and security",
			DefaultCanEdit:  false,
			DefaultCanDelete: false,
			DefaultCanAssignPIC: false,
			IsActive:        true,
		},
	}

	// Create all roles
	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}