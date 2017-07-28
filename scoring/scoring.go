package scoring

import "github.com/activedefense/submarine/models"

type Scoring interface {
	GetRanking() []Rank
}

type Rank struct {
	Rank  int
	Team  models.Team
	Score int
}
