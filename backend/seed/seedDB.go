package seed

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"OJ-backend/config"
	model "OJ-backend/models"

	"github.com/google/uuid"
)

func SeedDB() error {
	db := config.DB

	file, err := os.ReadFile("seed/data.json")
	if err != nil {
		return fmt.Errorf("failed to read seed data file: %w", err)
	}

	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(file, &rawData); err != nil {
		return fmt.Errorf("failed to unmarshal raw data: %w", err)
	}

	// Seed Languages
	var languages []model.Language
	if err := json.Unmarshal(rawData["languages"], &languages); err != nil {
		log.Fatalf("Failed to unmarshal languages: %v", err)
		return err
	}
	for _, lang := range languages {
		if err := db.Create(&lang).Error; err != nil {
			return fmt.Errorf("failed to seed language %s: %w", lang.Name, err)
		}
	}

	// Seed Users
	var usersJSON []map[string]interface{}
	if err := json.Unmarshal(rawData["users"], &usersJSON); err != nil {
		log.Fatalf("Failed to unmarshal users: %v", err)
		return err
	}
	for _, u := range usersJSON {
		id, _ := uuid.Parse(u["id"].(string))
		createdAt, _ := time.Parse(time.RFC3339, u["created_at"].(string))
		user := model.User{
			ID:        id,
			Username:  u["username"].(string),
			Email:     u["email"].(string),
			OauthID:   u["oauth_id"].(string),
			Provider:  u["provider"].(string),
			Image:     u["image"].(string),
			CreatedAt: createdAt,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Username, err)
		}
	}

	// Seed Contests
	var contestsJSON []map[string]interface{}
	if err := json.Unmarshal(rawData["contests"], &contestsJSON); err != nil {
		log.Fatalf("Failed to unmarshal contests: %v", err)
		return err
	}
	for _, c := range contestsJSON {
		id, _ := uuid.Parse(c["id"].(string))
		startTime, _ := time.Parse(time.RFC3339, c["start_time"].(string))
		endTime, _ := time.Parse(time.RFC3339, c["end_time"].(string))
		createdAt, _ := time.Parse(time.RFC3339, c["created_at"].(string))
		contest := model.Contest{
			ID:          id,
			Name:        c["name"].(string),
			Description: c["description"].(string),
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedAt:   createdAt,
		}
		if err := db.Create(&contest).Error; err != nil {
			return fmt.Errorf("failed to seed contest %s: %w", contest.Name, err)
		}
	}

	// Seed Problems
	var problemsJSON []map[string]interface{}
	if err := json.Unmarshal(rawData["problems"], &problemsJSON); err != nil {
		log.Fatalf("Failed to unmarshal problems: %v", err)
		return err
	}
	for _, p := range problemsJSON {
		id, _ := uuid.Parse(p["id"].(string))
		contestID, _ := uuid.Parse(p["contest_id"].(string))
		createdAt, _ := time.Parse(time.RFC3339, p["created_at"].(string))
		problem := model.Problem{
			ID:          id,
			ContestID:   contestID,
			Title:       p["title"].(string),
			Description: p["description"].(string),
			CreatedAt:   createdAt,
		}
		if err := db.Create(&problem).Error; err != nil {
			return fmt.Errorf("failed to seed problem %s: %w", problem.Title, err)
		}
	}

	// Seed TestCases
	var testCasesJSON []map[string]interface{}
	if err := json.Unmarshal(rawData["test_cases"], &testCasesJSON); err != nil {
		log.Fatalf("Failed to unmarshal test cases: %v", err)
		return err
	}
	for _, tc := range testCasesJSON {
		id, _ := uuid.Parse(tc["id"].(string))
		problemID, _ := uuid.Parse(tc["problem_id"].(string))
		createdAt, _ := time.Parse(time.RFC3339, tc["created_at"].(string))
		testCase := model.TestCase{
			ID:        id,
			ProblemID: problemID,
			Input:     tc["input"].(string),
			Output:    tc["output"].(string),
			CreatedAt: createdAt,
		}
		if err := db.Create(&testCase).Error; err != nil {
			return fmt.Errorf("failed to seed test case for problem %s: %w", testCase.ProblemID, err)
		}
	}

	// Seed Submissions
	var submissionsJSON []map[string]interface{}
	if err := json.Unmarshal(rawData["submissions"], &submissionsJSON); err != nil {
		log.Fatalf("Failed to unmarshal submissions: %v", err)
		return err
	}
	for _, s := range submissionsJSON {
		id, _ := uuid.Parse(s["id"].(string))
		problemID, _ := uuid.Parse(s["problem_id"].(string))
		userID, _ := uuid.Parse(s["user_id"].(string))
		contestID, _ := uuid.Parse(s["contest_id"].(string))
		submittedAt, _ := time.Parse(time.RFC3339, s["submitted_at"].(string))
		createdAt, _ := time.Parse(time.RFC3339, s["created_at"].(string))
		wrongTestCase, _ := uuid.Parse(s["wrong_test_case"].(string))
		submission := model.Submission{
			ID:            id,
			ProblemID:     problemID,
			UserID:        userID,
			ContestID:     contestID,
			SubmittedAt:   submittedAt,
			Result:        s["result"].(string),
			Language:      s["language"].(string),
			SourceCode:    s["source_code"].(string),
			Score:         int(s["score"].(float64)),
			StdOutput:     s["std_output"].(string),
			StdError:      s["std_error"].(string),
			WrongTestCase: wrongTestCase,
			CreatedAt:     createdAt,
			CompileOutput: s["compile_output"].(string),
			ExitCode:      int(s["exit_code"].(float64)),
		}
		if err := db.Create(&submission).Error; err != nil {
			return fmt.Errorf("failed to seed submission %s: %w", submission.ID, err)
		}
	}

	log.Println("Database seeding completed.")
	return nil
}
