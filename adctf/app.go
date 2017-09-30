package adctf

import (
	"reflect"
	"strings"

	"github.com/activedefense/submarine/adctf/handlers"
	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/ctf"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

type ADCTFConfig struct {
	DriverName     string
	DataSourceName string
	JWTSecret      []byte
	Debug          bool
}

const (
	JWTKey = "jwt"
)

func New(config ADCTFConfig) *echo.Echo {
	db, _ := gorm.Open(config.DriverName, config.DataSourceName)
	db.AutoMigrate(&models.Challenge{}, &models.Submission{}, &models.Team{}, &models.Category{})

	enforcer := initEnforcer(config)

	jeopardy := &Jeopardy{
		DB: db,
	}

	e := echo.New()
	e.Debug = config.Debug
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey: JWTKey,
		SigningKey: config.JWTSecret,
		Skipper: func(c echo.Context) bool {
			if c.Request().Header.Get("Authorization") == "" {
				return true
			}
			return false
		},
	}))
	e.Use(CasbinMiddleware(enforcer))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("secret", config.JWTSecret)
			c.Set("jeopardy", jeopardy)
			return next(c)
		}
	})

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	e.Validator = &CustomValidator{validate}

	{
		teams := e.Group("/api/v1/teams")
		teams.GET("", handlers.GetTeams)
		teams.POST("", handlers.CreateTeam)
		teams.GET("/:id", handlers.GetTeamByID)
	}

	{
		chals := e.Group("/api/v1/challenges")
		chals.GET("", handlers.GetChallenges)
		chals.POST("", handlers.CreateChallenge)
		chals.GET("/solves", handlers.GetSolves)
		chals.GET("/:id", handlers.GetChallengeByID)
		chals.PUT("/:id", handlers.UpdateChallenge)
		chals.DELETE("/:id", handlers.DeleteChallenge)
		chals.POST("/:id/submit", handlers.Submit)
	}

	{
		cate := e.Group("/api/v1/categories")
		cate.GET("", handlers.GetCategories)
		cate.POST("", handlers.CreateCategory)
		cate.PUT("/:id", handlers.UpdateCategory)
		cate.DELETE("/:id", handlers.DeleteCategory)
	}

	{
		users := e.Group("/api/v1/users")
		users.POST("/signup", handlers.Signup)
		users.POST("/signin", handlers.Signin)
	}

	e.GET("/api/v1/scoreboard", handlers.GetScoreboard)
	{
		me := e.Group("/api/v1/me")
		me.GET("", handlers.Me)
	}

	return e
}

type CustomValidator struct {
	validate *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

type Jeopardy struct {
	DB *gorm.DB
}

func (j Jeopardy) GetDB() *gorm.DB {
	return j.DB
}

func (j Jeopardy) GetSubmissions() ([]ctf.Submission, error) {
	sub, err := models.GetSubmissions(j.DB)
	if err != nil {
		return nil, err
	}
	ret := make([]ctf.Submission, len(sub))
	for i, _ := range sub {
		ret[i] = &sub[i]
	}
	return ret, nil
}

func (j Jeopardy) GetTeams() ([]ctf.Team, error) {
	teams, err := models.GetTeams(j.DB)
	if err != nil {
		return nil, err
	}
	ret := make([]ctf.Team, len(teams))
	for i, _ := range teams {
		ret[i] = &teams[i]
	}
	return ret, nil
}

func (j Jeopardy) GetTeam(id int) (ctf.Team, error) {
	return models.GetTeam(j.DB, id)
}

func initEnforcer(config ADCTFConfig) *casbin.Enforcer {
	enforcer := casbin.NewEnforcer("adctf/policy.conf", "adctf/policy.csv")

	return enforcer
}
