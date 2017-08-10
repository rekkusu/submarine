package handlers

import (
	"net/http"
	"sort"

	"github.com/activedefense/submarine/rules"
	"github.com/activedefense/submarine/scoring"
	"github.com/labstack/echo"
)

func GetScoreboard(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)

	scoring := scoring.FixedJeopardy{j}
	scores := scoring.GetScores()

	sort.Stable(scores)

	return c.JSON(http.StatusOK, scores)
}
