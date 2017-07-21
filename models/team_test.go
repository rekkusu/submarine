package models

import (
	"encoding/json"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTeamRepositoryAllTeam(t *testing.T) {
	db, _, _, expect := initDB()
	repo := DefaultTeamRepository{db}

	teams, err := repo.AllTeams()
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

func TestTeamRepositoryGetTeam(t *testing.T) {
	db, _, _, expect := initDB()
	repo := DefaultTeamRepository{db}

	tests := []struct {
		id     int
		expect Team
	}{
		{1, expect[0]},
		{100, nil},
	}

	for _, test := range tests {
		team, _ := repo.GetTeam(test.id)
		expect, _ := json.Marshal(test.expect)
		actual, _ := json.Marshal(team)
		assert.JSONEq(t, string(expect), string(actual))
	}
}

func TestTeamRepositorySaveTeam(t *testing.T) {
	db, _, _, _ := initDB()
	repo := DefaultTeamRepository{db}

	tests := []struct {
		team   Team
		expect error
	}{
		{
			&team{TeamName: "test1"},
			nil,
		},
		{
			&team{ID: 1, TeamName: "test2"},
			sqlite3.ErrConstraintPrimaryKey,
		},
		{
			&fakeTeam{},
			ErrModelMismatched,
		},
	}

	for _, test := range tests {
		err := repo.SaveTeam(test.team)

		if actual, ok := err.(sqlite3.Error); ok {
			assert.Error(t, actual)
			assert.Equal(t, test.expect, actual.ExtendedCode, "db error")
		} else {
			assert.Equal(t, test.expect, err, "")
		}
	}
}

type fakeTeam struct {
	team
}
