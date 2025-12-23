package dto

// MusicianResponse represents a musician with their pelayanan assignments
// This matches the frontend Musician interface
type MusicianResponse struct {
	ID          string   `json:"id"`
	Nama        string   `json:"nama"`
	Email       string   `json:"email"`
	Telepon     string   `json:"telepon"`
	Instruments []string `json:"instruments"` // Array of pelayanan names (Gitaris, Drummer, etc.)
	Status      string   `json:"status"`      // "active" or "inactive"
	Avatar      string   `json:"avatar"`
}

// ToggleActiveRequest is used to toggle the active status of a pelayanan assignment
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// PelayananRoleResponse represents a pelayanan role available in the music department
type PelayananRoleResponse struct {
	ID          string `json:"id"`
	Pelayanan   string `json:"pelayanan"`    // e.g., "Gitaris", "Drummer"
	Description string `json:"description"`
	IsPic       bool   `json:"is_pic"`
}

// MusicianDetailResponse provides detailed information about a musician
type MusicianDetailResponse struct {
	ID         string                             `json:"id"`
	Nama       string                             `json:"nama"`
	Email      string                             `json:"email"`
	Telepon    string                             `json:"telepon"`
	Avatar     string                             `json:"avatar"`
	ChurchID   string                             `json:"church_id"`
	ChurchName string                             `json:"church_name"`
	Pelayanan  []MusicianPelayananDetailResponse `json:"pelayanan"` // All pelayanan assignments
}

// MusicianPelayananDetailResponse represents a single pelayanan assignment for a musician
type MusicianPelayananDetailResponse struct {
	AssignmentID string `json:"assignment_id"` // PersonPelayananGereja ID
	PelayananID  string `json:"pelayanan_id"`
	Pelayanan    string `json:"pelayanan"`    // Name of pelayanan
	ChurchID     string `json:"church_id"`
	ChurchName   string `json:"church_name"`
	IsActive     bool   `json:"is_active"`
	IsPic        bool   `json:"is_pic"`
}

// AvailablePersonResponse represents a person available to be added to music ministry
type AvailablePersonResponse struct {
	ID         string `json:"id"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	Telepon    string `json:"telepon"`
	ChurchID   string `json:"church_id"`
	ChurchName string `json:"church_name"`
	Avatar     string `json:"avatar"`
}
