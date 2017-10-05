package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-sqlite3"
)

func TestGetChallenges(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()

	chals, err := GetChallenges(db)
	if err != nil {
		t.Error(err)
		return
	}

	for i, e := range expect {
		expect, _ := json.Marshal(e)
		actual, _ := json.Marshal(chals[i])
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestGetChallenge(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()

	tests := []struct {
		id     int
		expect *Challenge
	}{
		{1, expect[0]},
		{5, expect[4]},
		{100, nil},
	}

	for _, test := range tests {
		chal, _ := GetChallenge(db, test.id)
		assert.Equal(t, test.expect, chal, "")
	}
}

func TestChallengeCreate(t *testing.T) {
	db, _, _, _ := initDB()
	defer db.Close()

	tests := []struct {
		chal   *Challenge
		expect error
	}{
		{
			&Challenge{Title: "title", Point: 100, Description: "desc", Flag: "flag"},
			nil,
		},
		{
			&Challenge{ID: 1, Title: "title", Point: 100, Description: "desc", Flag: "flag"},
			sqlite3.ErrConstraintPrimaryKey,
		},
	}

	for _, test := range tests {
		err := test.chal.Create(db)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.NotEmpty(t, test.chal.ID, "")
			assert.Equal(t, test.expect, err, "")
		}
	}
}

func TestChallengeSave(t *testing.T) {
	db, _, _, _ := initDB()
	defer db.Close()

	tests := []struct {
		chal   *Challenge
		expect error
	}{
		{
			&Challenge{Title: "title", Point: 100, Description: "desc", Flag: "flag"},
			nil,
		},
		{
			&Challenge{ID: 1, Title: "title2", Point: 100, Description: "desc", Flag: "flag"},
			nil,
		},
	}

	for _, test := range tests {
		err := test.chal.Save(db)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.NotEmpty(t, test.chal.ID, "")
			assert.NotEmpty(t, test.chal.ID, "")
			assert.Equal(t, test.expect, err, "")
		}
	}
}
