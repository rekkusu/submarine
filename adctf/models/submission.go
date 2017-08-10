package models

import (
	"errors"
	"time"

	"github.com/activedefense/submarine/ctf"

	"github.com/jinzhu/gorm"
)

var (
	ErrChallengeHasAlreadySolved = errors.New("the challenge has already solved")
)

type Submission struct {
	ID          int        `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Team        *Team      `json:"team,omitempty" gorm:"ForeignKey:TeamID"`
	TeamID      int        `json:"-"`
	Challenge   *Challenge `json:"challenge,omitempty" gorm:"ForeignKey:ChallengeID"`
	ChallengeID int        `json:"-"`
	Answer      string     `json:"answer"`
	Score       int        `json:"score"`
	Correct     bool       `json:"is_correct"`
	CreatedAt   time.Time  `json:"submitted_at"`
}

func (s Submission) GetID() int {
	return s.ID
}

func (s Submission) GetTeam() ctf.Team {
	return s.Team
}

func (s Submission) GetUser() ctf.User {
	return s.Team
}

func (s Submission) GetChallenge() ctf.Challenge {
	return s.Challenge
}

func (s Submission) GetAnswer() string {
	return s.Answer
}

func (s Submission) GetScore() int {
	return s.Score
}

func (s Submission) IsCorrect() bool {
	return s.Correct
}

func (s Submission) GetSubmittedAt() time.Time {
	return s.CreatedAt
}

type SubmissionStore struct {
	DB *gorm.DB
}

func (store *SubmissionStore) All() ([]ctf.Submission, error) {
	var submissions []Submission
	if err := store.DB.Preload("Team").Preload("Challenge").Find(&submissions).Error; err != nil {
		return nil, err
	}

	result := make([]ctf.Submission, len(submissions))
	for i, _ := range submissions {
		result[i] = &submissions[i]
	}

	return result, nil
}

func (store *SubmissionStore) Get(id int) (ctf.Submission, error) {
	var sub Submission
	if err := store.DB.Preload("Team").Preload("Challenge").First(&sub, id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (store *SubmissionStore) Save(s ctf.Submission) error {
	sub, ok := s.(*Submission)
	if !ok {
		return ctf.ErrModelMismatched
	}

	tx := store.DB.Begin()

	// check if the chal is already solved
	solved := !tx.Where("team_id = ? AND challenge_id = ? AND correct = 1", sub.Team.GetID(), sub.Challenge.GetID()).Find(&Submission{}).RecordNotFound()

	if solved {
		tx.Commit()
		return ErrChallengeHasAlreadySolved
	}

	if err := tx.Create(&sub).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
