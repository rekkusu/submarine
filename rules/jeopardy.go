package rules

import "github.com/activedefense/submarine/ctf"

type JeopardyRule interface {
	GetChallengeStore() ctf.ChallengeStore
	GetTeamStore() ctf.TeamStore
	GetSubmissionStore() ctf.SubmissionStore
}
