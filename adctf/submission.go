package adctf

import (
	"time"

	"github.com/activedefense/submarine/models"

	"github.com/jinzhu/gorm"
)

type Submission struct {
	ID          int        `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Team        *Team      `json:"team" gorm:"ForeignKey:TeamID"`
	TeamID      int        `json:"-"`
	User        *User      `json:"user" gorm:"ForeignKey:UserID"`
	UserID      int        `json:"-"`
	Challenge   *Challenge `json:"challenge" gorm:"ForeignKey:ChallengeID"`
	ChallengeID int        `json:"-"`
	Answer      string     `json:"answer"`
	Score       int        `json:"score"`
	Correct     bool       `json:"is_correct"`
	CreatedAt   time.Time  `json:"submitted_at"`
}

func (s Submission) GetID() int {
	return s.ID
}

func (s Submission) GetTeam() models.Team {
	return s.Team
}

func (s Submission) GetUser() models.User {
	return s.User
}

func (s Submission) GetChallenge() models.Challenge {
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

func (store *SubmissionStore) All() ([]models.Submission, error) {
	var submissions []Submission
	if err := store.DB.Preload("Team").Preload("Challenge").Find(&submissions).Error; err != nil {
		return nil, err
	}

	result := make([]models.Submission, len(submissions))
	for i, _ := range submissions {
		result[i] = &submissions[i]
	}

	return result, nil
}

func (store *SubmissionStore) Get(id int) (models.Submission, error) {
	var sub Submission
	if err := store.DB.Preload("Team").Preload("Challenge").First(&sub, id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (store *SubmissionStore) Save(s models.Submission) error {
	sub, ok := s.(*Submission)
	if !ok {
		return models.ErrModelMismatched
	}
	return store.DB.Create(&sub).Error
}
