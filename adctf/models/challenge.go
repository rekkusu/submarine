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

func (c Challenge) TableName() string {
	return "challenges"
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

func (c Challenge) Submit(team ctf.Team, answer string) ctf.Submission {
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

type ChallengeStore struct {
	DB *gorm.DB
}

func (repo *ChallengeStore) All() ([]ctf.Challenge, error) {
	var chals []Challenge
	if err := repo.DB.Find(&chals).Error; err != nil {
		return nil, err
	}

	result := make([]ctf.Challenge, len(chals))
	for i, _ := range chals {
		result[i] = &chals[i]
	}

	return result, nil
}

func (repo *ChallengeStore) Get(id int) (ctf.Challenge, error) {
	var chal Challenge
	if err := repo.DB.First(&chal, id).Error; err != nil {
		return nil, err
	}
	return &chal, nil
}

func (repo *ChallengeStore) Save(c ctf.Challenge) error {
	chal, ok := c.(*Challenge)
	if !ok {
		return ctf.ErrModelMismatched
	}
	return repo.DB.Create(&chal).Error
}
