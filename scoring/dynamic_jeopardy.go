package scoring

import (
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
)

type DynamicJeopardy struct {
	Jeopardy   rules.JeopardyRule
	Expression func(base int, weight int) int
}

func (scoring DynamicJeopardy) CalcScore(chal ctf.Challenge) int {
	return scoring.Expression(chal.GetPoint(), chal.GetSolves())
}

func (scoring DynamicJeopardy) Recalculate() {
}

func (scoring DynamicJeopardy) GetScores() ctf.Scores {
	chals, _ := scoring.Jeopardy.GetChallenges()
	submissions, _ := scoring.Jeopardy.GetSubmissions()
	teams, _ := scoring.Jeopardy.GetTeams()
	teams_index := make(map[int]int)
	var scores ctf.Scores

	for i, team := range teams {
		teams_index[team.GetID()] = i
		scores = append(scores, &score{team, 0, time.Time{}})
	}

	for _, item := range submissions {
		if !item.IsCorrect() {
			continue
		}

		score := scores[teams_index[item.GetTeam().GetID()]].(*score)
		var chal ctf.Challenge
		for _, c := range chals {
			if c.GetID() == item.GetChallenge().GetID() {
				chal = c
				break
			}
		}

		score.Score += scoring.CalcScore(chal)
		if score.LastSubmission.Before(item.GetSubmittedAt()) {
			score.LastSubmission = item.GetSubmittedAt()
		}
	}

	return scores
}
