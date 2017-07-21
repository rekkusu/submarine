package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattn/go-sqlite3"
)

func TestSubmissionRepositoryAll(t *testing.T) {
	db, _, expect, _ := initDB()
	repo := DefaultSubmissionRepository{db}

	submissions, err := repo.All()
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

func TestSubmissionRepositoryGet(t *testing.T) {
	db, _, expect, _ := initDB()
	repo := DefaultSubmissionRepository{db}

	tests := []struct {
		id     int
		expect Submission
	}{
		{1, expect[0]},
		{100, nil},
	}

	for _, test := range tests {
		sub, _ := repo.Get(test.id)
		expect, _ := json.Marshal(test.expect)
		actual, _ := json.Marshal(sub)
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestSubmissionRepositorySave(t *testing.T) {
	db, chals, _, _ := initDB()
	repo := DefaultSubmissionRepository{db}

	tests := []struct {
		sub    Submission
		expect error
	}{
		{
			&submission{Challenge: chals[0], Answer: "test", Score: 100, Correct: false},
			nil,
		},
		{
			&submission{ID: 1, Answer: "test", Score: 100, Correct: false},
			sqlite3.ErrConstraintPrimaryKey,
		},
		{
			&fakeSubmission{},
			ErrModelMismatched,
		},
	}

	for _, test := range tests {
		err := repo.Save(test.sub)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.Equal(t, test.expect, err, "")
		}
	}
}

type fakeSubmission struct {
	submission
}
