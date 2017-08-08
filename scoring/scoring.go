package scoring

import "github.com/activedefense/submarine/ctf"

type Scoring interface {
	GetRanking() []Rank
}

type Rank struct {
	Rank  int      `json:"rank"`
	Team  ctf.Team `json:"team"`
	Score int      `json:"score"`
}
