package handlers

import (
	"net/http"

	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

func Signup(c echo.Context) error {
	var form struct {
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required,eqfield=Password2"`
		Password2 string `json:"password2" validate:"required"`
	}

	if err := c.Bind(&form); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return err
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	team := &models.Team{
		Username: form.Username,
		Password: string(passhash),
		Role:     "normal",
	}

	store, _ := c.Get("jeopardy").(rules.JeopardyRule).GetTeamStore().(*models.TeamStore)
	if _, err := store.GetTeamByName(team.GetName()); err != gorm.ErrRecordNotFound {
		if err == nil {
			return echo.NewHTTPError(http.StatusConflict, "duplicate")
		}
		return err
	}

	if err := store.SaveTeam(team); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, team)
}

func Signin(c echo.Context) error {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&form); err != nil {
		return err
	}

	store, _ := c.Get("jeopardy").(rules.JeopardyRule).GetTeamStore().(*models.TeamStore)
	team, err := store.GetTeamByName(form.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	hash := []byte(team.(ctf.User).GetPassword())
	if err := bcrypt.CompareHashAndPassword(hash, []byte(form.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = team.(ctf.User).GetID()
	claims["role"] = team.(*models.Team).Role

	key := c.Get("secret").([]byte)

	t, err := token.SignedString(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Me(c echo.Context) error {
	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	id, ok := claims["user"].(float64)
	if !ok {
		return echo.ErrNotFound
	}

	store := c.Get("jeopardy").(rules.JeopardyRule).GetTeamStore()
	team, err := store.GetTeam(int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, team)
}
