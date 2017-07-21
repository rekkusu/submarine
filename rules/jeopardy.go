package rules

import "github.com/activedefense/submarine/models"

type JeopardyRule interface {
	GetChallengeRepository() models.ChallengeRepository
	GetTeamRepository() models.TeamRepository
	GetSubmissionRepository() models.SubmissionRepository
}
