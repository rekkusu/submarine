package handlers

import (
	"github.com/activedefense/submarine/adctf/scoring"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
	//Jeopardy rules.JeopardyRule
	Scoring scoring.ScoringRule
}
