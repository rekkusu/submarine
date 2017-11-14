package handlers

import (
	"net/http"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/rules"
	"github.com/labstack/echo"
)

func GetContestInfo(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	info, err := models.GetContestInfo(db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info)
}

func PutContestInfo(c echo.Context) error {
	var info models.ContestInfo
	if err := c.Bind(&info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := models.SetContestStatus(db, info); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
