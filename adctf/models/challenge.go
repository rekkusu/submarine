package models

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type Challenge struct {
	ID          int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CategoryID  int     `json:"category_id"`
	Title       string  `json:"title"`
	Point       int     `json:"point"`
	Description string  `json:"description"`
	Flag        *string `json:"flag,omitempty"`
}

type ChallengeWithSolves struct {
	Challenge
	Solves int `json:"solves"`
}

func (c Challenge) GetID() int {
	return c.ID
}

func (c Challenge) GetTitle() string {
	return c.Title
}

func (c Challenge) GetPoint() int {
	return c.Point
}

func (c *Challenge) GetSolves() int {
	return 0
}

func (c *ChallengeWithSolves) GetSolves() int {
	return c.Solves
}

func (c Challenge) GetDescription() string {
	return c.Description
}

func (c Challenge) GetFlag() string {
	if c.Flag == nil {
		return ""
	}
	return *c.Flag
}

func (c *Challenge) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Challenge) Save(db *gorm.DB) error {
	return db.Save(c).Error
}

func (c *Challenge) Delete(db *gorm.DB) error {
	return db.Delete(c).Error
}

func (c Challenge) Submit(db *gorm.DB, team ctf.Team, answer string) (*Submission, error) {
	correct := c.GetFlag() == answer

	s := &Submission{
		Team:      team.(*Team),
		Challenge: &c,
		Answer:    answer,
		Correct:   correct,
	}

	tx := db.Begin()

	solved := !tx.Where("team_id = ? AND challenge_id = ? AND correct = 1", s.Team.ID, s.Challenge.ID).Find(&Submission{}).RecordNotFound()

	if solved {
		tx.Rollback()
		return nil, ErrChallengeHasAlreadySolved
	}

	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return s, nil
}

func GetChallenges(db *gorm.DB) (chals []Challenge, err error) {
	err = db.Find(&chals).Error
	return
}

func GetChallengesWithSolves(db *gorm.DB) (chals []ChallengeWithSolves, err error) {
	err = db.Table("challenges").Select("challenges.*, solves").Joins("INNER JOIN (SELECT challenge_id, COUNT(DISTINCT team_id) as solves FROM submissions WHERE correct=1 GROUP BY challenge_id) ON challenge_id = challenges.id").Scan(&chals).Error
	return
}

func GetChallenge(db *gorm.DB, id int) (*Challenge, error) {
	var chal Challenge
	err := db.First(&chal, id).Error
	if err != nil {
		return nil, err
	}
	return &chal, err
}
