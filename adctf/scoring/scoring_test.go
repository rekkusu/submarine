package scoring

import (
	"errors"
	"sort"
	"time"

	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type jeopardyChallenge struct {
	ID          int
	Title       string
	Point       int
	Description string
	Flag        string
}

func (c jeopardyChallenge) GetID() int {
	return c.ID
}

func (c jeopardyChallenge) GetTitle() string {
	return c.Title
}

func (c jeopardyChallenge) GetPoint() int {
	return c.Point
}

func (c jeopardyChallenge) GetDescription() string {
	return c.Description
}

func (c jeopardyChallenge) GetFlag() string {
	return c.Flag
}

type jeopardySubmission struct {
	ID          int
	Team        ctf.Team
	Challenge   ctf.Challenge
	Answer      string
	Correct     bool
	SubmittedAt time.Time
}

func (s jeopardySubmission) GetID() int {
	return s.ID
}

func (s jeopardySubmission) GetTeam() ctf.Team {
	return s.Team
}

func (s jeopardySubmission) GetUser() ctf.User {
	return nil
}

func (s jeopardySubmission) GetChallenge() ctf.Challenge {
	return s.Challenge
}

func (s jeopardySubmission) GetAnswer() string {
	return s.Answer
}

func (s jeopardySubmission) GetScore() int {
	return s.Challenge.GetPoint()
}

func (s jeopardySubmission) IsCorrect() bool {
	return s.Correct
}

func (s jeopardySubmission) GetSubmittedAt() time.Time {
	return s.SubmittedAt
}

type team struct {
	ID   int
	Name string
}

func (t team) GetID() int {
	return t.ID
}

func (t team) GetName() string {
	return t.Name
}

type challengeStore struct {
	records map[int]ctf.Challenge
}

func (store *challengeStore) All() ([]ctf.Challenge, error) {
	array := make([]ctf.Challenge, len(store.records))
	i := 0
	for _, item := range store.records {
		array[i] = item
		i += 1
	}
	sort.Slice(array, func(i, j int) bool {
		return array[i].GetID() < array[j].GetID()
	})
	return array, nil
}

func (store *challengeStore) Get(id int) (ctf.Challenge, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *challengeStore) Save(c ctf.Challenge) error {
	store.records[c.GetID()] = c
	return nil
}

type submissionStore struct {
	records map[int]ctf.Submission
}

func (store *submissionStore) All() ([]ctf.Submission, error) {
	array := make([]ctf.Submission, len(store.records))
	i := 0
	for _, item := range store.records {
		array[i] = item
		i += 1
	}
	sort.Slice(array, func(i, j int) bool {
		return array[i].GetID() < array[j].GetID()
	})
	return array, nil
}

func (store *submissionStore) Get(id int) (ctf.Submission, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *submissionStore) Save(s ctf.Submission) error {
	store.records[s.GetID()] = s
	return nil
}

type teamStore struct {
	records map[int]ctf.Team
}

func (store *teamStore) AllTeams() ([]ctf.Team, error) {
	array := make([]ctf.Team, len(store.records))
	i := 0
	for _, item := range store.records {
		array[i] = item
		i += 1
	}
	sort.Slice(array, func(i, j int) bool {
		return array[i].GetID() < array[j].GetID()
	})
	return array, nil
}

func (store *teamStore) GetTeam(id int) (ctf.Team, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *teamStore) SaveTeam(t ctf.Team) error {
	store.records[t.GetID()] = t
	return nil
}

type jeopardy struct {
	Team       []ctf.Team
	Submission []ctf.Submission
}

func (j jeopardy) GetDB() *gorm.DB {
	return nil
}

func (j jeopardy) GetSubmissions() ([]ctf.Submission, error) {
	return j.Submission, nil
}

func (j jeopardy) GetTeams() ([]ctf.Team, error) {
	return j.Team, nil
}

func (j jeopardy) GetTeam(id int) (ctf.Team, error) {
	for _, item := range j.Team {
		if item.GetID() == id {
			return item, nil
		}
	}
	return nil, errors.New("not found")
}
