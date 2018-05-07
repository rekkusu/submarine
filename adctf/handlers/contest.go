package handlers

import (
	"net/http"
	"strconv"
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

func GetAnnouncement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	announcement, err := models.GetAnnouncement(db, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcement)
}

func NewAnnouncement(c echo.Context) error {
	announcement := parseAnnouncement(c)
	if announcement == nil {
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

func EditAnnouncement(c echo.Context) error {
	announcement := parseAnnouncement(c)
	if announcement == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	announcement.ID = id

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := announcement.Save(db); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcement)
}

func parseAnnouncement(c echo.Context) *models.Announcement {
	var announcement models.Announcement
	if err := c.Bind(&announcement); err != nil {
		return nil
	}

	if err := c.Validate(announcement); err != nil {
		return nil
	}

	return &announcement
}
