package handlers

import (
	"fmt"
	"github.com/activedefense/submarine/adctf/config"
	"net/http"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

func (h *Handler) Signup(c echo.Context) error {
	var form struct {
		Username  string `json:"Username" validate:"required,min=4,max=32"`
		Password  string `json:"Password" validate:"required,eqfield=Password2"`
		Password2 string `json:"Password2" validate:"required"`
		Team      struct {
			Mode         string `validate:"required"`
			TeamName     string `json:"name" validate:"required,min=2"`
			TeamPassword string `json:"Password" validate:"required,min=8"`
		} `json:"Team"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Printf("%s\n", form.Team.Mode)
	passhash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tx := h.DB.Begin()
	if !config.CTF.Team {
		form.Team.Mode = "create"
		form.Team.TeamName = form.Username
		form.Team.TeamPassword = form.Password
	}

	if form.Team.Mode == "create" {
		_, err := models.GetTeamFromName(tx, form.Team.TeamName)
		if err == nil { // Found existing team
			return echo.NewHTTPError(http.StatusConflict, "duplicate team name")
		}

		_, err = models.CreateTeam(tx, form.Team.TeamName, form.Team.TeamPassword)
		if err != nil {
			return err
		}
	}

	user := &models.User{
		Username: form.Username,
		Password: string(passhash),
		Role:     "normal",
	}

	if _, err := models.GetUserByName(tx, user.GetName()); err != gorm.ErrRecordNotFound {
		tx.Rollback()
		if err == nil {
			return echo.NewHTTPError(http.StatusConflict, "duplicate")
		}
		return err
	}

	if err := user.Create(tx); err != nil {
		tx.Rollback()
		return err
	}

	err = models.JoinTeam(tx, user.ID, form.Team.TeamName, form.Team.TeamPassword)
	if err != nil {
		tx.Rollback()
		return err
		return echo.NewHTTPError(http.StatusUnauthorized, "Team error")
	}

	tx.Commit()

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Signin(c echo.Context) error {
	var form struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	team, err := models.GetUserByName(h.DB, form.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	hash := []byte(team.GetPassword())
	if err := bcrypt.CompareHashAndPassword(hash, []byte(form.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = team.GetID()
	claims["role"] = team.Role

	key := c.Get("secret").([]byte)

	t, err := token.SignedString(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (h *Handler) Me(c echo.Context) error {
	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	id, ok := claims["user"].(float64)
	if !ok {
		return echo.ErrNotFound
	}

	team, err := models.GetUser(h.DB, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, team)
}

func (h *Handler) SetPrivilege(c echo.Context) error {
	var form struct {
		Password string `json:"Password"`
	}
	if err := c.Bind(&form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if form.Password != c.Get("Password").(string) {
		return echo.ErrForbidden
	}

	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	id, ok := claims["user"].(float64)
	if !ok {
		return echo.ErrNotFound
	}

	team, err := models.GetUser(h.DB, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	team.Role = "admin"
	team.Save(h.DB)

	return c.NoContent(http.StatusNoContent)
}
