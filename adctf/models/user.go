package models

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username" gorm:"unique_index"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	Attributes string `json:"attrs"`
}

func (t User) GetID() int {
	return t.ID
}

func (t User) GetName() string {
	return t.Username
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetPassword() string {
	return u.Password
}

func (u User) GetTeam() ctf.Team {
	return u
}

func (u User) GetRole() string {
	return u.Role
}

func (t *User) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *User) Save(db *gorm.DB) error {
	return db.Save(t).Error
}

func GetTeams(db *gorm.DB) (teams []User, err error) {
	err = db.Find(&teams).Error
	return
}

func GetTeam(db *gorm.DB, id int) (*User, error) {
	var t User
	if err := db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTeamByName(db *gorm.DB, name string) (*User, error) {
	var t User
	if err := db.First(&t, "username=?", name).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
