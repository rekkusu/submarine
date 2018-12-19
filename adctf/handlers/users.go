package handlers

import (
	"net/http"

	"github.com/activedefense/submarine/adctf/models"
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
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	tx := db.Begin()
	if _, err := models.GetTeamByName(tx, team.GetName()); err != gorm.ErrRecordNotFound {
		tx.Rollback()
		if err == nil {
			return echo.NewHTTPError(http.StatusConflict, "duplicate")
		}
		return err
	}

	if err := team.Create(tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return c.JSON(http.StatusCreated, team)
}

func Signin(c echo.Context) error {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	team, err := models.GetTeamByName(db, form.Username)
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

func Me(c echo.Context) error {
	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	id, ok := claims["user"].(float64)
	if !ok {
		return echo.ErrNotFound
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	team, err := models.GetTeam(db, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, team)
}

func SetPrivilege(c echo.Context) error {
	var form struct {
		Password string `json:"password"`
	}
	if err := c.Bind(&form); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if form.Password != c.Get("password").(string) {
		return echo.ErrForbidden
	}

	claims := c.Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	id, ok := claims["user"].(float64)
	if !ok {
		return echo.ErrNotFound
	}

	db := c.Get("jeopardy").(rules.JeopardyRule).GetDB()
	team, err := models.GetTeam(db, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	team.Role = "admin"
	team.Save(db)

	return c.NoContent(http.StatusNoContent)
}
