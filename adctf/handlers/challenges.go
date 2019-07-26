package handlers

import (
	"net/http"
	"strconv"

	"github.com/activedefense/submarine/adctf/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetChallenges(c echo.Context) error {
	team := c.Get("team").(*models.User)

	if models.IsContestClosed(h.DB) && !models.IsAdmin(team) {
		return echo.ErrForbidden
	}

	chals, err := models.GetChallengesWithSolves(h.DB)
	if err != nil {
		return err
	}

	for i, c := range chals {
		chals[i].Flag = nil
		chals[i].Point = h.Scoring.CalcScore(c)
	}

	return c.JSON(http.StatusOK, chals)
}

func (h *Handler) GetChallengeByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	team := c.Get("team").(*models.User)
	if models.IsContestClosed(h.DB) && !models.IsAdmin(team) {
		return echo.ErrForbidden
	}

	chal, err := models.GetChallenge(h.DB, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, chal)
}

func (h *Handler) CreateChallenge(c echo.Context) error {
	var form struct {
		CategoryID  int    `json:"category_id"`
		Title       string `json:"title"`
		Point       int    `json:"point"`
		Description string `json:"description"`
		Flag        string `json:"flag"`
	}

	if err := c.Bind(&form); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "bad request"}
	}

	chal := &models.Challenge{
		CategoryID:  form.CategoryID,
		Title:       form.Title,
		Point:       form.Point,
		Description: form.Description,
		Flag:        &form.Flag,
	}

	if err := chal.Save(h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, chal)
}

func (h *Handler) UpdateChallenge(c echo.Context) error {
	var form struct {
		CategoryID  int    `json:"category_id"`
		Title       string `json:"title"`
		Point       int    `json:"point"`
		Description string `json:"description"`
		Flag        string `json:"flag"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	chal := &models.Challenge{
		ID:          id,
		CategoryID:  form.CategoryID,
		Title:       form.Title,
		Point:       form.Point,
		Description: form.Description,
		Flag:        &form.Flag,
	}

	if err := chal.Save(h.DB); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteChallenge(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	chal, err := models.GetChallenge(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	} else if err != nil {
		return err
	}

	if err := chal.Delete(h.DB); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Submit(c echo.Context) error {
	if !models.IsContestOpen(h.DB) {
		return echo.ErrForbidden
	}

	var form struct {
		Answer string `json:"answer"`
	}

	if err := c.Bind(&form); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	team_id := int(claims["user"].(float64))
	team, err := models.GetTeam(h.DB, team_id)
	if err != nil {
		return err
	}

	chal, err := models.GetChallenge(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	sub, err := chal.Submit(h.DB, team, team, form.Answer)
	if err != nil {
		if err == models.ErrChallengeHasAlreadySolved {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return err
	}

	h.Scoring.Update()

	return c.JSON(http.StatusAccepted, sub)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	var category models.Category

	if err := c.Bind(&category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(category); err != nil {
		return err
	}

	if err := category.Create(h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *Handler) UpdateCategory(c echo.Context) error {
	var category models.Category

	if err := c.Bind(&category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	category.ID = id

	if err := c.Validate(category); err != nil {
		return err
	}

	_, err = models.GetCategory(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := category.Save(h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	category, err := models.GetCategory(h.DB, id)
	if err == gorm.ErrRecordNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	if err := category.Delete(h.DB); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetCategories(c echo.Context) error {
	cates, err := models.GetCategories(h.DB)
	if err != nil {
		return err
	}

	result := make(map[int]string)
	for _, item := range cates {
		result[item.ID] = item.Name
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) GetSolvedChallenges(c echo.Context) error {
	team := c.Get("team").(*models.User)
	sub, err := models.GetSolvedChallenges(h.DB, team.GetID())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sub)
}
