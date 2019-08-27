package scoring

import (
	"github.com/activedefense/submarine/adctf/models"
	"time"

	"github.com/activedefense/submarine/ctf"
)

type ScoringRule interface {
	CalcScore(chal models.ChallengeWithSolves) int
	GetScores(chals []models.ChallengeWithSolves, submissions []models.Submission, teams []ctf.Team) []*ScoreRecord
	Update()
}

type ScoreRecord struct {
	Team           ctf.Team  `json:"team"`
	Score          int       `json:"score"`
	LastSubmission time.Time `json:"last_submission"`
}

func (s ScoreRecord) GetTeam() ctf.Team {
	return s.Team
}

func (s ScoreRecord) GetScore() int {
	return s.Score
}

func (s ScoreRecord) GetLastSubmission() time.Time {
	return s.LastSubmission
}

func (s ScoreRecord) LessThan(s2 ScoreRecord) bool {
	if s.GetScore() == s2.GetScore() {
		return s.GetLastSubmission().Before(s2.GetLastSubmission())
	}
	return s.GetScore() > s2.GetScore()
}
