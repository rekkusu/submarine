package adctf

import "github.com/jinzhu/gorm"

func initDB() (*gorm.DB, []*challenge, []*submission, []*team) {
	chals := []*challenge{
		&challenge{Title: "Test1", Point: 100, Description: "Desc1", Flag: "Flag1"},
		&challenge{Title: "Test2", Point: 200, Description: "Desc2", Flag: "Flag2"},
		&challenge{Title: "Test3", Point: 300, Description: "Desc3", Flag: "Flag3"},
		&challenge{Title: "Test4", Point: 400, Description: "Desc4", Flag: "Flag4"},
		&challenge{Title: "Test5", Point: 500, Description: "Desc5", Flag: "Flag5"},
	}

	teams := []*team{
		&team{TeamName: "team1"},
		&team{TeamName: "team2"},
		&team{TeamName: "team3"},
		&team{TeamName: "team4"},
	}

	submissions := []*submission{
		&submission{Team: teams[0], Challenge: chals[0], Answer: "sample", Score: 0, Correct: false},
	}

	db, _ := gorm.Open("sqlite3", ":memory:?parseTime=true")
	db.LogMode(false)
	db.CreateTable(&challenge{}, &submission{}, &team{})

	for _, chal := range chals {
		db.Create(chal)
	}

	for i, _ := range teams {
		db.Create(&teams[i])
	}

	for i, _ := range submissions {
		db.Create(&submissions[i])
	}

	return db, chals, submissions, teams
}
