package handlers

import (
	"github.com/activedefense/submarine/adctf/models"
	"github.com/labstack/echo"
	"net/http"
)

func (h *Handler) GetSubmissions(c echo.Context) error {
	subs, err := models.GetSubmissions(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, subs)
}
