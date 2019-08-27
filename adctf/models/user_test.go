package models

import (
	"encoding/json"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetTeams(t *testing.T) {
	db, _, _, expect := initDB()

	teams, err := GetUsers(db)
	if err != nil {
		t.Error(err)
		return
	}

	for i, e := range expect {
		expect, _ := json.Marshal(e)
		actual, _ := json.Marshal(teams[i])
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestGetTeam(t *testing.T) {
	db, _, _, expect := initDB()

	tests := []struct {
		id     int
		expect *User
	}{
		{1, expect[0]},
		{100, nil},
	}

	for _, test := range tests {
		team, _ := GetUser(db, test.id)
		expect, _ := json.Marshal(test.expect)
		actual, _ := json.Marshal(team)
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestTeamCreate(t *testing.T) {
	db, _, _, _ := initDB()

	tests := []struct {
		team   *User
		expect error
	}{
		{
			&User{Username: "test1"},
			nil,
		},
		{
			&User{ID: 1, Username: "test2"},
			sqlite3.ErrConstraintPrimaryKey,
		},
	}

	for _, test := range tests {
		err := test.team.Create(db)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.Equal(t, test.expect, err, "")
		}
	}
}
