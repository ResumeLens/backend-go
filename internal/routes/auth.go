package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
	"github.com/resumelens/authservice/internal/utils"
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
	if err := db.DB.Where("name = ?", req.OrganizationName).First(&org).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Organization already exists. Please contact your admin or use an invite."})
		return
	} else if err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking organization"})
		return
	}

	org = models.Organization{
		Name:      req.OrganizationName,
		CreatedBy: req.Email,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create User (Admin by default)
	user := models.User{
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		Role:           "admin",
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
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user_id":      user.ID,
		"role":         user.Role,
		"organization": user.OrganizationID,
	})
}
