package scoring

import "github.com/activedefense/submarine/ctf"

type Scoring interface {
	GetRanking() []Rank
}

type Rank struct {
	Rank  int
	Team  ctf.Team
	Score int
}
