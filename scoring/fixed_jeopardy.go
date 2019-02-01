package scoring

import (
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
)

type FixedJeopardy struct {
	Jeopardy rules.JeopardyRule
}

func (scoring FixedJeopardy) CalcScore(chal ctf.Challenge) int {
	return chal.GetPoint()
}

func (scoring FixedJeopardy) Recalculate() {
}

func (scoring FixedJeopardy) GetScores(chals []ctf.Challenge, submissions []ctf.Submission, teams []ctf.Team) ctf.Scores {
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
		score.Score += item.GetChallenge().GetPoint()
		if score.LastSubmission.Before(item.GetSubmittedAt()) {
			score.LastSubmission = item.GetSubmittedAt()
		}
	}

	return scores
}
