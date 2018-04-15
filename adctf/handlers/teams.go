package handlers

import (
	"net/http"
	"strconv"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/rules"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func GetTeams(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	teams, err := j.GetTeams()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, teams)
}

func GetTeamByID(c echo.Context) error {
	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	team, err := models.GetTeam(j.GetDB(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	solved, err := models.GetSolvedChallenges(j.GetDB(), team.GetID())
	solves, err := models.GetSolves(j.GetDB())

	for _, chal := range solved {
		var count int
		for _, sol := range solves {
			if sol.ChallengeID == chal.Challenge.ID {
				count = sol.Solves
			}
		}
		chal.Challenge.Point = j.GetScoring().CalcScore(&models.ChallengeWithSolves{
			Challenge: *chal.Challenge,
			Solves:    count,
		})
	}

	return c.JSON(http.StatusOK, struct {
		*models.Team
		Solved []models.Submission `json:"solved"`
	}{team, solved})
}

func CreateTeam(c echo.Context) error {
	team := &models.Team{}
	if err := c.Bind(team); err != nil {
		return err
	}

	team.ID = 0

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	if err := team.Create(db); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, team)
}
