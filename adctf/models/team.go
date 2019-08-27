package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name" gorm:"unique_index"`
	Password string `json:"-"`
}

func (t Team) GetID() int {
	return t.ID
}

func (t Team) GetName() string {
	return t.Name
}

func (t *Team) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *Team) Save(db *gorm.DB) error {
	return db.Save(t).Error
}

func CreateTeam(db *gorm.DB, name, password string) (*Team, error) {
	var passhash []byte
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	t := &Team{
		ID:       0,
		Name:     name,
		Password: string(passhash),
	}

	if err := t.Create(db); err != nil {
		return nil, err
	}

	return t, nil
}

func JoinTeam(db *gorm.DB, userID int, name, password string) error {
	team, err := GetTeamFromName(db, name)
	if err != nil {
		return errors.New("the team is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(team.Password), []byte(password))
	if err != nil {
		return errors.New("the password is wrong")
	}

	user, err := GetUser(db, userID)
	if err != nil {
		return errors.New("the user is not found")
	}

	user.TeamID = team.ID
	user.Save(db)

	return nil
}

func GetTeams(db *gorm.DB) (teams []Team, err error) {
	err = db.Find(&teams).Error
	return
}

func GetTeam(db *gorm.DB, id int) (*Team, error) {
	var t Team
	if err := db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTeamFromName(db *gorm.DB, name string) (*Team, error) {
	var t Team
	if err := db.Where("name = ?", name).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
