package services

import (
	"fmt"
	"net/http"
	"strings"
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
	Title           string   `json:"title" binding:"required,min=1,max=200"`
	Description     string   `json:"description" binding:"required,min=10,max=5000"`
	Location        []string `json:"location" binding:"required,min=1,dive,min=1,max=100"`
	ExperienceLevel string   `json:"experience_level" binding:"required"`
	SkillsRequired  []string `json:"skills_required" binding:"required,min=1,dive,min=1,max=50"`
	EmploymentType  []string `json:"employment_type" binding:"required,min=1,dive,min=1,max=50"`
	SalaryRange     []string `json:"salary_range" binding:"required,min=1,dive,min=1,max=50"`
	IsActive        bool     `json:"is_active"`
}

type UpdateJobRequest struct {
	Title           *string   `json:"title" binding:"omitempty,min=1,max=200"`
	Description     *string   `json:"description" binding:"omitempty,min=10,max=5000"`
	Location        *[]string `json:"location" binding:"omitempty,min=1,dive,min=1,max=100"`
	ExperienceLevel *string   `json:"experience_level" binding:"omitempty"`
	SkillsRequired  *[]string `json:"skills_required" binding:"omitempty,min=1,dive,min=1,max=50"`
	EmploymentType  *[]string `json:"employment_type" binding:"omitempty,min=1,dive,min=1,max=50"`
	SalaryRange     *[]string `json:"salary_range" binding:"omitempty,min=1,dive,min=1,max=50"`
	IsActive        *bool     `json:"is_active"`
}

type JobListResponse struct {
	Jobs       []models.Job `json:"jobs"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

func (s *JobHostingService) CreateJob(req CreateJobRequest, orgID, userID string) (gin.H, int) {
	var org models.Organization
	if err := db.DB.Where("id = ?", orgID).First(&org).Error; err != nil {
		return gin.H{"error": "Organization not found"}, http.StatusNotFound
	}

	job := models.Job{
		Title:            strings.TrimSpace(req.Title),
		OrganizationID:   orgID,
		CreatedByID:      userID,
		Description:      strings.TrimSpace(req.Description),
		Location:         sanitizeStringArray(req.Location),
		ExperienceLevel:  strings.TrimSpace(req.ExperienceLevel),
		SkillsRequired:   sanitizeStringArray(req.SkillsRequired),
		EmploymentType:   sanitizeStringArray(req.EmploymentType),
		SalaryRange:      sanitizeStringArray(req.SalaryRange),
		IsActive:         req.IsActive,
		ApplicationCount: 0,
		CreatedAt:        time.Now(),
	}

	if err := db.DB.Create(&job).Error; err != nil {
		return gin.H{"error": "Failed to create job"}, http.StatusInternalServerError
	}

	job.PublicLink = fmt.Sprintf("https://resumelens.com/job/%s/%s", job.OrganizationID, job.ID)
	job.ShortLink = fmt.Sprintf("https://resumelens.com/job/%s", job.ID)

	if err := db.DB.Model(&job).Updates(map[string]interface{}{
		"public_link": job.PublicLink,
		"short_link":  job.ShortLink,
	}).Error; err != nil {
		return gin.H{"error": "Failed to update job links"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Job created successfully", "job_id": job.ID}, http.StatusCreated
}

func (s *JobHostingService) GetJob(id string) (gin.H, int) {
	if !isValidUUID(id) {
		return gin.H{"error": "Invalid job ID format"}, http.StatusBadRequest
	}

	var job models.Job
	if err := db.DB.Where("id = ? AND is_active = ?", id, true).First(&job).Error; err != nil {
		return gin.H{"error": "Job not found"}, http.StatusNotFound
	}

	return gin.H{"job": job}, http.StatusOK
}

func (s *JobHostingService) GetJobsForOrganization(orgID string, page, pageSize int) (gin.H, int) {
	if !isValidUUID(orgID) {
		return gin.H{"error": "Invalid organization ID format"}, http.StatusBadRequest
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var jobs []models.Job
	var total int64

	if err := db.DB.Model(&models.Job{}).Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return gin.H{"error": "Failed to count jobs"}, http.StatusInternalServerError
	}

	if err := db.DB.Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&jobs).Error; err != nil {
		return gin.H{"error": "Failed to fetch jobs"}, http.StatusInternalServerError
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := JobListResponse{
		Jobs:       jobs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return gin.H{"data": response}, http.StatusOK
}

func (s *JobHostingService) UpdateJob(jobID string, updates map[string]interface{}) (gin.H, int) {
	if !isValidUUID(jobID) {
		return gin.H{"error": "Invalid job ID format"}, http.StatusBadRequest
	}

	allowedFields := map[string]bool{
		"title": true, "description": true, "location": true,
		"experience_level": true, "skills_required": true,
		"employment_type": true, "salary_range": true, "is_active": true,
	}

	sanitizedUpdates := make(map[string]interface{})
	for field, value := range updates {
		if allowedFields[field] {
			switch field {
			case "title", "description", "experience_level":
				if str, ok := value.(string); ok {
					sanitizedUpdates[field] = strings.TrimSpace(str)
				}
			case "location", "skills_required", "employment_type", "salary_range":
				if arr, ok := value.([]string); ok {
					sanitizedUpdates[field] = sanitizeStringArray(arr)
				}
			case "is_active":
				if _, ok := value.(bool); ok {
					sanitizedUpdates[field] = value
				}
			}
		}
	}

	if len(sanitizedUpdates) == 0 {
		return gin.H{"error": "No valid fields to update"}, http.StatusBadRequest
	}

	if err := db.DB.Model(&models.Job{}).Where("id = ?", jobID).Updates(sanitizedUpdates).Error; err != nil {
		return gin.H{"error": "Failed to update job"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Job updated successfully"}, http.StatusOK
}

func (s *JobHostingService) DeleteJob(jobID string) (gin.H, int) {
	if !isValidUUID(jobID) {
		return gin.H{"error": "Invalid job ID format"}, http.StatusBadRequest
	}

	if err := db.DB.Model(&models.Job{}).Where("id = ?", jobID).Update("is_active", false).Error; err != nil {
		return gin.H{"error": "Failed to delete job"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Job deleted successfully"}, http.StatusOK
}

func (s *JobHostingService) GetJobApplications(jobID string, page, pageSize int) (gin.H, int) {
	if !isValidUUID(jobID) {
		return gin.H{"error": "Invalid job ID format"}, http.StatusBadRequest
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var applications []models.JobApplication
	var total int64

	if err := db.DB.Model(&models.JobApplication{}).Where("job_id = ?", jobID).Count(&total).Error; err != nil {
		return gin.H{"error": "Failed to count applications"}, http.StatusInternalServerError
	}

	if err := db.DB.Preload("Candidate").
		Where("job_id = ?", jobID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications).Error; err != nil {
		return gin.H{"error": "Failed to fetch applications"}, http.StatusInternalServerError
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return gin.H{
		"applications": applications,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_pages":  totalPages,
	}, http.StatusOK
}

func sanitizeStringArray(arr []string) []string {
	var sanitized []string
	for _, s := range arr {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			sanitized = append(sanitized, trimmed)
		}
	}
	return sanitized
}

func isValidUUID(id string) bool {
	if len(id) != 36 {
		return false
	}
	parts := strings.Split(id, "-")
	return len(parts) == 5 && len(parts[0]) == 8 && len(parts[1]) == 4 && len(parts[2]) == 4 && len(parts[3]) == 4 && len(parts[4]) == 12
}
