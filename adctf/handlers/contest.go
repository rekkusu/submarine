package handlers

import (
	"net/http"
	"time"

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

func GetAllAnnouncements(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	announcements, err := models.GetAllAnnouncements(db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcements)
}

func GetCurrentAnnouncements(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	team := c.Get("team").(*models.Team)

	var announcements []models.Announcement
	var err error
	if models.IsAdmin(team) {
		announcements, err = models.GetAllAnnouncements(db)
	} else {
		announcements, err = models.GetCurrentAnnouncements(db)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcements)
}

func NewAnnouncement(c echo.Context) error {
	var announcement models.Announcement
	if err := c.Bind(&announcement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if announcement.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	announcement.ID = 0
	if announcement.PostedAt.Before(time.Now()) {
		announcement.PostedAt = time.Now()
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := announcement.Create(db); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcement)
}
