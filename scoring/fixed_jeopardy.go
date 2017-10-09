package scoring

import (
	"time"

	"github.com/activedefense/submarine/rules"
)

type FixedJeopardy struct {
	Jeopardy rules.JeopardyRule
}

func (scoring FixedJeopardy) GetScores() Scores {
	submissions, _ := scoring.Jeopardy.GetSubmissions()
	teams, _ := scoring.Jeopardy.GetTeams()
	teams_index := make(map[int]int)
	var scores Scores

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
