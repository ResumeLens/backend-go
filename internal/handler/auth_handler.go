package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req services.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.authService.Signup(req)
	c.JSON(statusCode, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.authService.Login(req)
	c.JSON(statusCode, response)
}

func (h *AuthHandler) Invite(c *gin.Context) {
	inviterRole, _ := c.Get("role")
	inviterOrgID, _ := c.Get("organizationID")

	var req services.InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.authService.Invite(req, inviterRole.(string), inviterOrgID.(string))
	c.JSON(statusCode, response)
}

func (h *AuthHandler) ValidateInvite(c *gin.Context) {
	token := c.Query("token")
	response, statusCode := h.authService.ValidateInvite(token)
	c.JSON(statusCode, response)
}

func (h *AuthHandler) AcceptInvite(c *gin.Context) {
	var req services.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.authService.AcceptInvite(req)
	c.JSON(statusCode, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req services.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, statusCode := h.authService.RefreshToken(req)
	c.JSON(statusCode, response)
}
