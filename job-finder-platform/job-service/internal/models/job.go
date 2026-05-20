package models

import (
	"strings"
	"time"
)

type Job struct {
	ID               string    `db:"id"`
	EmployerID       string    `db:"employer_id"`
	Title            string    `db:"title"`
	Description      string    `db:"description"`
	Location         string    `db:"location"`
	JobType          string    `db:"job_type"`
	Category         string    `db:"category"`
	SalaryMin        float64   `db:"salary_min"`
	SalaryMax        float64   `db:"salary_max"`
	Currency         string    `db:"currency"`
	Requirements     []string
	RequirementsStr  string    `db:"requirements"`
	ExperienceLevel  string    `db:"experience_level"`
	Status           string    `db:"status"`
	ApplicationCount int       `db:"application_count"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (j *Job) ParseRequirements() {
	if j.RequirementsStr != "" {
		j.Requirements = strings.Split(j.RequirementsStr, "||")
	}
}

type Application struct {
	ID          string    `db:"id"`
	JobID       string    `db:"job_id"`
	ApplicantID string    `db:"applicant_id"`
	CoverLetter string    `db:"cover_letter"`
	ResumeURL   string    `db:"resume_url"`
	Status      string    `db:"status"`
	AppliedAt   time.Time `db:"applied_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateJobInput struct {
	EmployerID      string
	Title           string
	Description     string
	Location        string
	JobType         string
	Category        string
	SalaryMin       float64
	SalaryMax       float64
	Currency        string
	Requirements    []string
	ExperienceLevel string
}

type UpdateJobInput struct {
	Title       string
	Description string
	Location    string
	JobType     string
	Category    string
	SalaryMin   float64
	SalaryMax   float64
	Status      string
}

type SearchInput struct {
	Query           string
	Location        string
	Category        string
	JobType         string
	ExperienceLevel string
	SalaryMin       float64
	SalaryMax       float64
	Page            int
	Limit           int
}
