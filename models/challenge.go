package models

type Challenge interface {
	GetID() int
	GetTitle() string
	GetPoint() int
	GetDescription() string
	GetFlag() string
}

type ChallengeStore interface {
	All() ([]Challenge, error)
	Get(id int) (Challenge, error)
	Save(c Challenge) error
}
