package models

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type Challenge struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CategoryID  int    `json:"category_id"`
	Title       string `json:"title"`
	Point       int    `json:"point"`
	Description string `json:"description"`
	Flag        string `json:"-"`
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

func (c Challenge) GetDescription() string {
	return c.Description
}

func (c Challenge) GetFlag() string {
	return c.Flag
}

func (c *Challenge) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Challenge) Save(db *gorm.DB) error {
	return db.Save(c).Error
}

func (c Challenge) Submit(team ctf.Team, answer string) *Submission {
	score := 0
	correct := false

	if c.GetFlag() == answer {
		score = c.GetPoint()
		correct = true
	}

	return &Submission{
		Team:      team.(*Team),
		Challenge: &c,
		Answer:    answer,
		Score:     score,
		Correct:   correct,
	}
}

func GetChallenges(db *gorm.DB) (chals []Challenge, err error) {
	err = db.Find(&chals).Error
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
