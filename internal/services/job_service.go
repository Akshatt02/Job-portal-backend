package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Akshatt02/job-portal-backend/internal/db"
	"github.com/Akshatt02/job-portal-backend/internal/models"
	"github.com/google/uuid"
)

var ErrJobNotFound = errors.New("job not found")

// CreateJob creates a job row and returns the created job id.
func CreateJob(title, description string, skills []string, salary, location, userIDStr, paymentTx string) (string, error) {
	// Ensure user id is valid uuid
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", err
	}

	// Enforce a payment_tx_hash for posting (as per assignment). Remove if not desired.
	if paymentTx == "" {
		return "", errors.New("payment required before posting job (payment_tx_hash missing)")
	}

	jobID := uuid.New()
	skillsBytes := []byte("null")
	if skills != nil {
		b, _ := json.Marshal(skills)
		skillsBytes = b
	}

	_, err = db.Pool.Exec(context.Background(),
		`INSERT INTO jobs (id, title, description, skills, salary, location, user_id, payment_tx_hash, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		jobID, title, description, skillsBytes, salary, location, userID, paymentTx, time.Now(),
	)
	if err != nil {
		return "", err
	}

	return jobID.String(), nil
}

// ListJobs returns up to `limit` recent jobs. If limit is 0, defaults to 100.
func ListJobs(limit int) ([]*models.Job, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := db.Pool.Query(context.Background(),
		`SELECT id, title, description, skills, salary, location, user_id, payment_tx_hash, created_at
		 FROM jobs
		 ORDER BY created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []*models.Job{}
	for rows.Next() {
		var (
			id uuid.UUID
			title, description string
			skillsRaw []byte
			salary, location string
			userID uuid.UUID
			paymentTx *string
			createdAt time.Time
		)
		err := rows.Scan(&id, &title, &description, &skillsRaw, &salary, &location, &userID, &paymentTx, &createdAt)
		if err != nil {
			return nil, err
		}

		var skills []string
		if len(skillsRaw) > 0 {
			_ = json.Unmarshal(skillsRaw, &skills)
		}

		px := ""
		if paymentTx != nil {
			px = *paymentTx
		}

		job := &models.Job{
			ID: id,
			Title: title,
			Description: description,
			Skills: skills,
			Salary: salary,
			Location: location,
			UserID: userID,
			PaymentTxHash: px,
			CreatedAt: createdAt,
		}
		out = append(out, job)
	}
	return out, nil
}

func GetJobByID(jobIDStr string) (*models.Job, error) {
	id, err := uuid.Parse(jobIDStr)
	if err != nil {
		return nil, err
	}

	var (
		title, description string
		skillsRaw []byte
		salary, location string
		userID uuid.UUID
		paymentTx *string
		createdAt time.Time
	)
	err = db.Pool.QueryRow(context.Background(),
		`SELECT title, description, skills, salary, location, user_id, payment_tx_hash, created_at
		 FROM jobs WHERE id=$1`, id).
		Scan(&title, &description, &skillsRaw, &salary, &location, &userID, &paymentTx, &createdAt)
	if err != nil {
		return nil, ErrJobNotFound
	}

	var skills []string
	if len(skillsRaw) > 0 {
		_ = json.Unmarshal(skillsRaw, &skills)
	}

	pt := ""
	if paymentTx != nil {
		pt = *paymentTx
	}

	j := &models.Job{
		ID: id,
		Title: title,
		Description: description,
		Skills: skills,
		Salary: salary,
		Location: location,
		UserID: userID,
		PaymentTxHash: pt,
		CreatedAt: createdAt,
	}
	return j, nil
}
