package handlers

import (
	"net/http"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/rules"
	"github.com/labstack/echo"
)

func GetSubmissions(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	subs, err := j.GetSubmissions()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, subs)
}

func GetSolves(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	solves, err := models.GetSolves(db)
	if err != nil {
		return err
	}

	result := make(map[int]int)
	for _, item := range solves {
		result[item.ChallengeID] = item.Solves
	}

	return c.JSON(http.StatusOK, result)
}
