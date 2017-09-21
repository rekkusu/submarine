package models

import "github.com/jinzhu/gorm"

type Category struct {
	ID         int         `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name       string      `json:"name" validate:"required"`
	Challenges []Challenge `gorm:"ForeignKey:CategoryID" json:"challenges" json:",omitempty"`
}

func (c *Category) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Category) Save(db *gorm.DB) error {
	return db.Save(c).Error
}

func GetCategories(db *gorm.DB) (categories []Category, err error) {
	err = db.Preload("Challenges").Find(&categories).Error
	return
}

func GetCategoryByID(db *gorm.DB, id int) (category *Category, err error) {
	err = db.Preload("Challenges").First(category, id).Error
	return
}
