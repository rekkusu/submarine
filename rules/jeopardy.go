package rules

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type JeopardyRule interface {
	GetChallenges(db *gorm.DB) ([]ctf.Challenge, error)
	GetSubmissions(db *gorm.DB) ([]ctf.Submission, error)
	GetTeams(db *gorm.DB) ([]ctf.Team, error)
	GetTeam(db *gorm.DB, id int) (ctf.Team, error)
	GetScoring() ctf.Scoring
}
