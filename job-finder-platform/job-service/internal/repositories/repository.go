package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"jobfinder/job-service/internal/models"
)

// ── Job Repository ────────────────────────────────────────────────────────

type JobRepository interface {
	Create(ctx context.Context, in *models.CreateJobInput) (*models.Job, error)
	GetByID(ctx context.Context, id string) (*models.Job, error)
	Update(ctx context.Context, id string, in *models.UpdateJobInput) (*models.Job, error)
	Delete(ctx context.Context, id, employerID string) error
	List(ctx context.Context, page, limit int, status string) ([]*models.Job, int, error)
	Search(ctx context.Context, in *models.SearchInput) ([]*models.Job, int, error)
	GetByEmployer(ctx context.Context, employerID string, page, limit int) ([]*models.Job, int, error)
	IncrementApplicationCount(ctx context.Context, jobID string) error
}

type jobRepo struct{ db *pgxpool.Pool }

func NewJobRepository(db *pgxpool.Pool) JobRepository { return &jobRepo{db: db} }

const jobCols = `id,employer_id,title,description,location,job_type,category,salary_min,salary_max,currency,requirements,experience_level,status,application_count,created_at,updated_at`

func scanJob(row interface {
	Scan(...interface{}) error
}) (*models.Job, error) {
	j := &models.Job{}
	err := row.Scan(&j.ID, &j.EmployerID, &j.Title, &j.Description, &j.Location, &j.JobType,
		&j.Category, &j.SalaryMin, &j.SalaryMax, &j.Currency, &j.RequirementsStr,
		&j.ExperienceLevel, &j.Status, &j.ApplicationCount, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	j.ParseRequirements()
	return j, nil
}

func (r *jobRepo) Create(ctx context.Context, in *models.CreateJobInput) (*models.Job, error) {
	reqs := strings.Join(in.Requirements, "||")
	if in.Currency == "" { in.Currency = "USD" }
	if in.JobType == "" { in.JobType = "full-time" }
	if in.ExperienceLevel == "" { in.ExperienceLevel = "mid" }
	now := time.Now()
	row := r.db.QueryRow(ctx, `
		INSERT INTO jobs (id,employer_id,title,description,location,job_type,category,salary_min,salary_max,currency,requirements,experience_level,status,application_count,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'active',0,$13,$14)
		RETURNING `+jobCols,
		uuid.New().String(), in.EmployerID, in.Title, in.Description, in.Location, in.JobType,
		in.Category, in.SalaryMin, in.SalaryMax, in.Currency, reqs, in.ExperienceLevel, now, now)
	return scanJob(row)
}

func (r *jobRepo) GetByID(ctx context.Context, id string) (*models.Job, error) {
	row := r.db.QueryRow(ctx, `SELECT `+jobCols+` FROM jobs WHERE id=$1`, id)
	j, err := scanJob(row)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return j, nil
}

func (r *jobRepo) Update(ctx context.Context, id string, in *models.UpdateJobInput) (*models.Job, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE jobs SET title=$1,description=$2,location=$3,job_type=$4,category=$5,
		salary_min=$6,salary_max=$7,status=$8,updated_at=$9 WHERE id=$10
		RETURNING `+jobCols,
		in.Title, in.Description, in.Location, in.JobType, in.Category,
		in.SalaryMin, in.SalaryMax, in.Status, time.Now(), id)
	return scanJob(row)
}

func (r *jobRepo) Delete(ctx context.Context, id, employerID string) error {
	res, err := r.db.Exec(ctx, `DELETE FROM jobs WHERE id=$1 AND employer_id=$2`, id, employerID)
	if err != nil { return err }
	if res.RowsAffected() == 0 { return fmt.Errorf("job not found or unauthorized") }
	return nil
}

func (r *jobRepo) List(ctx context.Context, page, limit int, status string) ([]*models.Job, int, error) {
	var total int
	if status != "" {
		r.db.QueryRow(ctx, `SELECT COUNT(*) FROM jobs WHERE status=$1`, status).Scan(&total)
	} else {
		r.db.QueryRow(ctx, `SELECT COUNT(*) FROM jobs`).Scan(&total)
	}
	offset := (page - 1) * limit
	var rows interface{ Next() bool; Scan(...interface{}) error; Close() }
	var err error
	if status != "" {
		rows, err = r.db.Query(ctx, `SELECT `+jobCols+` FROM jobs WHERE status=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, `SELECT `+jobCols+` FROM jobs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil { return nil, 0, err }
	defer rows.Close()
	return scanJobs(rows), total, nil
}

func (r *jobRepo) Search(ctx context.Context, in *models.SearchInput) ([]*models.Job, int, error) {
	conds := []string{"status='active'"}
	args := []interface{}{}
	i := 1
	add := func(cond string, val interface{}) {
		conds = append(conds, fmt.Sprintf(cond, i))
		args = append(args, val)
		i++
	}
	if in.Query != ""           { add("(title ILIKE $%d OR description ILIKE $%d)", "%"+in.Query+"%"); i++ ; args = append(args, "%"+in.Query+"%") }
	if in.Location != ""        { add("location ILIKE $%d", "%"+in.Location+"%") }
	if in.Category != ""        { add("category=$%d", in.Category) }
	if in.JobType != ""         { add("job_type=$%d", in.JobType) }
	if in.ExperienceLevel != "" { add("experience_level=$%d", in.ExperienceLevel) }
	if in.SalaryMin > 0         { add("salary_max>=$%d", in.SalaryMin) }
	if in.SalaryMax > 0         { add("salary_min<=$%d", in.SalaryMax) }

	where := strings.Join(conds, " AND ")
	var total int
	r.db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM jobs WHERE %s`, where), args...).Scan(&total)

	args = append(args, in.Limit, (in.Page-1)*in.Limit)
	q := fmt.Sprintf(`SELECT `+jobCols+` FROM jobs WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, i, i+1)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	return scanJobs(rows), total, nil
}

func (r *jobRepo) GetByEmployer(ctx context.Context, empID string, page, limit int) ([]*models.Job, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM jobs WHERE employer_id=$1`, empID).Scan(&total)
	rows, err := r.db.Query(ctx, `SELECT `+jobCols+` FROM jobs WHERE employer_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		empID, limit, (page-1)*limit)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	return scanJobs(rows), total, nil
}

func (r *jobRepo) IncrementApplicationCount(ctx context.Context, jobID string) error {
	_, err := r.db.Exec(ctx, `UPDATE jobs SET application_count=application_count+1,updated_at=$1 WHERE id=$2`, time.Now(), jobID)
	return err
}

func scanJobs(rows interface{ Next() bool; Scan(...interface{}) error }) []*models.Job {
	var list []*models.Job
	for rows.Next() {
		j := &models.Job{}
		rows.Scan(&j.ID, &j.EmployerID, &j.Title, &j.Description, &j.Location, &j.JobType,
			&j.Category, &j.SalaryMin, &j.SalaryMax, &j.Currency, &j.RequirementsStr,
			&j.ExperienceLevel, &j.Status, &j.ApplicationCount, &j.CreatedAt, &j.UpdatedAt)
		j.ParseRequirements()
		list = append(list, j)
	}
	return list
}

// ── Application Repository ────────────────────────────────────────────────

type ApplicationRepository interface {
	Create(ctx context.Context, jobID, applicantID, coverLetter, resumeURL string) (*models.Application, error)
	GetByJob(ctx context.Context, jobID string, page, limit int) ([]*models.Application, int, error)
	GetByUser(ctx context.Context, userID string, page, limit int) ([]*models.Application, int, error)
	UpdateStatus(ctx context.Context, id, status string) (*models.Application, error)
	Exists(ctx context.Context, jobID, userID string) (bool, error)
}

type appRepo struct{ db *pgxpool.Pool }

func NewApplicationRepository(db *pgxpool.Pool) ApplicationRepository { return &appRepo{db: db} }

const appCols = `id,job_id,applicant_id,cover_letter,resume_url,status,applied_at,updated_at`

func scanApp(row interface{ Scan(...interface{}) error }) (*models.Application, error) {
	a := &models.Application{}
	err := row.Scan(&a.ID, &a.JobID, &a.ApplicantID, &a.CoverLetter, &a.ResumeURL, &a.Status, &a.AppliedAt, &a.UpdatedAt)
	return a, err
}

func (r *appRepo) Create(ctx context.Context, jobID, applicantID, coverLetter, resumeURL string) (*models.Application, error) {
	exists, _ := r.Exists(ctx, jobID, applicantID)
	if exists { return nil, fmt.Errorf("already applied") }
	now := time.Now()
	row := r.db.QueryRow(ctx, `
		INSERT INTO job_applications (id,job_id,applicant_id,cover_letter,resume_url,status,applied_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,'pending',$6,$7) RETURNING `+appCols,
		uuid.New().String(), jobID, applicantID, coverLetter, resumeURL, now, now)
	return scanApp(row)
}

func (r *appRepo) GetByJob(ctx context.Context, jobID string, page, limit int) ([]*models.Application, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM job_applications WHERE job_id=$1`, jobID).Scan(&total)
	rows, _ := r.db.Query(ctx, `SELECT `+appCols+` FROM job_applications WHERE job_id=$1 ORDER BY applied_at DESC LIMIT $2 OFFSET $3`,
		jobID, limit, (page-1)*limit)
	defer rows.Close()
	var list []*models.Application
	for rows.Next() {
		a := &models.Application{}
		rows.Scan(&a.ID, &a.JobID, &a.ApplicantID, &a.CoverLetter, &a.ResumeURL, &a.Status, &a.AppliedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, total, nil
}

func (r *appRepo) GetByUser(ctx context.Context, userID string, page, limit int) ([]*models.Application, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM job_applications WHERE applicant_id=$1`, userID).Scan(&total)
	rows, _ := r.db.Query(ctx, `SELECT `+appCols+` FROM job_applications WHERE applicant_id=$1 ORDER BY applied_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, (page-1)*limit)
	defer rows.Close()
	var list []*models.Application
	for rows.Next() {
		a := &models.Application{}
		rows.Scan(&a.ID, &a.JobID, &a.ApplicantID, &a.CoverLetter, &a.ResumeURL, &a.Status, &a.AppliedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, total, nil
}

func (r *appRepo) UpdateStatus(ctx context.Context, id, status string) (*models.Application, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE job_applications SET status=$1,updated_at=$2 WHERE id=$3 RETURNING `+appCols,
		status, time.Now(), id)
	return scanApp(row)
}

func (r *appRepo) Exists(ctx context.Context, jobID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM job_applications WHERE job_id=$1 AND applicant_id=$2)`, jobID, userID).Scan(&exists)
	return exists, err
}
