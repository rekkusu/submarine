package scoring

import (
	"time"

	"github.com/activedefense/submarine/ctf"
)

type Scoring interface {
	GetScores() Scores
}

type Scores []Score

func (s Scores) Len() int {
	return len(s)
}

func (s Scores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Scores) Less(i, j int) bool {
	return s[i].LessThan(s[j])
}

type Score interface {
	GetTeam() ctf.Team
	GetScore() int
	GetLastSubmission() time.Time
	LessThan(Score) bool
}

type score struct {
	Team           ctf.Team  `json:"team"`
	Score          int       `json:"score"`
	LastSubmission time.Time `json:"last_submission"`
}

func (s score) GetTeam() ctf.Team {
	return s.Team
}

func (s score) GetScore() int {
	return s.Score
}

func (s score) GetLastSubmission() time.Time {
	return s.LastSubmission
}

func (s score) LessThan(s2 Score) bool {
	if s.GetScore() == s2.GetScore() {
		return s.GetLastSubmission().Before(s2.GetLastSubmission())
	}
	return s.GetScore() > s2.GetScore()
}
