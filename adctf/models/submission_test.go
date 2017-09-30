package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-sqlite3"
)

func TestGetSubmissions(t *testing.T) {
	db, _, expect, _ := initDB()

	submissions, err := GetSubmissions(db)
	if err != nil {
		t.Error(err)
		return
	}

	for i, e := range expect {
		expect, _ := json.Marshal(e)
		actual, _ := json.Marshal(submissions[i])
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestGetSubmission(t *testing.T) {
	_, _, expect, _ := initDB()

	_ = []struct {
		id     int
		expect *Submission
	}{
		{1, expect[0]},
		{100, nil},
	}

	/*
		for _, test := range tests {
			sub, _ := repo.Get(test.id)
			expect, _ := json.Marshal(test.expect)
			actual, _ := json.Marshal(sub)
			assert.JSONEq(t, string(expect), string(actual))
		}
	*/
}

func TestSubmissionCreate(t *testing.T) {
	db, chals, _, teams := initDB()

	tests := []struct {
		sub    *Submission
		expect error
	}{
		{
			&Submission{Challenge: chals[0], Answer: "test", Score: 100, Correct: false, Team: teams[0]},
			nil,
		},
		{
			&Submission{ID: 1, Challenge: chals[1], Answer: "test", Score: 100, Correct: false, Team: teams[1]},
			sqlite3.ErrConstraintPrimaryKey,
		},
	}

	for _, test := range tests {
		err := test.sub.Create(db)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.Equal(t, test.expect, err, "")
		}
	}
}

func TestGetSolves(t *testing.T) {
	db, chals, _, teams := initDB()
	subs := []Submission{
		Submission{Challenge: chals[0], Answer: "sample", Score: 100, Correct: true, Team: teams[0]},
		Submission{Challenge: chals[1], Answer: "sample", Score: 100, Correct: true, Team: teams[0]},
		Submission{Challenge: chals[1], Answer: "sample", Score: 100, Correct: true, Team: teams[1]},
		Submission{Challenge: chals[1], Answer: "sample", Score: 100, Correct: true, Team: teams[1]},
	}

	expect := []Solves{
		Solves{1, 1},
		Solves{2, 2},
	}

	for _, item := range subs {
		db.Create(&item)
	}

	solves, err := GetSolves(db)
	assert.NoError(t, err)
	assert.Equal(t, expect, solves)
}
