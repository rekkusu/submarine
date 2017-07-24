package rules

import "github.com/activedefense/submarine/models"

type JeopardyRule interface {
	GetChallengeStore() models.ChallengeStore
	GetTeamStore() models.TeamStore
	GetSubmissionStore() models.SubmissionStore
}
