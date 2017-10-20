package models

import "github.com/jinzhu/gorm"

type Category struct {
	ID         int         `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name       string      `json:"name" validate:"required"`
	Challenges []Challenge `gorm:"ForeignKey:CategoryID" json:"challenges,omitempty"`
}

func (c *Category) Create(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *Category) Save(db *gorm.DB) error {
	return db.Save(c).Error
}

func (c *Category) Delete(db *gorm.DB) error {
	if err := db.Delete(c).Error; err != nil {
		return err
	}

	return nil
}

func GetCategories(db *gorm.DB) (categories []Category, err error) {
	err = db.Preload("Challenges").Find(&categories).Error
	return
}

func GetCategory(db *gorm.DB, id int) (*Category, error) {
	var category Category
	err := db.Preload("Challenges").First(&category, id).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}
