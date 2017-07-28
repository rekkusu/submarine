package models

import (
	"encoding/json"
	"testing"

	"github.com/activedefense/submarine/ctf"
	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-sqlite3"
)

func TestChallengeStoreAll(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()
	repo := ChallengeStore{db}

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

func TestChallengeStoreGet(t *testing.T) {
	db, expect, _, _ := initDB()
	defer db.Close()
	repo := ChallengeStore{db}

	tests := []struct {
		id     int
		expect ctf.Challenge
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

func TestChallengeStoreSave(t *testing.T) {
	db, _, _, _ := initDB()
	defer db.Close()
	repo := ChallengeStore{db}

	tests := []struct {
		chal   ctf.Challenge
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
		{
			&fakeChal{},
			ctf.ErrModelMismatched,
		},
	}

	for _, test := range tests {
		err := repo.Save(test.chal)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			if chal, ok := test.chal.(*Challenge); ok {
				assert.NotEmpty(t, chal.ID, "")
			}
			assert.Equal(t, test.expect, err, "")
		}
	}
}

type fakeChal struct {
	Challenge
}
