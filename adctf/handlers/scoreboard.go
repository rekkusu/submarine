package handlers

import (
	"net/http"

	"github.com/activedefense/submarine/rules"
	"github.com/activedefense/submarine/scoring"
	"github.com/labstack/echo"
)

func GetScoreboard(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)

	score := scoring.FixedJeopardy{j}
	ranks := score.GetRanking()
	return c.JSON(http.StatusOK, ranks)
}
