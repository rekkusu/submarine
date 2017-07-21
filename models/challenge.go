package models

import "github.com/jinzhu/gorm"

type Challenge interface {
	GetID() int
	GetTitle() string
	GetPoint() int
	GetDescription() string
	GetFlag() string
}

type challenge struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string `json:"title"`
	Point       int    `json:"point"`
	Description string `json:"description"`
	Flag        string `json:"-"`
}

func NewChallenge(id int, title string, point int, description string, flag string) Challenge {
	return challenge{id, title, point, description, flag}
}

func (c challenge) TableName() string {
	return "challenges"
}

func (c challenge) GetID() int {
	return c.ID
}

func (c challenge) GetTitle() string {
	return c.Title
}

func (c challenge) GetPoint() int {
	return c.Point
}

func (c challenge) GetDescription() string {
	return c.Description
}

func (c challenge) GetFlag() string {
	return c.Flag
}

type ChallengeRepository interface {
	All() ([]Challenge, error)
	Get(id int) (Challenge, error)
	Save(c Challenge) error
}

type DefaultChallengeRepository struct {
	DB *gorm.DB
}

func (repo *DefaultChallengeRepository) All() ([]Challenge, error) {
	var chals []challenge
	if err := repo.DB.Find(&chals).Error; err != nil {
		return nil, err
	}

	result := make([]Challenge, len(chals))
	for i, _ := range chals {
		result[i] = &chals[i]
	}

	return result, nil
}

func (repo *DefaultChallengeRepository) Get(id int) (Challenge, error) {
	var chal challenge
	if err := repo.DB.First(&chal, id).Error; err != nil {
		return nil, err
	}
	return &chal, nil
}

func (repo *DefaultChallengeRepository) Save(c Challenge) error {
	chal, ok := c.(*challenge)
	if !ok {
		return ErrModelMismatched
	}
	return repo.DB.Create(&chal).Error
}
