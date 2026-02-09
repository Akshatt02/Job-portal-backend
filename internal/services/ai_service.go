package services

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genai"
)

// NOTE: This file uses the official Google GenAI SDK for Go.
// Install it with: go get google.golang.org/genai
// The client reads the API key from GEMINI_API_KEY (or GOOGLE_API_KEY).

const (
	// model to use; quickstart uses "gemini-3-flash-preview" in examples.
	geminiModel = "gemini-3-flash-preview"
)

// ExtractSkillsFromText calls Gemini and asks it to return a JSON array of skills.
func ExtractSkillsFromText(ctx context.Context, bio string) ([]string, error) {
	if strings.TrimSpace(bio) == "" {
		return nil, errors.New("bio is empty")
	}

	system := "You are a helpful assistant that extracts relevant professional skills from a textual bio. Return the top skills as a JSON array only. Use short skill names (e.g., go, react, nodejs, postgres)."
	user := "Extract top skills from the following bio and return ONLY a JSON array (e.g. [\"go\",\"react\"]). Do not add any explanation or text.\n\nBIO:\n" + bio
	prompt := system + "\n\n" + user

	out, err := callGenAI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// parse JSON array from output (tolerant)
	return parseStringArray(out)
}

// ComputeMatchScore returns an integer 0..100 score for how well userSkills match the jobDescription.
func ComputeMatchScore(ctx context.Context, userSkills []string, jobDescription string) (int, error) {
	sys := "You are a helpful assistant that scores how well a candidate's skills match a job."
	userPrompt := "Given the user's skills JSON array:\n" + toJSONString(userSkills) + "\n\nAnd the job description below:\n" + jobDescription + "\n\nReturn ONLY a JSON object with a single numeric field `match_score` with an integer value between 0 and 100 indicating the match percentage. Example: {\"match_score\":78}. Return no other text."

	prompt := sys + "\n\n" + userPrompt

	out, err := callGenAI(ctx, prompt)
	if err != nil {
		return 0, err
	}

	return parseScore(out)
}

// ----------------- GenAI call helper -----------------

// callGenAI uses the official genai client to call Gemini. The client reads GEMINI_API_KEY from the environment.
// It includes a small retry loop for transient failures.
func callGenAI(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	maxRetries := 3
	backoff := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		res, err := client.Models.GenerateContent(
			ctx,
			geminiModel,
			genai.Text(prompt),
			nil,
		)
		if err == nil {
			return res.Text(), nil
		}

		if attempt == maxRetries {
			return "", err
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return "", errors.New("genai failed after retries")
}

// ----------------- parsing helpers (same tolerant logic you had) -----------------

func parseStringArray(out string) ([]string, error) {
	s := strings.TrimSpace(out)
	arrText, err := firstJSONFragment(s)
	if err != nil {
		// fallback: maybe the whole output is a JSON array
		var fallback []string
		if err2 := json.Unmarshal([]byte(s), &fallback); err2 == nil {
			return fallback, nil
		}
		return nil, err
	}

	var arr []string
	if err := json.Unmarshal([]byte(arrText), &arr); err != nil {
		// try single quotes -> double
		clean := strings.ReplaceAll(arrText, "'", "\"")
		if err2 := json.Unmarshal([]byte(clean), &arr); err2 == nil {
			return arr, nil
		}
		return nil, err
	}
	return arr, nil
}

func parseScore(out string) (int, error) {
	s := strings.TrimSpace(out)
	jsonText, err := firstJSONFragment(s)
	if err != nil {
		// fallback: find any integer in text
		if n, err2 := extractFirstInt(s); err2 == nil {
			return clampScore(n), nil
		}
		return 0, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &m); err != nil {
		if n, err2 := extractFirstInt(jsonText); err2 == nil {
			return clampScore(n), nil
		}
		return 0, err
	}

	if v, ok := m["match_score"]; ok {
		switch t := v.(type) {
		case float64:
			return clampScore(int(t)), nil
		case string:
			if n, err := strconv.Atoi(t); err == nil {
				return clampScore(n), nil
			}
		case int:
			return clampScore(t), nil
		}
	}

	// final fallback
	if n, err := extractFirstInt(jsonText); err == nil {
		return clampScore(n), nil
	}
	return 0, errors.New("could not parse match score")
}

func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func firstJSONFragment(s string) (string, error) {
	// find either a leading array [ ... ] or object { ... }
	idxObj := strings.Index(s, "{")
	idxArr := strings.Index(s, "[")
	if idxObj == -1 && idxArr == -1 {
		return "", errors.New("no json found")
	}

	// prefer whichever appears first
	if idxObj != -1 && (idxArr == -1 || idxObj < idxArr) {
		last := strings.LastIndex(s, "}")
		if last == -1 || last <= idxObj {
			return "", errors.New("malformed json object")
		}
		return s[idxObj : last+1], nil
	}

	last := strings.LastIndex(s, "]")
	if last == -1 || last <= idxArr {
		return "", errors.New("malformed json array")
	}
	return s[idxArr : last+1], nil
}

func extractFirstInt(s string) (int, error) {
	re := regexp.MustCompile(`\d{1,3}`)
	m := re.FindString(s)
	if m == "" {
		return 0, errors.New("no integer found")
	}
	return strconv.Atoi(m)
}

func clampScore(n int) int {
	if n < 0 {
		return 0
	}
	if n > 100 {
		return 100
	}
	return n
}
