package handlers

import (
	"net/http"
	"sort"
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

	challenges, _ := h.Jeopardy.GetChallenges(h.DB)
	submissions, _ := h.Jeopardy.GetSubmissions(h.DB)
	teams, _ := h.Jeopardy.GetTeams(h.DB)

	scores := h.Jeopardy.GetScoring().GetScores(challenges, submissions, teams)

	sort.Stable(scores)

	result := make([]record, len(scores))
	for i, item := range scores {
		result[i] = record{
			Order:          i + 1,
			Team:           item.GetTeam(),
			Score:          item.GetScore(),
			LastSubmission: item.GetLastSubmission(),
		}
	}

	return c.JSON(http.StatusOK, result)
}
