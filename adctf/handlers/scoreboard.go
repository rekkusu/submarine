package handlers

import (
	"github.com/activedefense/submarine/adctf/models"
	"net/http"
	"sort"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetScoreboard(c echo.Context) error {
	type record struct {
		Order          int       `json:"order"`
		Team           ctf.Team  `json:"team"`
		Score          int       `json:"score"`
		LastSubmission time.Time `json:"last"`
	}

	challenges, _ := models.GetChallengesWithSolves(h.DB)
	submissions, _ := models.GetCorrectSubmissions(h.DB)
	teams, _ := models.GetTeams(h.DB)

	scores := h.Scoring.GetScores(challenges, submissions, teams)

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].GetScore() == scores[j].GetScore() {
			return scores[i].GetLastSubmission().Before(scores[j].GetLastSubmission())
		}
		return scores[i].GetScore() > scores[j].GetScore()
	})

	result := make([]record, len(scores))
	for i, item := range scores {
		result[i] = record{
			Order:          i + 1,
			Team:           item.Team,
			Score:          item.Score,
			LastSubmission: item.LastSubmission,
		}
	}

	return c.JSON(http.StatusOK, result)
}
