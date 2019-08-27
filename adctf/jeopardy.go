package adctf

import (
	"github.com/activedefense/submarine/adctf/models"
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type Jeopardy struct {
}

func (j Jeopardy) GetChallenges(db *gorm.DB) ([]ctf.Challenge, error) {
	chals, err := models.GetChallengesWithSolves(db)
	if err != nil {
		return nil, err
	}
	ret := make([]ctf.Challenge, len(chals))
	for i := range chals {
		ret[i] = &chals[i]
	}
	return ret, nil
}

/*
func (j Jeopardy) GetSubmissions(db *gorm.DB) ([]ctf.Submission, error) {
	sub, err := models.GetSubmissions(db)
	if err != nil {
		return nil, err
	}
	ret := make([]ctf.Submission, len(sub))
	for i := range sub {
		sub[i].Challenge = nil
		sub[i].User = nil
		ret[i] = &sub[i]
	}
	return ret, nil
}
*/

func (j Jeopardy) GetTeam(db *gorm.DB, id int) (ctf.Team, error) {
	return models.GetUser(db, id)
}

func (j Jeopardy) GetScoring() ctf.Scoring {
	return nil
}
