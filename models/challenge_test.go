package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-sqlite3"
)

func TestChallengeRepositoryAll(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()
	repo := DefaultChallengeRepository{db}

	chals, err := repo.All()
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

func TestChallengeRepositoryGet(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()
	repo := DefaultChallengeRepository{db}

	tests := []struct {
		id     int
		expect Challenge
	}{
		{1, expect[0]},
		{5, expect[4]},
		{100, nil},
	}

	for _, test := range tests {
		chal, _ := repo.Get(test.id)
		assert.Equal(t, test.expect, chal, "")
	}
}

func TestChallengeRepositorySave(t *testing.T) {
	db, _, _, _ := initDB()
	defer db.Close()
	repo := DefaultChallengeRepository{db}

	tests := []struct {
		chal   Challenge
		expect error
	}{
		{
			&challenge{Title: "title", Point: 100, Description: "desc", Flag: "flag"},
			nil,
		},
		{
			&challenge{ID: 1, Title: "title", Point: 100, Description: "desc", Flag: "flag"},
			sqlite3.ErrConstraintPrimaryKey,
		},
		{
			&fakeChal{},
			ErrModelMismatched,
		},
	}

	for _, test := range tests {
		err := repo.Save(test.chal)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			if chal, ok := test.chal.(*challenge); ok {
				assert.NotEmpty(t, chal.ID, "")
			}
			assert.Equal(t, test.expect, err, "")
		}
	}
}

type fakeChal struct {
	challenge
}
