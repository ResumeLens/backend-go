package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
	"github.com/resumelens/authservice/internal/utils"
)

type AuthService struct {
	config *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{config: cfg}
}

type SignupRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	OrganizationName string `json:"organization_name" binding:"required"`
}

func (s *AuthService) Signup(req SignupRequest) (gin.H, int) {
	// Check if user exists first
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return gin.H{"error": "Email already registered"}, http.StatusConflict
	}

	// Check if org exists first
	var existingOrg models.Organization
	if err := db.DB.Where("name = ?", req.OrganizationName).First(&existingOrg).Error; err == nil {
		return gin.H{"error": "Organization already exists. Please contact your admin or use an invite."}, http.StatusConflict
	} else if err.Error() != "record not found" {
		return gin.H{"error": "Database error while checking organization"}, http.StatusInternalServerError
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return gin.H{"error": "Failed to hash password"}, http.StatusInternalServerError
	}

	org := models.Organization{
		Name:        req.OrganizationName,
		CreatedByID: nil, // will update after user is created
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&org).Error; err != nil {
		return gin.H{"error": "Failed to create organization"}, http.StatusInternalServerError
	}

	// Create admin role for this org
	adminRole := models.Role{
		Name:                "admin",
		OrganizationID:      org.ID,
		HomePermission:      true,
		CreateJobPermission: true,
		ViewJobPermission:   true,
		IamPermission:       true,
		CreatedAt:           time.Now(),
	}
	if err := db.DB.Create(&adminRole).Error; err != nil {
		return gin.H{"error": "Failed to create admin role"}, http.StatusInternalServerError
	}

	// Now create the user and assign the admin role
	user := models.User{
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		RoleID:         adminRole.ID,
		OrganizationID: org.ID,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return gin.H{"error": "Failed to create user"}, http.StatusInternalServerError
	}

	// Update org with created_by
	if err := db.DB.Model(&org).Update("created_by", user.ID).Error; err != nil {
		return gin.H{"error": "Failed to update organization with creator"}, http.StatusInternalServerError
	}

	return gin.H{
		"message":         "Signup successful",
		"user_id":         user.ID,
		"organization_id": org.ID,
		"role_id":         adminRole.ID,
	}, http.StatusCreated
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Login(req LoginRequest) (gin.H, int) {
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return gin.H{"error": "Invalid email or password"}, http.StatusUnauthorized
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return gin.H{"error": "Invalid email or password"}, http.StatusUnauthorized
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.RoleID, user.OrganizationID)
	if err != nil {
		return gin.H{"error": "Failed to generate token"}, http.StatusInternalServerError
	}

	return gin.H{
		"access_token": token,
		"user_id":      user.ID,
		"role":         user.RoleID,
		"organization": user.OrganizationID,
	}, http.StatusOK
}

type InviteRequest struct {
	Email  string `json:"email" binding:"required,email"`
	RoleID string `json:"role_id" binding:"required"`
}

func (s *AuthService) Invite(req InviteRequest, inviterRole, inviterOrgID string) (gin.H, int) {
	if inviterRole != "admin" {
		return gin.H{"error": "Only admins can invite members"}, http.StatusForbidden
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return gin.H{"error": "User with this email already exists"}, http.StatusConflict
	}

	inviteToken := utils.GenerateRandomToken(32)

	invite := models.Invite{
		Email:          req.Email,
		OrganizationID: inviterOrgID,
		RoleID:         req.RoleID,
		Token:          inviteToken,
		Expiry:         time.Now().Add(48 * time.Hour),
		IsAccepted:     false,
		CreatedAt:      time.Now(),
	}

	if err := db.DB.Create(&invite).Error; err != nil {
		return gin.H{"error": "Failed to create invite"}, http.StatusInternalServerError
	}

	err := utils.SendInviteEmail(req.Email, invite.Token, s.config)
	if err != nil {
		return gin.H{"error": "Failed to send invite email"}, http.StatusInternalServerError
	}

	return gin.H{
		"message":      "Invite created successfully",
		"invite_token": invite.Token,
	}, http.StatusOK
}

func (s *AuthService) ValidateInvite(token string) (gin.H, int) {
	if token == "" {
		return gin.H{"error": "Invite token is required"}, http.StatusBadRequest
	}

	var invite models.Invite
	if err := db.DB.Where("token = ? AND is_accepted = false", token).First(&invite).Error; err != nil {
		return gin.H{"error": "Invalid or already used invite token"}, http.StatusNotFound
	}

	if time.Now().After(invite.Expiry) {
		return gin.H{"error": "Invite token has expired"}, http.StatusBadRequest
	}

	return gin.H{
		"valid":           true,
		"email":           invite.Email,
		"organization_id": invite.OrganizationID,
		"role_id":         invite.RoleID,
	}, http.StatusOK
}

type AcceptInviteRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *AuthService) AcceptInvite(req AcceptInviteRequest) (gin.H, int) {
	var invite models.Invite
	if err := db.DB.Where("token = ? AND is_accepted = false", req.Token).First(&invite).Error; err != nil {
		return gin.H{"error": "Invalid or expired invite token"}, http.StatusNotFound
	}

	if invite.Expiry.Before(time.Now()) {
		return gin.H{"error": "Invite has expired"}, http.StatusUnauthorized
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", invite.Email).First(&existingUser).Error; err == nil {
		return gin.H{"error": "User with this email already exists"}, http.StatusConflict
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return gin.H{"error": "Failed to hash password"}, http.StatusInternalServerError
	}

	user := models.User{
		Email:          invite.Email,
		PasswordHash:   hashedPassword,
		RoleID:         invite.RoleID,
		OrganizationID: invite.OrganizationID,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return gin.H{"error": "Failed to create user"}, http.StatusInternalServerError
	}

	invite.IsAccepted = true
	db.DB.Save(&invite)

	return gin.H{
		"message":         "Account created successfully via invite",
		"user_id":         user.ID,
		"organization_id": user.OrganizationID,
	}, http.StatusCreated
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *AuthService) RefreshToken(req RefreshTokenRequest) (gin.H, int) {
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return gin.H{"error": "Invalid or expired refresh token"}, http.StatusUnauthorized
	}

	newAccessToken, err := utils.GenerateJWT(claims.UserID, claims.Email, claims.Role, claims.OrganizationID)
	if err != nil {
		return gin.H{"error": "Failed to generate new access token"}, http.StatusInternalServerError
	}

	return gin.H{
		"access_token": newAccessToken,
		"expires_in":   s.config.JWTExpiry,
	}, http.StatusOK
}
