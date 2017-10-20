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

func (c *Category) Delete(db *gorm.DB) error {
	tx := db.Begin()

	for _, item := range c.Challenges {
		err := item.Delete(db)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := db.Delete(c).Error; err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
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
