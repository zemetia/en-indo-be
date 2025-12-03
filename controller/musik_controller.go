package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type MusikController interface {
	GetPelayanMusik(ctx *gin.Context)
	GetPelayanMusikByID(ctx *gin.Context)
	AssignPelayanan(ctx *gin.Context)
	RemovePelayanan(ctx *gin.Context)
	ToggleActive(ctx *gin.Context)
	GetAvailablePelayanan(ctx *gin.Context)
}

type musikController struct {
	musikService service.MusikService
}

func NewMusikController(musikService service.MusikService) MusikController {
	return &musikController{
		musikService: musikService,
	}
}

// GetPelayanMusik handles GET /api/musik/pelayan
// Returns list of musicians based on PIC permissions
func (c *musikController) GetPelayanMusik(ctx *gin.Context) {
	// Get user_id from JWT middleware context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call service
	musicians, err := c.musikService.GetPelayanMusik(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, musicians)
}

// GetPelayanMusikByID handles GET /api/musik/pelayan/:person_id
// Returns detailed information about a specific musician
func (c *musikController) GetPelayanMusikByID(ctx *gin.Context) {
	// Get user_id from JWT middleware context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get person_id from URL param
	personIDStr := ctx.Param("person_id")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	// Call service
	musician, err := c.musikService.GetPelayanMusikByID(ctx.Request.Context(), userID, personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, musician)
}

// AssignPelayanan handles POST /api/musik/pelayan/:person_id/assign
// Assigns a pelayanan role to a musician
func (c *musikController) AssignPelayanan(ctx *gin.Context) {
	// Get user_id from JWT middleware context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get person_id from URL param
	personIDStr := ctx.Param("person_id")
	if personIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Person ID is required"})
		return
	}

	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	// Parse request body
	var req dto.AssignPelayananRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set person_id from URL param (overriding any value in body)
	req.PersonID = personID

	// Call service
	if err := c.musikService.AssignPelayananToMusician(ctx.Request.Context(), userID, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Pelayanan assigned successfully"})
}

// RemovePelayanan handles DELETE /api/musik/pelayan/assignment/:assignment_id
// Removes a pelayanan assignment from a musician
func (c *musikController) RemovePelayanan(ctx *gin.Context) {
	// Get user_id from JWT middleware context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get assignment_id from URL param
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Call service
	if err := c.musikService.RemovePelayananFromMusician(ctx.Request.Context(), userID, assignmentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Pelayanan removed successfully"})
}

// ToggleActive handles PUT /api/musik/pelayan/assignment/:assignment_id/toggle
// Toggles the active status of a pelayanan assignment
func (c *musikController) ToggleActive(ctx *gin.Context) {
	// Get user_id from JWT middleware context
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get assignment_id from URL param
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Parse request body
	var req dto.ToggleActiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	if err := c.musikService.ToggleActiveStatus(ctx.Request.Context(), userID, assignmentID, req.IsActive); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Active status updated successfully"})
}

// GetAvailablePelayanan handles GET /api/musik/pelayanan-roles
// Returns all available pelayanan roles in the music department
func (c *musikController) GetAvailablePelayanan(ctx *gin.Context) {
	// Call service (no authentication required for getting available roles)
	pelayananRoles, err := c.musikService.GetAvailableMusikPelayanan(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pelayananRoles)
}
