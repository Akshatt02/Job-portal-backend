package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Akshatt02/job-portal-backend/internal/db"
	"github.com/Akshatt02/job-portal-backend/internal/models"
	"github.com/Akshatt02/job-portal-backend/pkg/utils"
	"github.com/google/uuid"
)

func RegisterUser(name, email, password string) (string, error) {
	// check if email exists
	var exists bool
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("email already registered")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	id := uuid.New()

	_, err = db.Pool.Exec(context.Background(),
		`INSERT INTO users (id, name, email, password_hash, created_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		id, name, email, hash, time.Now(),
	)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func LoginUser(email, password string) (string, error) {
	var id uuid.UUID
	var hash string

	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE email=$1",
		email,
	).Scan(&id, &hash)

	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, hash) {
		return "", errors.New("invalid credentials")
	}

	return id.String(), nil
}

func GetUserByID(userID string) (*models.User, error) {
	var (
		id           uuid.UUID
		name         string
		email        string
		bio          *string
		linkedin     *string
		skillsRaw    []byte
		wallet       *string
		createdAt    time.Time
	)

	err := db.Pool.QueryRow(context.Background(),
		`SELECT id, name, email, bio, linkedin_url, skills, wallet_address, created_at
		 FROM users WHERE id=$1`, userID,
	).Scan(&id, &name, &email, &bio, &linkedin, &skillsRaw, &wallet, &createdAt)

	if err != nil {
		return nil, err
	}

	var skills []string
	if len(skillsRaw) > 0 {
		_ = json.Unmarshal(skillsRaw, &skills)
	}

	u := &models.User{
		ID:            id,
		Name:          name,
		Email:         email,
		Bio:           safeStr(bio),
		LinkedinURL:   safeStr(linkedin),
		Skills:        skills,
		WalletAddress: safeStr(wallet),
		CreatedAt:     createdAt,
	}
	return u, nil
}

func UpdateUser(userID string, updates map[string]interface{}) error {
	// Build update dynamically but safely.
	// Allowed fields: name, bio, linkedin_url, skills ([]string), wallet_address
	args := []interface{}{}
	setClauses := []string{}
	argIdx := 1

	if v, ok := updates["name"].(string); ok {
		setClauses = append(setClauses, `name = $`+itoa(argIdx))
		args = append(args, v); argIdx++
	}
	if v, ok := updates["bio"].(string); ok {
		setClauses = append(setClauses, `bio = $`+itoa(argIdx))
		args = append(args, v); argIdx++
	}
	if v, ok := updates["linkedin_url"].(string); ok {
		setClauses = append(setClauses, `linkedin_url = $`+itoa(argIdx))
		args = append(args, v); argIdx++
	}
	if v, ok := updates["wallet_address"].(string); ok {
		setClauses = append(setClauses, `wallet_address = $`+itoa(argIdx))
		args = append(args, v); argIdx++
	}
	if v, ok := updates["skills"].([]string); ok {
		// marshal to JSON and set
		skillsBytes, _ := json.Marshal(v)
		setClauses = append(setClauses, `skills = $`+itoa(argIdx))
		args = append(args, skillsBytes); argIdx++
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	// Append user id
	args = append(args, userID)
	query := `UPDATE users SET ` + join(setClauses, ", ") + ` WHERE id = $` + itoa(argIdx)

	_, err := db.Pool.Exec(context.Background(), query, args...)
	return err
}

// small helpers
func safeStr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
func itoa(i int) string { return fmt.Sprintf("%d", i) }
func join(arr []string, sep string) string {
	out := ""
	for i, s := range arr {
		if i != 0 {
			out += sep
		}
		out += s
	}
	return out
}
