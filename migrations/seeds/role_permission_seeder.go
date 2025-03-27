package seeds

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func RolePermissionSeeder(db *gorm.DB) error {
	// Create Permissions
	permissions := []entity.Permission{
		{
			Name:        "manage_users",
			Description: "Can manage users",
		},
		{
			Name:        "manage_roles",
			Description: "Can manage roles and permissions",
		},
		{
			Name:        "manage_content",
			Description: "Can manage website content",
		},
		{
			Name:        "manage_media",
			Description: "Can manage media files",
		},
		{
			Name:        "manage_events",
			Description: "Can manage church events",
		},
		{
			Name:        "manage_sermons",
			Description: "Can manage sermons",
		},
		{
			Name:        "manage_donations",
			Description: "Can manage donations",
		},
		{
			Name:        "view_reports",
			Description: "Can view reports",
		},
	}

	for _, permission := range permissions {
		if err := db.Create(&permission).Error; err != nil {
			return err
		}
	}

	// Create Roles
	roles := []entity.Role{
		{
			Name:        "Admin",
			Description: "System Administrator",
			Permissions: permissions[:], // All permissions
		},
		{
			Name:        "Pastor",
			Description: "Church Pastor",
			Permissions: permissions[2:7], // Most permissions except user/role management
		},
		{
			Name:        "Leader",
			Description: "Church Leader",
			Permissions: permissions[2:5], // Content, media, and events management
		},
		{
			Name:        "LG Leader",
			Description: "Life Group Leader",
			Permissions: permissions[2:4], // Content and media management
		},
		{
			Name:        "Pemusik",
			Description: "Church Musician",
			Permissions: permissions[3:4], // Media management
		},
		{
			Name:        "Multimedia",
			Description: "Multimedia Team",
			Permissions: permissions[3:4], // Media management
		},
		{
			Name:        "Dokumentasi",
			Description: "Documentation Team",
			Permissions: permissions[3:4], // Media management
		},
	}

	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
