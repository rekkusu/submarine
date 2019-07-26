package models

import (
	"github.com/jinzhu/gorm"
)

func IsAdmin(team *User) bool {
	return team != nil && team.Role == "admin"
}

func IsContestOpen(db *gorm.DB) bool {
	info, err := GetContestInfo(db)
	if err != nil {
		panic(err)
	}

	return info.Status == ContestOpen
}

func IsContestClosed(db *gorm.DB) bool {
	info, err := GetContestInfo(db)
	if err != nil {
		panic(err)
	}

	return info.Status == ContestClosed
}
