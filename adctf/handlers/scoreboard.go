package handlers

import (
	"github.com/activedefense/submarine/adctf/models"
	"net/http"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/labstack/echo"
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

	//sort.Stable(scores)

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
