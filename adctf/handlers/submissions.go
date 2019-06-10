package handlers

import (
	"github.com/activedefense/submarine/adctf/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func (h *Handler) GetSubmissions(c echo.Context) error {
	offset := 0
	if param, err := strconv.Atoi(c.QueryParam("offset")); err == nil {
		offset = param
	}

	limit := 50
	if param, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = param
	}

	count, err := models.GetSubmissionCount(h.DB)
	if err != nil {
		return err
	}

	subs, err := models.GetSubmissions(h.DB, offset, limit)
	if err != nil {
		return err
	}

	result := struct{
		Submissions []models.Submission `json:"submissions"`
		Total int `json:"total"`
	}{
		subs,
		count,
	}

	return c.JSON(http.StatusOK, result)
}
