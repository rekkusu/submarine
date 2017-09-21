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
		DB:         db,
		Challenge:  &models.ChallengeStore{db},
		Submission: &models.SubmissionStore{db},
		Team:       &models.TeamStore{db},
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
		chals.POST("", handlers.NewChallenge)
		chals.GET("/:id", handlers.GetChallengeByID)
		chals.POST("/:id/submit", handlers.Submit)
	}

	{
		cate := e.Group("/api/v1/categories")
		cate.GET("", handlers.GetCategories)
		cate.POST("", handlers.CreateCategory)
		cate.PUT("/:id", handlers.UpdateCategory)
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
	DB         *gorm.DB
	Challenge  ctf.ChallengeStore
	Team       ctf.TeamStore
	Submission ctf.SubmissionStore
}

func (j Jeopardy) GetDB() *gorm.DB {
	return j.DB
}

func (j Jeopardy) GetChallengeStore() ctf.ChallengeStore {
	return j.Challenge
}

func (j Jeopardy) GetTeamStore() ctf.TeamStore {
	return j.Team
}

func (j Jeopardy) GetSubmissionStore() ctf.SubmissionStore {
	return j.Submission
}

func initEnforcer(config ADCTFConfig) *casbin.Enforcer {
	enforcer := casbin.NewEnforcer("adctf/policy.conf", "adctf/policy.csv")

	return enforcer
}
