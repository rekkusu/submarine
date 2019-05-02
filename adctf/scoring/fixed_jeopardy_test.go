package scoring

import (
	"sort"
	"testing"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
)

func genJeopardy(chals []ctf.Challenge, teams []ctf.Team, submissions []ctf.Submission) rules.JeopardyRule {
	rule := jeopardy{
		Submission: make([]ctf.Submission, 0),
		Team:       make([]ctf.Team, 0),
	}
	for _, item := range submissions {
		rule.Submission = append(rule.Submission, item)
	}

	for _, item := range teams {
		rule.Team = append(rule.Team, item)
	}

	return rule
}

func Test_FixedJeopardy_GetRanking(t *testing.T) {
	chals := []ctf.Challenge{
		&jeopardyChallenge{1, "chal1", 100, "desc", "flag1"},
		&jeopardyChallenge{2, "chal2", 200, "desc", "flag2"},
		&jeopardyChallenge{3, "chal3", 400, "desc", "flag3"},
	}

	teams := []ctf.Team{
		&team{1, "team1"},
		&team{2, "team2"},
		&team{3, "team3"},
		&team{4, "team4"},
		&team{5, "team5"},
		&team{6, "team6"},
		&team{7, "team7"},
	}

	tests := []struct {
		submissions []ctf.Submission
		expect      Scores
	}{
		{
			[]ctf.Submission{
				&jeopardySubmission{1, teams[0], chals[0], "flag1", true, time.Date(2017, 8, 1, 12, 0, 0, 0, time.UTC)},
				&jeopardySubmission{2, teams[1], chals[0], "flag1", true, time.Date(2017, 8, 1, 12, 30, 0, 0, time.UTC)},
				&jeopardySubmission{3, teams[2], chals[0], "flag1", true, time.Date(2017, 8, 1, 13, 30, 0, 0, time.UTC)},
				&jeopardySubmission{4, teams[2], chals[2], "flagX", false, time.Date(2017, 8, 1, 14, 0, 0, 0, time.UTC)},
				&jeopardySubmission{5, teams[2], chals[2], "flag3", true, time.Date(2017, 8, 1, 15, 0, 0, 0, time.UTC)},
				&jeopardySubmission{6, teams[4], chals[0], "flag1", true, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
				&jeopardySubmission{7, teams[5], chals[0], "flag1", true, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
				&jeopardySubmission{8, teams[6], chals[0], "flag1", true, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
			},
			Scores{
				score{teams[2], 500, time.Date(2017, 8, 1, 15, 0, 0, 0, time.UTC)},
				score{teams[0], 100, time.Date(2017, 8, 1, 12, 0, 0, 0, time.UTC)},
				score{teams[1], 100, time.Date(2017, 8, 1, 15, 30, 0, 0, time.UTC)},
				score{teams[4], 100, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
				score{teams[5], 100, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
				score{teams[6], 100, time.Date(2017, 8, 1, 18, 0, 0, 0, time.UTC)},
				score{teams[3], 0, time.Time{}},
			},
		},
		{
			[]ctf.Submission{
				&jeopardySubmission{1, teams[0], chals[0], "flag1", true, time.Time{}},
				&jeopardySubmission{2, teams[1], chals[0], "flag1", true, time.Time{}},
				&jeopardySubmission{3, teams[2], chals[0], "flag1", true, time.Time{}},
			},
			Scores{
				score{teams[0], 100, time.Time{}},
				score{teams[1], 100, time.Time{}},
				score{teams[2], 100, time.Time{}},
				score{teams[3], 0, time.Time{}},
				score{teams[4], 0, time.Time{}},
				score{teams[5], 0, time.Time{}},
				score{teams[6], 0, time.Time{}},
			},
		},
	}

	for _, test := range tests {
		rule := genJeopardy(chals, teams, test.submissions)
		scoring := FixedJeopardy{rule}
		ranks := scoring.GetScores()
		sort.Stable(ranks)

		if len(test.expect) != len(ranks) {
			t.Errorf("len(GetRanking()) = %d, want %d\n", len(ranks), len(test.expect))
		}
		hasError := false
		for i := 0; i < len(test.expect); i++ {
			if test.expect[i].GetScore() != ranks[i].GetScore() || test.expect[i].GetTeam().GetName() != ranks[i].GetTeam().GetName() {
				hasError = true
			}
		}
		if hasError {
			for i := 0; i < len(test.expect); i++ {
				if test.expect[i].GetScore() != ranks[i].GetScore() || test.expect[i].GetTeam().GetName() != ranks[i].GetTeam().GetName() {
					t.Errorf("%s/%d != %s/%d\n",
						test.expect[i].GetTeam().GetName(), test.expect[i].GetScore(),
						ranks[i].GetTeam().GetName(), ranks[i].GetScore())
				} else {
					t.Errorf("%s/%d == %s/%d\n",
						test.expect[i].GetTeam().GetName(), test.expect[i].GetScore(),
						ranks[i].GetTeam().GetName(), ranks[i].GetScore())
				}
			}
		}
	}
}
