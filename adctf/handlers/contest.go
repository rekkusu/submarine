package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetContestInfo(c echo.Context) error {
	info, err := models.GetContestInfo(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info)
}

func (h *Handler) PutContestInfo(c echo.Context) error {
	var info models.ContestInfo
	if err := c.Bind(&info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if err := models.SetContestStatus(h.DB, info); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetAllAnnouncements(c echo.Context) error {
	announcements, err := models.GetAllAnnouncements(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcements)
}

func (h *Handler) GetCurrentAnnouncements(c echo.Context) error {
	announcements, err := models.GetCurrentAnnouncements(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, announcements)
}

func (h *Handler) GetAnnouncement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	announcement, err := models.GetAnnouncement(h.DB, id)
	if err != nil {
		return err
	}

	if announcement == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, announcement)
}

func (h *Handler) NewAnnouncement(c echo.Context) error {
	announcement := parseAnnouncement(c)
	if announcement == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	announcement.ID = 0
	if announcement.PostedAt.Before(time.Now()) {
		announcement.PostedAt = time.Now()
	}

	if err := announcement.Create(h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, announcement)
}

func (h *Handler) EditAnnouncement(c echo.Context) error {
	announcement := parseAnnouncement(c)
	if announcement == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	announcement.ID = id

	if err := announcement.Save(h.DB); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteAnnouncement(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	tx := h.DB.Begin()
	announcement, err := models.GetAnnouncement(tx, id)
	if err != nil {
		tx.Rollback()
		return echo.ErrNotFound
	}

	if err := announcement.Delete(tx); err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "error")
	}

	tx.Commit()

	return c.NoContent(http.StatusNoContent)
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
