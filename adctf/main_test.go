package adctf

import "github.com/jinzhu/gorm"

func initDB() (*gorm.DB, []*Challenge, []*Submission, []*Team) {
	chals := []*Challenge{
		&Challenge{Title: "Test1", Point: 100, Description: "Desc1", Flag: "Flag1"},
		&Challenge{Title: "Test2", Point: 200, Description: "Desc2", Flag: "Flag2"},
		&Challenge{Title: "Test3", Point: 300, Description: "Desc3", Flag: "Flag3"},
		&Challenge{Title: "Test4", Point: 400, Description: "Desc4", Flag: "Flag4"},
		&Challenge{Title: "Test5", Point: 500, Description: "Desc5", Flag: "Flag5"},
	}

	teams := []*Team{
		&Team{TeamName: "team1"},
		&Team{TeamName: "team2"},
		&Team{TeamName: "team3"},
		&Team{TeamName: "team4"},
	}

	submissions := []*Submission{
		&Submission{Team: teams[0], Challenge: chals[0], Answer: "sample", Score: 0, Correct: false},
	}

	db, _ := gorm.Open("sqlite3", ":memory:?parseTime=true")
	db.LogMode(false)
	db.CreateTable(&Challenge{}, &Submission{}, &Team{})

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
