package migrations

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

// MigratePicField handles the migration of PIC functionality from PersonPelayananGereja to Pelayanan
func MigratePicField(db *gorm.DB) error {
	// Step 1: Create PIC pelayanan for all existing departments
	if err := createPicPelayananForDepartments(db); err != nil {
		return fmt.Errorf("failed to create PIC pelayanan: %v", err)
	}

	// Step 2: Migrate existing PIC assignments to the new PIC pelayanan
	if err := migratePicAssignments(db); err != nil {
		return fmt.Errorf("failed to migrate PIC assignments: %v", err)
	}

	// Step 3: Drop is_pic column from person_pelayanan_gereja table
	if err := dropIsPicFromPersonPelayanan(db); err != nil {
		return fmt.Errorf("failed to drop is_pic column: %v", err)
	}

	return nil
}

// createPicPelayananForDepartments creates PIC pelayanan for all existing departments
func createPicPelayananForDepartments(db *gorm.DB) error {
	var departments []entity.Department
	if err := db.Find(&departments).Error; err != nil {
		return err
	}

	for _, dept := range departments {
		// Check if PIC pelayanan already exists for this department
		var existingPic entity.Pelayanan
		err := db.Where("department_id = ? AND is_pic = ?", dept.ID, true).First(&existingPic).Error
		
		if err == gorm.ErrRecordNotFound {
			// Create new PIC pelayanan
			picPelayanan := entity.Pelayanan{
				ID:           uuid.New(),
				Pelayanan:    fmt.Sprintf("PIC %s", dept.Name),
				Description:  fmt.Sprintf("Person in Charge untuk departemen %s", dept.Name),
				DepartmentID: dept.ID,
				IsPic:        true,
			}
			
			if err := db.Create(&picPelayanan).Error; err != nil {
				return fmt.Errorf("failed to create PIC pelayanan for department %s: %v", dept.Name, err)
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

// migratePicAssignments moves existing PIC assignments to the new PIC pelayanan
func migratePicAssignments(db *gorm.DB) error {
	// Find all existing PIC assignments
	var picAssignments []PersonPelayananGerejaForMigration
	if err := db.Where("is_pic = ?", true).Find(&picAssignments).Error; err != nil {
		return err
	}

	for _, assignment := range picAssignments {
		// Find the corresponding pelayanan to get the department
		var pelayanan entity.Pelayanan
		if err := db.First(&pelayanan, assignment.PelayananID).Error; err != nil {
			continue // Skip if pelayanan not found
		}

		// Find the PIC pelayanan for this department
		var picPelayanan entity.Pelayanan
		if err := db.Where("department_id = ? AND is_pic = ?", pelayanan.DepartmentID, true).First(&picPelayanan).Error; err != nil {
			continue // Skip if PIC pelayanan not found
		}

		// Check if assignment to PIC pelayanan already exists
		var existingAssignment entity.PersonPelayananGereja
		err := db.Where("person_id = ? AND pelayanan_id = ? AND church_id = ?", 
			assignment.PersonID, picPelayanan.ID, assignment.ChurchID).First(&existingAssignment).Error

		if err == gorm.ErrRecordNotFound {
			// Create new assignment to PIC pelayanan
			newAssignment := entity.PersonPelayananGereja{
				ID:          uuid.New(),
				PersonID:    assignment.PersonID,
				PelayananID: picPelayanan.ID,
				ChurchID:    assignment.ChurchID,
			}
			
			if err := db.Create(&newAssignment).Error; err != nil {
				return fmt.Errorf("failed to create PIC assignment: %v", err)
			}
		}
	}

	return nil
}

// dropIsPicFromPersonPelayanan drops the is_pic column from person_pelayanan_gereja table
func dropIsPicFromPersonPelayanan(db *gorm.DB) error {
	// Check if column exists before trying to drop it
	if db.Migrator().HasColumn(&PersonPelayananGerejaForMigration{}, "is_pic") {
		if err := db.Migrator().DropColumn(&PersonPelayananGerejaForMigration{}, "is_pic"); err != nil {
			return err
		}
	}
	return nil
}

// PersonPelayananGerejaForMigration is a temporary struct for migration purposes
type PersonPelayananGerejaForMigration struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	PersonID    uuid.UUID `gorm:"type:char(36);not null"`
	PelayananID uuid.UUID `gorm:"type:char(36);not null"`
	ChurchID    uuid.UUID `gorm:"type:char(36);not null"`
	IsPic       bool      `gorm:"column:is_pic"`
}

func (PersonPelayananGerejaForMigration) TableName() string {
	return "person_pelayanan_gerejas"
}