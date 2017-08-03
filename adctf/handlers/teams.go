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
	teams, err := j.GetTeamStore().AllTeams()
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

	team, err := j.GetTeamStore().GetTeam(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, team)
}

func CreateTeam(c echo.Context) error {
	team := &models.Team{}
	if err := c.Bind(team); err != nil {
		return err
	}

	team.ID = 0

	j, _ := c.Get("jeopardy").(rules.JeopardyRule)
	if err := j.GetTeamStore().SaveTeam(team); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, team)
}
