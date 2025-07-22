package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
	"github.com/resumelens/authservice/internal/utils"
	"github.com/spf13/viper"
)

// SignupRequest represents incoming signup data
type SignupRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	OrganizationName string `json:"organization_name" binding:"required"`
}

// SignupHandler handles new user registration
func SignupHandler(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Check if organization exists
	var org models.Organization
	orgExists := false
	if err := db.DB.Where("name = ?", req.OrganizationName).First(&org).Error; err == nil {
		orgExists = true
	} else if err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking organization"})
		return
	}

	if orgExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Organization already exists. Please contact your admin or use an invite."})
		return
	}

	// Organization does not exist, create it
	org = models.Organization{
		Name:      req.OrganizationName,
		CreatedBy: req.Email,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Create admin role for the organization
	adminRole := models.Role{
		Name:                "admin",
		OrganizationID:      org.ID,
		CreatedAt:           time.Now(),
		HomePermission:      true,
		CreateJobPermission: true,
		ViewJobPermission:   true,
		IAMPermission:       true,
	}
	if err := db.DB.Create(&adminRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin role"})
		return
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create User with admin role
	user := models.User{
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		RoleID:         adminRole.ID,
		OrganizationID: org.ID,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Signup successful",
		"user_id":         user.ID,
		"organization_id": org.ID,
		"role_id":         adminRole.ID,
	})
}

// LoginRequest represents incoming login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler handles user authentication
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.RoleID, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user_id":      user.ID,
		"role":         user.RoleID,
		"organization": user.OrganizationID,
	})
}

// InviteRequest represents the request body to invite a member
type InviteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

// InviteHandler allows an admin to invite another member
func InviteHandler(c *gin.Context) {
	// Extract the inviter's details from context
	inviterRole, exists := c.Get("role")
	if !exists || inviterRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can invite members"})
		return
	}

	inviterOrgID, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization ID not found in context"})
		return
	}

	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Check if user with this email already exists in the system
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Generate unique token for invite
	inviteToken := utils.GenerateRandomToken(32)

	// In InviteHandler, set RoleID to an empty string or add a TODO for role lookup
	invite := models.Invite{
		Email:          req.Email,
		OrganizationID: inviterOrgID.(string),
		RoleID:         "", // TODO: Lookup role ID by name and organization
		Token:          inviteToken,
		Expiry:         time.Now().Add(48 * time.Hour), // 2-day expiry
		IsAccepted:     false,
		CreatedAt:      time.Now(),
	}

	if err := db.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invite"})
		return
	}

	// In production: Send invite via email (TODO: SMTP)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Invite created successfully",
		"invite_token": invite.Token,
	})

	err := utils.SendInviteEmail(req.Email, invite.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invite email"})
		return
	}

}

// ValidateInviteHandler checks if the provided invite token is valid and not expired
func ValidateInviteHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invite token is required"})
		return
	}

	var invite models.Invite
	if err := db.DB.Where("token = ? AND is_accepted = false", token).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or already used invite token"})
		return
	}

	if time.Now().After(invite.Expiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invite token has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":           true,
		"email":           invite.Email,
		"organization_id": invite.OrganizationID,
		"role_id":         invite.RoleID, // Changed from invite.Role to invite.RoleID
	})
}

// AcceptInviteRequest represents the request body to accept an invite
type AcceptInviteRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// AcceptInviteHandler allows a user to accept an invite
func AcceptInviteHandler(c *gin.Context) {
	var req AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invite models.Invite
	if err := db.DB.Where("token = ? AND is_accepted = false", req.Token).First(&invite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired invite token"})
		return
	}

	if invite.Expiry.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invite has expired"})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", invite.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:          invite.Email,
		PasswordHash:   hashedPassword,
		RoleID:         invite.RoleID, // Changed from invite.Role to invite.RoleID
		OrganizationID: invite.OrganizationID,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	invite.IsAccepted = true
	db.DB.Save(&invite)
	c.JSON(http.StatusCreated, gin.H{
		"message":         "Account created successfully via invite",
		"user_id":         user.ID,
		"organization_id": user.OrganizationID,
	})
}

// RefreshTokenRequest represents the incoming refresh token data
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenHandler issues a new access token
func RefreshTokenHandler(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Generate a new access token
	newAccessToken, err := utils.GenerateJWT(claims.UserID, claims.Email, claims.Role, claims.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"expires_in":   viper.GetInt("JWT_EXPIRY"),
	})
}
