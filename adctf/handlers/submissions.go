package handlers

import (
	"net/http"

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
