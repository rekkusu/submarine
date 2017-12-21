package handlers

import (
	"net/http"
	"sort"

	"github.com/activedefense/submarine/rules"
	"github.com/labstack/echo"
)

func GetScoreboard(c echo.Context) error {
	type record struct {
		Order int    `json:"order"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	scores := j.GetScoring().GetScores()

	sort.Stable(scores)

	result := make([]record, len(scores))
	for i, item := range scores {
		result[i] = record{
			Order: i + 1,
			ID:    item.GetTeam().GetID(),
			Name:  item.GetTeam().GetName(),
			Score: item.GetScore(),
		}
	}

	return c.JSON(http.StatusOK, result)
}
