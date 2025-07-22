package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/config"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
)

type JobHostingService struct {
	config *config.Config
}

func NewJobHostingService(cfg *config.Config) *JobHostingService {
	return &JobHostingService{config: cfg}
}

type CreateJobRequest struct {
	Title           string   `json:"title" binding:"required"`
	OrganizationID  string   `json:"organization_id" binding:"required"`
	CreatedByID     string   `json:"created_by_id" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Location        []string `json:"location" binding:"required"`
	ExperienceLevel string   `json:"experience_level" binding:"required"`
	SkillsRequired  []string `json:"skills_required" binding:"required"`
	EmploymentType  []string `json:"employment_type" binding:"required"`
	SalaryRange     []string `json:"salary_range" binding:"required"`
	IsActive        bool     `json:"is_active" binding:"required"`
}

func (s *JobHostingService) CreateJob(req CreateJobRequest) (gin.H, int) {
	var org models.Organization
	var job models.Job

	if err := db.DB.Where("id = ?", req.OrganizationID).First(&org).Error; err != nil {
		return gin.H{"error": "User not found"}, http.StatusNotFound
	}

	job.Title = req.Title
	job.OrganizationID = req.OrganizationID
	job.CreatedByID = req.CreatedByID
	job.Description = req.Description
	job.Location = req.Location
	job.ExperienceLevel = req.ExperienceLevel
	job.SkillsRequired = req.SkillsRequired
	job.EmploymentType = req.EmploymentType
	job.SalaryRange = req.SalaryRange
	job.IsActive = req.IsActive
	job.ApplicationCount = 0
	job.CreatedAt = time.Now()

	if err := db.DB.Create(&job).Error; err != nil {
		return gin.H{"error": err.Error()}, http.StatusInternalServerError
	}

	job.PublicLink = fmt.Sprintf("https://resumelens.com/job/%s/%s", job.OrganizationID, job.ID)
	job.ShortLink = fmt.Sprintf("https://resumelens.com/job/%s", job.ID)

	if err := db.DB.Model(&job).Updates(map[string]interface{}{
		"public_link": job.PublicLink,
		"short_link":  job.ShortLink,
	}).Error; err != nil {
		return gin.H{"error": err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"message": "Job created successfully"}, http.StatusOK
}

func (s *JobHostingService) GetJob(id string) (gin.H, int) {
	var job models.Job
	if err := db.DB.Where("id = ?", id).First(&job).Error; err != nil {
		return gin.H{"error": "Job not found"}, http.StatusNotFound
	}

	return gin.H{"job": job}, http.StatusOK
}
