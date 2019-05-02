package ctf

import "time"

type Submission interface {
	GetID() int
	GetTeamID() int
	GetTeam() Team
	GetUser() User
	GetChallengeID() int
	GetChallenge() Challenge
	GetAnswer() string
	IsCorrect() bool
	GetSubmittedAt() time.Time
}

type SubmissionStore interface {
	All() ([]Submission, error)
	Get(id int) (Submission, error)
	Save(s Submission) error
}
