package handlers

import (
	"github.com/activedefense/submarine/ctf"
	"net/http"
	"strconv"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetTeams(c echo.Context) error {
	teams, err := models.GetTeams(h.DB)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, teams)
}

func (h *Handler) GetTeamByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	team, err := models.GetTeam(h.DB, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	solved, err := models.GetSolvedChallenges(h.DB, team.GetID())
	solves, err := models.GetSolves(h.DB)

	for _, chal := range solved {
		var count int
		for _, sol := range solves {
			if sol.ChallengeID == chal.Challenge.ID {
				count = sol.Solves
			}
		}
		chal.Challenge.Point = h.Scoring.CalcScore(models.ChallengeWithSolves{
			Challenge: *chal.Challenge,
			Solves:    count,
		})
	}

	return c.JSON(http.StatusOK, struct {
		ctf.Team
		Solved []models.Submission `json:"solved"`
	}{team, solved})
}

func (h *Handler) CreateTeam(c echo.Context) error {
	team := &models.Team{}
	if err := c.Bind(team); err != nil {
		return err
	}

	team.ID = 0

	if err := team.Create(h.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, team)
}

func (h *Handler) UpdateTeam(c echo.Context) error {
	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	team_id := int(claims["user"].(float64))
	team, err := models.GetUser(h.DB, team_id)

	if err != nil {
		return echo.ErrNotFound
	}

	var form struct {
		Attributes string `json:"attrs"`
	}

	if err := c.Bind(&form); err != nil {
		return err
	}

	team.Attributes = form.Attributes

	if err := team.Save(h.DB); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
