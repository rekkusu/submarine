package scoring

import (
	"errors"
	"sort"
	"time"

	"github.com/activedefense/submarine/models"
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
	Team        models.Team
	Challenge   models.Challenge
	Answer      string
	Correct     bool
	SubmittedAt time.Time
}

func (s jeopardySubmission) GetID() int {
	return s.ID
}

func (s jeopardySubmission) GetTeam() models.Team {
	return s.Team
}

func (s jeopardySubmission) GetUser() models.User {
	return nil
}

func (s jeopardySubmission) GetChallenge() models.Challenge {
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
	records map[int]models.Challenge
}

func (store *challengeStore) All() ([]models.Challenge, error) {
	array := make([]models.Challenge, len(store.records))
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

func (store *challengeStore) Get(id int) (models.Challenge, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *challengeStore) Save(c models.Challenge) error {
	store.records[c.GetID()] = c
	return nil
}

type submissionStore struct {
	records map[int]models.Submission
}

func (store *submissionStore) All() ([]models.Submission, error) {
	array := make([]models.Submission, len(store.records))
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

func (store *submissionStore) Get(id int) (models.Submission, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *submissionStore) Save(s models.Submission) error {
	store.records[s.GetID()] = s
	return nil
}

type teamStore struct {
	records map[int]models.Team
}

func (store *teamStore) AllTeams() ([]models.Team, error) {
	array := make([]models.Team, len(store.records))
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

func (store *teamStore) GetTeam(id int) (models.Team, error) {
	val, ok := store.records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return val, nil
}

func (store *teamStore) SaveTeam(t models.Team) error {
	store.records[t.GetID()] = t
	return nil
}

type jeopardy struct {
	Challenge  models.ChallengeStore
	Team       models.TeamStore
	Submission models.SubmissionStore
}

func (j jeopardy) GetChallengeStore() models.ChallengeStore {
	return j.Challenge
}

func (j jeopardy) GetTeamStore() models.TeamStore {
	return j.Team
}

func (j jeopardy) GetSubmissionStore() models.SubmissionStore {
	return j.Submission
}
