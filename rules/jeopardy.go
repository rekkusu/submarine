package rules

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type JeopardyRule interface {
	GetDB() *gorm.DB
	GetSubmissions() ([]ctf.Submission, error)
	GetTeams() ([]ctf.Team, error)
	GetTeam(id int) (ctf.Team, error)
}
