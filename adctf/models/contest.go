package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ContestInfo struct {
	ID     int           `json:"-"`
	Status ContestStatus `json:"status"`
}

func (ContestInfo) TableName() string {
	return "contest_info"
}

const (
	contestInfoID = 1
)

type Announce struct {
	ID       int       `json:"id"`
	Message  string    `json:"message"`
	PostedAt time.Time `json:"posted_at"`
}

//go:generate jsonenums -type=ContestStatus
type ContestStatus int

const (
	Undefined ContestStatus = iota
	ContestClosed
	ContestOpen
	ContestFinished
)

func (status ContestStatus) String() string {
	switch status {
	case ContestClosed:
		return "closed"
	case ContestOpen:
		return "open"
	case ContestFinished:
		return "finished"
	default:
		return "undefined"
	}
}

func GetContestInfo(db *gorm.DB) (*ContestInfo, error) {
	var info ContestInfo
	err := db.Where(ContestInfo{ID: contestInfoID}).Attrs(ContestInfo{Status: ContestClosed}).FirstOrCreate(&info).Error
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func SetContestStatus(db *gorm.DB, info ContestInfo) error {
	info.ID = contestInfoID
	return db.Save(&info).Error
}
