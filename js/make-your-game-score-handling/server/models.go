package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Score represents a player's score entry
type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Time  string `json:"time"`
}

// ScoreWithRank includes ranking information for display
type ScoreWithRank struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Time     string `json:"time"`
	Rank     int    `json:"rank"`
	RankText string `json:"rankText"` // "1st", "2nd", "3rd", etc.
}

// ScoresResponse contains paginated scores data
type ScoresResponse struct {
	Scores      []ScoreWithRank `json:"scores"`
	CurrentPage int             `json:"currentPage"`
	TotalPages  int             `json:"totalPages"`
	TotalScores int             `json:"totalScores"`
	PerPage     int             `json:"perPage"`
}

// SubmitScoreResponse contains ranking feedback after score submission
type SubmitScoreResponse struct {
	Status      string  `json:"status"`
	PlayerRank  int     `json:"playerRank"`
	TotalScores int     `json:"totalScores"`
	Percentage  float64 `json:"percentage"`
	RankMessage string  `json:"rankMessage"`
}

const SCORES_FILE = "scores/scores.json"
const SCORES_PER_PAGE = 5

// LoadScores reads all scores from the JSON file
func LoadScores() ([]Score, error) {
	if _, err := os.Stat(SCORES_FILE); os.IsNotExist(err) {
		return []Score{}, nil
	}

	data, err := os.ReadFile(SCORES_FILE)
	if err != nil {
		return nil, fmt.Errorf("failed to read scores file: %w", err)
	}

	if len(data) == 0 {
		return []Score{}, nil
	}

	var scores []Score
	if err := json.Unmarshal(data, &scores); err != nil {
		return nil, fmt.Errorf("failed to parse scores JSON: %w", err)
	}

	return scores, nil
}

// SaveScores writes all scores to the JSON file
func SaveScores(scores []Score) error {
	if err := os.MkdirAll("scores", 0755); err != nil {
		return fmt.Errorf("failed to create scores directory: %w", err)
	}

	data, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scores: %w", err)
	}

	if err := os.WriteFile(SCORES_FILE, data, 0644); err != nil {
		return fmt.Errorf("failed to write scores file: %w", err)
	}

	return nil
}

// SortScoresByHighest sorts scores in descending order (highest first)
func SortScoresByHighest(scores []Score) {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
}

// GetRankText converts a numeric rank to ordinal text (1st, 2nd, 3rd, etc.)
func GetRankText(rank int) string {
	if rank >= 11 && rank <= 13 {
		return fmt.Sprintf("%dth", rank)
	}
	
	switch rank % 10 {
	case 1:
		return fmt.Sprintf("%dst", rank)
	case 2:
		return fmt.Sprintf("%dnd", rank)
	case 3:
		return fmt.Sprintf("%drd", rank)
	default:
		return fmt.Sprintf("%dth", rank)
	}
}

// AddScoreWithRanking adds a score and returns ranking information
func AddScoreWithRanking(newScore Score) (*SubmitScoreResponse, error) {
	scores, err := LoadScores()
	if err != nil {
		return nil, err
	}

	// Add the new score
	scores = append(scores, newScore)
	SortScoresByHighest(scores)

	// Find the player's rank
	var playerRank int
	for i, score := range scores {
		if score.Name == newScore.Name && score.Score == newScore.Score && score.Time == newScore.Time {
			playerRank = i + 1
			break
		}
	}

	totalScores := len(scores)
	percentage := float64(totalScores-playerRank+1) / float64(totalScores) * 100

	// Create ranking message
	var rankMessage string
	if playerRank == 1 {
		rankMessage = fmt.Sprintf("🥇 Amazing! You achieved 1st place with %d points!", newScore.Score)
	} else if playerRank <= 3 {
		rankMessage = fmt.Sprintf("🏆 Fantastic! You're in %s place - that's top %d%%!", GetRankText(playerRank), int(percentage))
	} else if percentage >= 90 {
		rankMessage = fmt.Sprintf("🌟 Excellent! You're in the top %.0f%%, ranking %s out of %d players!", percentage, GetRankText(playerRank), totalScores)
	} else if percentage >= 75 {
		rankMessage = fmt.Sprintf("👏 Great job! You're in the top %.0f%%, ranking %s out of %d players!", percentage, GetRankText(playerRank), totalScores)
	} else if percentage >= 50 {
		rankMessage = fmt.Sprintf("💪 Good work! You're in the top %.0f%%, ranking %s out of %d players!", percentage, GetRankText(playerRank), totalScores)
	} else {
		rankMessage = fmt.Sprintf("🎯 Nice try! You ranked %s out of %d players. Keep practicing!", GetRankText(playerRank), totalScores)
	}

	// Save updated scores
	if err := SaveScores(scores); err != nil {
		return nil, err
	}

	return &SubmitScoreResponse{
		Status:      "success",
		PlayerRank:  playerRank,
		TotalScores: totalScores,
		Percentage:  percentage,
		RankMessage: rankMessage,
	}, nil
}

// GetPaginatedScores returns scores for a specific page
func GetPaginatedScores(page int) (*ScoresResponse, error) {
	scores, err := LoadScores()
	if err != nil {
		return nil, err
	}

	SortScoresByHighest(scores)

	totalScores := len(scores)
	totalPages := (totalScores + SCORES_PER_PAGE - 1) / SCORES_PER_PAGE

	// Handle edge cases
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate pagination bounds
	startIdx := (page - 1) * SCORES_PER_PAGE
	endIdx := startIdx + SCORES_PER_PAGE
	if endIdx > totalScores {
		endIdx = totalScores
	}

	// Convert to ScoreWithRank
	var scoresWithRank []ScoreWithRank
	if startIdx < totalScores {
		for i := startIdx; i < endIdx; i++ {
			score := scores[i]
			rank := i + 1
			scoresWithRank = append(scoresWithRank, ScoreWithRank{
				Name:     score.Name,
				Score:    score.Score,
				Time:     score.Time,
				Rank:     rank,
				RankText: GetRankText(rank),
			})
		}
	}

	return &ScoresResponse{
		Scores:      scoresWithRank,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalScores: totalScores,
		PerPage:     SCORES_PER_PAGE,
	}, nil
}