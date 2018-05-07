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

type Announcement struct {
	ID       int       `json:"id"`
	Title    string    `json:"title" validate:"required"`
	Content  string    `json:"content"`
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

func (a *Announcement) Create(db *gorm.DB) error {
	err := db.Create(a).Error
	return err
}

func (a *Announcement) Save(db *gorm.DB) error {
	err := db.Save(a).Error
	return err
}

func GetAllAnnouncements(db *gorm.DB) ([]Announcement, error) {
	var announcements []Announcement
	if err := db.Order("posted_at").Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}

func GetCurrentAnnouncements(db *gorm.DB) ([]Announcement, error) {
	var announcements []Announcement
	if err := db.Where("posted_at <= ?", time.Now()).Order("posted_at").Find(&announcements).Error; err != nil {
		return nil, err
	}
	return announcements, nil
}

func GetAnnouncement(db *gorm.DB, id int) (*Announcement, error) {
	var announcement Announcement
	if err := db.Find(&announcement, id).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}
