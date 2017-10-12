package scoring

import (
	"time"

	"github.com/activedefense/submarine/ctf"
)

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

func (s score) LessThan(s2 ctf.Score) bool {
	if s.GetScore() == s2.GetScore() {
		return s.GetLastSubmission().Before(s2.GetLastSubmission())
	}
	return s.GetScore() > s2.GetScore()
}
