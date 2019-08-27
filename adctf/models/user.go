package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID         int    `json:"id"`
	TeamID     int    `json:"team_id"`
	Username   string `json:"username" gorm:"unique_index"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Attributes string `json:"attrs"`
}

func (u User) GetID() int {
	return u.ID
}

func (u User) GetTeamID() int {
	return u.TeamID
}

func (u User) GetName() string {
	return u.Username
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetPassword() string {
	return u.Password
}

func (u User) GetRole() string {
	return u.Role
}

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) Save(db *gorm.DB) error {
	return db.Save(u).Error
}

func GetUsers(db *gorm.DB) (users []*User, err error) {
	err = db.Find(&users).Error
	return
}

func GetUser(db *gorm.DB, id int) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByName(db *gorm.DB, name string) (*User, error) {
	var u User
	if err := db.First(&u, "username=?", name).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
