package rules

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type JeopardyRule interface {
	GetDB() *gorm.DB
	GetChallengeStore() ctf.ChallengeStore
	GetTeamStore() ctf.TeamStore
	GetSubmissionStore() ctf.SubmissionStore
}
