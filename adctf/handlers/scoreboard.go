package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/activedefense/submarine/rules"
	"github.com/labstack/echo"
)

func GetScoreboard(c echo.Context) error {
	type record struct {
		Order          int       `json:"order"`
		ID             int       `json:"id"`
		Name           string    `json:"name"`
		Score          int       `json:"score"`
		LastSubmission time.Time `json:"last"`
	}

	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	scores := j.GetScoring().GetScores()

	sort.Stable(scores)

	result := make([]record, len(scores))
	for i, item := range scores {
		result[i] = record{
			Order:          i + 1,
			ID:             item.GetTeam().GetID(),
			Name:           item.GetTeam().GetName(),
			Score:          item.GetScore(),
			LastSubmission: item.GetLastSubmission(),
		}
	}

	return c.JSON(http.StatusOK, result)
}
