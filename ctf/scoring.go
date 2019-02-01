package ctf

import "time"

type Scoring interface {
	CalcScore(chal Challenge) int
	GetScores(chals []Challenge, submissions []Submission, teams []Team) Scores
	Recalculate()
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
	GetTeam() Team
	GetScore() int
	GetLastSubmission() time.Time
	LessThan(Score) bool
}
