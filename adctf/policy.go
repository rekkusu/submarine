package adctf

import (
	"github.com/activedefense/submarine/adctf/models"
	"github.com/casbin/casbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CasbinConfig struct {
	Skipper  middleware.Skipper
	Enforcer *casbin.Enforcer
}

var DefaultCasbinConfig = CasbinConfig{
	Skipper: middleware.DefaultSkipper,
}

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
	t, ok := c.Get("user").(*models.User)
	if !ok || t == nil {
		return NotAuthorized
	}
	return t.Role
}

func (conf *CasbinConfig) CheckPermission(c echo.Context) bool {
	role := conf.GetRole(c)
	method := c.Request().Method
	path := c.Request().URL.Path
	return conf.Enforcer.Enforce(role, path, method)
}
