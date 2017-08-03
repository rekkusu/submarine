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
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	chals, err := j.GetChallengeStore().All()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, chals)
}

func GetChallengeByID(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	chal, err := j.GetChallengeStore().Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, chal)
}

func NewChallenge(c echo.Context) error {
	var form struct {
		Title       string `json:"title"`
		Point       int    `json:"point"`
		Description string `json:"description"`
		Flag        string `json:"flag"`
	}

	if err := c.Bind(&form); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "bad request"}
	}

	chal := &models.Challenge{
		Title:       form.Title,
		Point:       form.Point,
		Description: form.Description,
		Flag:        form.Flag,
	}

	store := c.Get("jeopardy").(rules.JeopardyRule).GetChallengeStore()
	if err := store.Save(chal); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, chal)
}

func Submit(c echo.Context) error {
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

	j, _ := c.Get("jeopardy").(rules.JeopardyRule)

	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	team_id := int(claims["user"].(float64))
	team, err := j.GetTeamStore().GetTeam(team_id)
	if err != nil {
		return err
	}

	chal, err := j.GetChallengeStore().Get(id)
	if err == gorm.ErrRecordNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	sub := chal.(*models.Challenge).Submit(team, form.Answer)

	if err := j.GetSubmissionStore().Save(sub); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, sub)
}
