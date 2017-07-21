package models

import "github.com/jinzhu/gorm"

type Team interface {
	GetID() int
	GetName() string
}

type team struct {
	ID       int    `json:"id"`
	TeamName string `json:"team_name"`
}

func NewTeam(id int, teamName string) Team {
	return &team{id, teamName}
}

func (t team) GetID() int {
	return t.ID
}

func (t team) GetName() string {
	return t.TeamName
}

type User interface {
	GetID() int
	GetTeam() Team
	GetUsername() string
	GetPassword() string
}

type user struct {
	ID       int
	Team     *team
	Username string
	Password string
}

func (u *user) GetID() int {
	return u.ID
}

func (u *user) GetTeam() Team {
	return u.Team
}

func (u *user) GetUsername() string {
	return u.Username
}

func (u *user) GetPassword() string {
	return u.Password
}

type TeamRepository interface {
	AllTeams() ([]Team, error)
	GetTeam(id int) (Team, error)
	SaveTeam(t Team) error
}

type DefaultTeamRepository struct {
	DB *gorm.DB
}

func (repo *DefaultTeamRepository) AllTeams() ([]Team, error) {
	var teams []team
	if err := repo.DB.Find(&teams).Error; err != nil {
		return nil, err
	}

	result := make([]Team, len(teams))
	for i, _ := range teams {
		result[i] = &teams[i]
	}

	return result, nil
}

func (repo *DefaultTeamRepository) GetTeam(id int) (Team, error) {
	var t team
	if err := repo.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *DefaultTeamRepository) SaveTeam(t Team) error {
	team, ok := t.(*team)
	if !ok {
		return ErrModelMismatched
	}

	return repo.DB.Create(team).Error
}
