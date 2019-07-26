package scoring

import (
	"github.com/activedefense/submarine/adctf/models"
	"time"
)

type FixedJeopardy struct {
}

func (scoring FixedJeopardy) CalcScore(chal models.ChallengeWithSolves) int {
	return chal.Point
}

func (scoring FixedJeopardy) Update() {
}

func (scoring FixedJeopardy) GetScores(chals []models.ChallengeWithSolves, submissions []models.Submission, teams []models.User) []*ScoreRecord {
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

		var chal models.ChallengeWithSolves
		for _, c := range chals {
			if c.GetID() == sub.GetChallengeID() {
				chal = c
				break
			}
		}

		score := scores[teams_index[sub.TeamID]]
		score.Score += scoring.CalcScore(chal)

		if score.LastSubmission.Before(sub.GetSubmittedAt()) {
			score.LastSubmission = sub.GetSubmittedAt()
		}
	}

	return scores
}
