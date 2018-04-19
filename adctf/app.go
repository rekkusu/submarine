package adctf

import (
	"math"
	"reflect"
	"strings"

	"github.com/activedefense/submarine/adctf/handlers"
	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
	"github.com/activedefense/submarine/scoring"
	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
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
	JWTKey        = "jwt"
	NotAuthorized = "noauth"
)

func New(config ADCTFConfig) *echo.Echo {
	db, _ := gorm.Open(config.DriverName, config.DataSourceName)
	db.AutoMigrate(&models.Challenge{}, &models.Submission{}, &models.Team{}, &models.Category{}, &models.ContestInfo{}, &models.Announcement{})

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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("secret", config.JWTSecret)
			c.Set("jeopardy", jeopardy)
			c.Set("team", getTeamFromJWT(c))
			return next(c)
		}
	})

	e.Use(CasbinMiddleware(enforcer))

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
		teams.PATCH("", handlers.UpdateTeam)
		teams.GET("/:id", handlers.GetTeamByID)
	}

	{
		chals := e.Group("/api/v1/challenges")
		chals.GET("", handlers.GetChallenges)
		chals.POST("", handlers.CreateChallenge)
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
		submissions := e.Group("/api/v1/submissions")
		submissions.GET("/solved", handlers.GetSolvedChallenges)
	}

	{
		users := e.Group("/api/v1/users")
		users.POST("/signup", handlers.Signup)
		users.POST("/signin", handlers.Signin)
	}

	{
		me := e.Group("/api/v1/me")
		me.GET("", handlers.Me)
	}

	e.GET("/api/v1/scoreboard", handlers.GetScoreboard)
	e.GET("/api/v1/contest", handlers.GetContestInfo)
	e.PUT("/api/v1/contest", handlers.PutContestInfo)

	{
		announcements := e.Group("/api/v1/announcements")
		announcements.GET("", handlers.GetCurrentAnnouncements)
		announcements.POST("", handlers.NewAnnouncement)
		announcements.GET("/:id", handlers.GetAnnouncement)
		announcements.PUT("/:id", handlers.EditAnnouncement)
	}

	return e
}

func getTeamFromJWT(c echo.Context) *models.Team {
	if c.Get(JWTKey) == nil {
		return nil
	}

	token := c.Get(JWTKey).(*jwt.Token)
	if !token.Valid {
		return nil
	}

	user, ok := token.Claims.(jwt.MapClaims)["user"]
	if !ok {
		return nil
	}

	jeopardy := c.Get("jeopardy").(rules.JeopardyRule)
	db := jeopardy.GetDB()
	team, err := models.GetTeam(db, (int)(user.(float64)))
	if err != nil {
		return nil
	}

	return team
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

func (j Jeopardy) GetChallenges() ([]ctf.Challenge, error) {
	chals, err := models.GetChallengesWithSolves(j.DB)
	if err != nil {
		return nil, err
	}
	ret := make([]ctf.Challenge, len(chals))
	for i, _ := range chals {
		ret[i] = &chals[i]
	}
	return ret, nil
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

func (j Jeopardy) GetScoring() ctf.Scoring {
	return &scoring.DynamicJeopardy{
		Jeopardy: j,
		Expression: func(base, solves int) int {
			if solves == 0 {
				return base
			}
			return int(math.Max(float64(base/100), float64(base)/math.Cbrt(float64(solves))))
		},
	}
}

func initEnforcer(config ADCTFConfig) *casbin.Enforcer {
	enforcer := casbin.NewEnforcer("adctf/policy.conf", "adctf/policy.csv")

	return enforcer
}
