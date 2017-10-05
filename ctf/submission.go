package ctf

import "time"

type Submission interface {
	GetID() int
	GetTeam() Team
	GetUser() User
	GetChallenge() Challenge
	GetAnswer() string
	GetScore() int
	IsCorrect() bool
	GetSubmittedAt() time.Time
}

type SubmissionStore interface {
	All() ([]Submission, error)
	Get(id int) (Submission, error)
	Save(s Submission) error
}
