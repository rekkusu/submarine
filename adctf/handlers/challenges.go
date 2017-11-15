package handlers

import (
	"net/http"
	"strconv"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/rules"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func GetChallenges(c echo.Context) error {
	jeopardy := c.Get("jeopardy").(rules.JeopardyRule)
	db := jeopardy.GetDB()
	team := c.Get("team").(*models.Team)

	if models.IsContestClosed(db) && !models.IsAdmin(team) {
		return echo.ErrForbidden
	}

	chals, err := models.GetChallengesWithSolves(db)
	if err != nil {
		return err
	}

	for i, c := range chals {
		chals[i].Flag = nil
		chals[i].Point = jeopardy.GetScoring().CalcScore(&c)
	}

	return c.JSON(http.StatusOK, chals)
}

func GetChallengeByID(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	team := c.Get("team").(*models.Team)
	if models.IsContestClosed(db) && !models.IsAdmin(team) {
		return echo.ErrForbidden
	}

	chal, err := models.GetChallenge(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, chal)
}

func CreateChallenge(c echo.Context) error {
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

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := chal.Save(db); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, chal)
}

func UpdateChallenge(c echo.Context) error {
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

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := chal.Save(db); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteChallenge(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	chal, err := models.GetChallenge(db, id)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	} else if err != nil {
		return err
	}

	if err := chal.Delete(db); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func Submit(c echo.Context) error {
	jeopardy := c.Get("jeopardy").(rules.JeopardyRule)
	db := jeopardy.GetDB()

	if !models.IsContestOpen(db) {
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
	team, err := models.GetTeam(db, team_id)
	if err != nil {
		return err
	}

	chal, err := models.GetChallenge(db, id)
	if err == gorm.ErrRecordNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	sub, err := chal.Submit(db, team, form.Answer)
	if err != nil {
		if err == models.ErrChallengeHasAlreadySolved {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return err
	}

	jeopardy.GetScoring().Recalculate()

	return c.JSON(http.StatusAccepted, sub)
}

func CreateCategory(c echo.Context) error {
	var category models.Category

	if err := c.Bind(&category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(category); err != nil {
		return err
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := category.Create(db); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c echo.Context) error {
	var category models.Category

	if err := c.Bind(&category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	category.ID = id

	if err := c.Validate(category); err != nil {
		return err
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := category.Save(db); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	category, err := models.GetCategory(db, id)
	if err == gorm.ErrRecordNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	} else if err != nil {
		return err
	}

	if err := category.Delete(db); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func GetCategories(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	cates, err := models.GetCategories(db)
	if err != nil {
		return err
	}

	result := make(map[int]string)
	for _, item := range cates {
		result[item.ID] = item.Name
	}

	return c.JSON(http.StatusOK, result)
}

func GetSolvedChallenges(c echo.Context) error {
	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	team := c.Get("team").(*models.Team)
	sub, err := models.GetSolvedChallenges(db, team.GetID())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sub)
}
