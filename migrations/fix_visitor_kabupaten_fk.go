package migrations

import (
	"gorm.io/gorm"
)

// FixVisitorKabupatenFK fixes the foreign key constraint issue with visitors table
func FixVisitorKabupatenFK(db *gorm.DB) error {
	// Check if table exists
	if !db.Migrator().HasTable("visitors") {
		return nil
	}

	// Drop existing foreign key constraint if it exists
	db.Exec(`ALTER TABLE visitors DROP FOREIGN KEY IF EXISTS fk_visitors_kabupaten`)

	// Make sure kabupaten_id column is nullable
	db.Exec(`ALTER TABLE visitors MODIFY COLUMN kabupaten_id INT NULL`)

	// Clean up invalid foreign key references - set to NULL if kabupaten doesn't exist
	db.Exec(`
		UPDATE visitors
		SET kabupaten_id = NULL
		WHERE kabupaten_id IS NOT NULL
		  AND kabupaten_id NOT IN (SELECT id FROM kabupatens)
	`)

	// Recreate foreign key with proper constraint
	db.Exec(`
		ALTER TABLE visitors
		ADD CONSTRAINT fk_visitors_kabupaten
		FOREIGN KEY (kabupaten_id)
		REFERENCES kabupatens(id)
		ON UPDATE CASCADE
		ON DELETE SET NULL
	`)

	return nil
}
