package adctf

import (
	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CasbinConfig struct {
	Skipper  middleware.Skipper
	Enforcer *casbin.Enforcer
}

var DefaultCasbinConfig = CasbinConfig{
	Skipper: middleware.DefaultSkipper,
}

const NotAuthorized = "noauth"

func CasbinMiddleware(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	c := DefaultCasbinConfig
	c.Enforcer = enforcer
	return CasbinMiddlewareWithConfig(c)
}

func CasbinMiddlewareWithConfig(config CasbinConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultCasbinConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) || config.CheckPermission(c) {
				return next(c)
			}

			if config.GetRole(c) == NotAuthorized {
				return echo.ErrUnauthorized
			} else {
				return echo.ErrForbidden
			}
		}
	}
}

func (conf *CasbinConfig) GetRole(c echo.Context) string {
	if c.Get(JWTKey) == nil {
		return NotAuthorized
	}

	token := c.Get(JWTKey).(*jwt.Token)
	if !token.Valid {
		return NotAuthorized
	}

	username, ok := token.Claims.(jwt.MapClaims)["role"]
	if !ok {
		return NotAuthorized
	}
	return username.(string)
}

func (conf *CasbinConfig) CheckPermission(c echo.Context) bool {
	role := conf.GetRole(c)
	method := c.Request().Method
	path := c.Request().URL.Path
	return conf.Enforcer.Enforce(role, path, method)
}
