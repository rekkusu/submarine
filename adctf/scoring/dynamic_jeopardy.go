package scoring

import (
	"github.com/activedefense/submarine/adctf/models"
	"time"
)

type DynamicJeopardy struct {
	Expression func(base int, weight int) int
}

func (scoring DynamicJeopardy) CalcScore(chal models.ChallengeWithSolves) int {
	return scoring.Expression(chal.Point, chal.Solves)
}

func (scoring DynamicJeopardy) Update() {
}

func (scoring DynamicJeopardy) GetScores(chals []models.ChallengeWithSolves, submissions []models.Submission, teams []models.Team) []*ScoreRecord {
	teams_index := make(map[int]int)
	var scores []*ScoreRecord

	for i, team := range teams {
		teams_index[team.GetID()] = i
		scores = append(scores, &ScoreRecord{team, 0, time.Time{}})
	}

	for _, sub := range submissions {
		if !sub.IsCorrect() {
			continue
		}

		score := scores[teams_index[sub.TeamID]]

		// Convert Challenge(WithoutSolves) -> ChallengeWithSolves
		var chal models.ChallengeWithSolves
		for _, c := range chals {
			if c.GetID() == sub.GetChallengeID() {
				chal = c
				break
			}
		}

		score.Score += scoring.CalcScore(chal)
		if score.LastSubmission.Before(sub.GetSubmittedAt()) {
			score.LastSubmission = sub.GetSubmittedAt()
		}
	}

	return scores
}
