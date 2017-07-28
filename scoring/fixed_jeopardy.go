package scoring

import (
	"sort"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
)

type FixedJeopardy struct {
	Jeopardy rules.JeopardyRule
}

type status struct {
	team  ctf.Team
	score int
	last  time.Time
}
type ranking []*status

func (r ranking) Len() int {
	return len(r)
}

func (r ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ranking) Less(i, j int) bool {
	a := r[i]
	b := r[j]
	if a.score == b.score {
		if a.last == b.last {
			return false
		} else {
			return a.last.Before(b.last)
		}
	} else {
		return a.score > b.score
	}
}

func (score FixedJeopardy) GetRanking() []Rank {
	submissions, _ := score.Jeopardy.GetSubmissionStore().All()
	teams, _ := score.Jeopardy.GetTeamStore().AllTeams()
	teams_index := make(map[int]int)
	var ranking ranking = make([]*status, len(teams))

	for i, team := range teams {
		teams_index[team.GetID()] = i
		ranking[i] = &status{team, 0, time.Time{}}
	}

	for _, item := range submissions {
		if !item.IsCorrect() {
			continue
		}

		team := ranking[teams_index[item.GetTeam().GetID()]]
		team.score += item.GetScore()
		if team.last.Before(item.GetSubmittedAt()) {
			team.last = item.GetSubmittedAt()
		}
	}

	sort.Stable(ranking)

	result := make([]Rank, len(ranking))
	for i, item := range ranking {
		result[i] = Rank{i + 1, item.team, item.score}
	}

	return result
}
