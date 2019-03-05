package adctf

import (
	"reflect"
	"strings"

	"github.com/activedefense/submarine/adctf/handlers"
	"github.com/activedefense/submarine/adctf/models"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type ADCTFConfig struct {
	DriverName     string
	DataSourceName string
	JWTSecret      []byte
	MasterPassword string
	Debug          bool
}

const (
	JWTKey        = "jwt"
	NotAuthorized = "noauth"
)

func New(config ADCTFConfig) *echo.Echo {
	db, err := gorm.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Challenge{}, &models.Submission{}, &models.Team{}, &models.Category{}, &models.ContestInfo{}, &models.Announcement{}).Error

	if err != nil {
		panic(err)
	}

	enforcer := initEnforcer(config)

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
			c.Set("password", config.MasterPassword)
			c.Set("team", getTeamFromJWT(db, c))
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

	handler := handlers.Handler{db, &Jeopardy{}}

	{
		teams := e.Group("/api/v1/teams")
		teams.GET("", handler.GetTeams)
		teams.POST("", handler.CreateTeam)
		teams.PATCH("", handler.UpdateTeam)
		teams.GET("/:id", handler.GetTeamByID)
	}

	{
		chals := e.Group("/api/v1/challenges")
		chals.GET("", handler.GetChallenges)
		chals.POST("", handler.CreateChallenge)
		chals.GET("/solved", handler.GetSolvedChallenges)
		chals.GET("/:id", handler.GetChallengeByID)
		chals.PUT("/:id", handler.UpdateChallenge)
		chals.DELETE("/:id", handler.DeleteChallenge)
		chals.POST("/:id/submit", handler.Submit)
	}

	{
		cate := e.Group("/api/v1/categories")
		cate.GET("", handler.GetCategories)
		cate.POST("", handler.CreateCategory)
		cate.PUT("/:id", handler.UpdateCategory)
		cate.DELETE("/:id", handler.DeleteCategory)
	}

	{
		submissions := e.Group("/api/v1/submissions")
		submissions.GET("", handler.GetSubmissions)
	}

	{
		users := e.Group("/api/v1/users")
		users.POST("/signup", handler.Signup)
		users.POST("/signin", handler.Signin)
		users.PATCH("/priv", handler.SetPrivilege)
	}

	{
		me := e.Group("/api/v1/me")
		me.GET("", handler.Me)
	}

	e.GET("/api/v1/scoreboard", handler.GetScoreboard)
	e.GET("/api/v1/contest", handler.GetContestInfo)
	e.PUT("/api/v1/contest", handler.PutContestInfo)

	{
		announcements := e.Group("/api/v1/announcements")
		announcements.GET("", handler.GetCurrentAnnouncements)
		announcements.POST("", handler.NewAnnouncement)
		announcements.GET("/:id", handler.GetAnnouncement)
		announcements.PUT("/:id", handler.EditAnnouncement)
		announcements.DELETE("/:id", handler.DeleteAnnouncement)
		announcements.GET("/all", handler.GetAllAnnouncements)
	}

	return e
}

func getTeamFromJWT(db *gorm.DB, c echo.Context) *models.Team {
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

func initEnforcer(config ADCTFConfig) *casbin.Enforcer {
	enforcer := casbin.NewEnforcer("adctf/policy.conf", "adctf/policy.csv")

	return enforcer
}
