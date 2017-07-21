package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

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

type submission struct {
	ID          int        `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Team        *team      `json:"team" gorm:"ForeignKey:TeamID"`
	TeamID      int        `json:"-"`
	User        *user      `json:"user" gorm:"ForeignKey:UserID"`
	UserID      int        `json:"-"`
	Challenge   *challenge `json:"challenge" gorm:"ForeignKey:ChallengeID"`
	ChallengeID int        `json:"-"`
	Answer      string     `json:"answer"`
	Score       int        `json:"score"`
	Correct     bool       `json:"is_correct"`
	CreatedAt   time.Time  `json:"submitted_at"`
}

func (s submission) GetID() int {
	return s.ID
}

func (s submission) GetTeam() Team {
	return s.Team
}

func (s submission) GetUser() User {
	return s.User
}

func (s submission) GetChallenge() Challenge {
	return s.Challenge
}

func (s submission) GetAnswer() string {
	return s.Answer
}

func (s submission) GetScore() int {
	return s.Score
}

func (s submission) IsCorrect() bool {
	return s.Correct
}

func (s submission) GetSubmittedAt() time.Time {
	return s.CreatedAt
}

type SubmissionRepository interface {
	All() ([]Submission, error)
	Get(id int) (Submission, error)
	Save(s Submission) error
}

type DefaultSubmissionRepository struct {
	DB *gorm.DB
}

func (repo *DefaultSubmissionRepository) All() ([]Submission, error) {
	var submissions []submission
	if err := repo.DB.Preload("Team").Preload("Challenge").Find(&submissions).Error; err != nil {
		return nil, err
	}

	result := make([]Submission, len(submissions))
	for i, _ := range submissions {
		result[i] = &submissions[i]
	}

	return result, nil
}

func (repo *DefaultSubmissionRepository) Get(id int) (Submission, error) {
	var sub submission
	if err := repo.DB.Preload("Team").Preload("Challenge").First(&sub, id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (repo *DefaultSubmissionRepository) Save(s Submission) error {
	sub, ok := s.(*submission)
	if !ok {
		return ErrModelMismatched
	}
	return repo.DB.Create(&sub).Error
}
