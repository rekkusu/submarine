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
	Answer      *string    `json:"answer,omitempty"`
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
	return *s.Answer
}

func (s Submission) IsCorrect() bool {
	return s.Correct
}

func (s Submission) GetSubmittedAt() time.Time {
	return s.CreatedAt
}

func (s *Submission) Create(db *gorm.DB) error {
	tx := db.Begin()

	solved := !tx.Where("team_id = ? AND challenge_id = ? AND correct = 1", s.Team.ID, s.Challenge.ID).Find(&Submission{}).RecordNotFound()

	if solved {
		tx.Rollback()
		return ErrChallengeHasAlreadySolved
	}

	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GetSubmissions(db *gorm.DB) (submissions []Submission, err error) {
	err = db.Preload("Team").Preload("Challenge").Find(&submissions).Error
	return
}

type Solves struct {
	ChallengeID int `json:"challenge_id"`
	Solves      int `json:"solves"`
}

func GetSolves(db *gorm.DB) ([]Solves, error) {
	var solves []Solves
	err := db.Select("challenge_id, COUNT(DISTINCT team_id) as solves").Where("correct=?", true).Group("challenge_id").Table("submissions").Find(&solves).Error
	return solves, err
}

func GetSolvedChallenges(db *gorm.DB, id int) ([]Submission, error) {
	var subs []Submission
	err := db.Preload("Challenge").Where("correct = 1 AND team_id = ?", id).Find(&subs).Error
	if err != nil {
		return nil, err
	}
	for i, _ := range subs {
		subs[i].Answer = nil
		subs[i].Challenge.Flag = nil
	}
	return subs, nil
}
