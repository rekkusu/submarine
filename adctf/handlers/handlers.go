package handlers

import (
	"github.com/activedefense/submarine/rules"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB       *gorm.DB
	Jeopardy rules.JeopardyRule
}
