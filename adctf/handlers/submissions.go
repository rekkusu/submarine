package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) GetSubmissions(c echo.Context) error {
	subs, err := h.Jeopardy.GetSubmissions(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, subs)
}
