package scoring

import (
	"reflect"
	"testing"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/activedefense/submarine/rules"
)

func genJeopardy(chals []ctf.Challenge, teams []ctf.Team, submissions []ctf.Submission) rules.JeopardyRule {
	rule := jeopardy{
		Challenge:  &challengeStore{make(map[int]ctf.Challenge)},
		Submission: &submissionStore{make(map[int]ctf.Submission)},
		Team:       &teamStore{make(map[int]ctf.Team)},
	}
	for _, item := range chals {
		rule.Challenge.Save(item)
	}

	for _, item := range teams {
		rule.Team.SaveTeam(item)
	}

	for _, item := range submissions {
		rule.Submission.Save(item)
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
		expect      []Rank
	}{
		{
			[]ctf.Submission{
				&jeopardySubmission{1, teams[0], chals[0], "flag1", true, time.Unix(1500883944, 0)},
				&jeopardySubmission{2, teams[1], chals[0], "flag1", true, time.Unix(1500883950, 0)},
				&jeopardySubmission{3, teams[2], chals[0], "flag1", true, time.Unix(1500884003, 0)},
				&jeopardySubmission{4, teams[2], chals[2], "flagX", false, time.Unix(1500884024, 0)},
				&jeopardySubmission{5, teams[2], chals[2], "flag3", true, time.Unix(1500884028, 0)},
				&jeopardySubmission{6, teams[4], chals[0], "flag1", true, time.Unix(1500884071, 0)},
				&jeopardySubmission{7, teams[5], chals[0], "flag1", true, time.Unix(1500884072, 0)},
				&jeopardySubmission{8, teams[6], chals[0], "flag1", true, time.Unix(1500884092, 0)},
			},
			[]Rank{
				Rank{1, teams[2], 500},
				Rank{2, teams[0], 100},
				Rank{3, teams[1], 100},
				Rank{4, teams[4], 100},
				Rank{5, teams[5], 100},
				Rank{6, teams[6], 100},
				Rank{7, teams[3], 0},
			},
		},
		{
			[]ctf.Submission{
				&jeopardySubmission{1, teams[0], chals[0], "flag1", true, time.Unix(1500883944, 0)},
				&jeopardySubmission{2, teams[1], chals[0], "flag1", true, time.Unix(1500883944, 0)},
				&jeopardySubmission{3, teams[2], chals[0], "flag1", true, time.Unix(1500883944, 0)},
			},
			[]Rank{
				Rank{1, teams[0], 100},
				Rank{2, teams[1], 100},
				Rank{3, teams[2], 100},
				Rank{4, teams[3], 0},
				Rank{5, teams[4], 0},
				Rank{6, teams[5], 0},
				Rank{7, teams[6], 0},
			},
		},
	}

	for _, test := range tests {
		rule := genJeopardy(chals, teams, test.submissions)
		scoring := FixedJeopardy{rule}
		ranks := scoring.GetRanking()

		if len(test.expect) != len(ranks) {
			t.Errorf("len(GetRanking()) = %d, want %d\n", len(ranks), len(test.expect))
		}
		if !reflect.DeepEqual(test.expect, ranks) {
			for i := 0; i < len(test.expect); i++ {
				if !reflect.DeepEqual(test.expect[i], ranks[i]) {
					t.Errorf("(%d)%s/%d != (%d)%s/%d\n",
						test.expect[i].Rank, test.expect[i].Team.GetName(), test.expect[i].Score,
						ranks[i].Rank, ranks[i].Team.GetName(), ranks[i].Score)
				} else {
					t.Logf("(%d)%s/%d == (%d)%s/%d\n",
						test.expect[i].Rank, test.expect[i].Team.GetName(), test.expect[i].Score,
						ranks[i].Rank, ranks[i].Team.GetName(), ranks[i].Score)
				}
			}
		}
	}
}
