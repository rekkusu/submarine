package models

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type Team struct {
	ID       int    `json:"id"`
	Username string `json:"username" gorm:"unique_index"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

func (t Team) GetID() int {
	return t.ID
}

func (t Team) GetName() string {
	return t.Username
}

func (u Team) GetUsername() string {
	return u.Username
}

func (u Team) GetPassword() string {
	return u.Password
}

func (u Team) GetTeam() ctf.Team {
	return u
}

func (t *Team) Create(db *gorm.DB) error {
	return db.Create(t).Error
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

func GetTeamByName(db *gorm.DB, name string) (*Team, error) {
	var t Team
	if err := db.First(&t, "username=?", name).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
