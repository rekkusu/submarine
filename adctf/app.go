package adctf

import (
	"github.com/activedefense/submarine/adctf/handlers"
	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ADCTFConfig struct {
	DriverName     string
	DataSourceName string
	JWTSecret      []byte
	Debug          bool
}

func New(config ADCTFConfig) *echo.Echo {
	db, _ := gorm.Open(config.DriverName, config.DataSourceName)
	db.AutoMigrate(&models.Challenge{}, &models.Submission{}, &models.Team{})

	jeopardy := &Jeopardy{
		Challenge:  &models.ChallengeStore{db},
		Submission: &models.SubmissionStore{db},
		Team:       &models.TeamStore{db},
	}

	e := echo.New()
	e.Debug = config.Debug
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtconf := middleware.DefaultJWTConfig
	jwtconf.ContextKey = "jwt"
	jwtconf.SigningKey = config.JWTSecret

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("jeopardy", jeopardy)
			return next(c)
		}
	})

	{
		teams := e.Group("/api/v1/teams")
		teams.GET("", handlers.GetTeams)
		teams.POST("", handlers.CreateTeam)
		teams.GET("/:id", handlers.GetTeamByID)
	}

	{
		chals := e.Group("/api/v1/challenges")
		chals.Use(middleware.JWTWithConfig(jwtconf))
		chals.GET("", handlers.GetChallenges)
		chals.POST("", handlers.NewChallenge)
		chals.GET("/:id", handlers.GetChallengeByID)
		chals.POST("/:id/submit", handlers.Submit)
	}

	{
		users := e.Group("/api/v1/users")
		users.POST("/signup", handlers.Signup)
		users.POST("/signin", handlers.Signin)
	}

	{
		me := e.Group("/api/v1/me")
		me.Use(middleware.JWTWithConfig(jwtconf))
		me.GET("", handlers.Me)
	}

	return e
}

type Jeopardy struct {
	Challenge  ctf.ChallengeStore
	Team       ctf.TeamStore
	Submission ctf.SubmissionStore
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
