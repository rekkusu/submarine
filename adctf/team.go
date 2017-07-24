package adctf

import (
	"github.com/activedefense/submarine/models"
	"github.com/jinzhu/gorm"
)

type Team struct {
	ID       int    `json:"id"`
	TeamName string `json:"team_name"`
}

func (t Team) GetID() int {
	return t.ID
}

func (t Team) GetName() string {
	return t.TeamName
}

type User struct {
	ID       int
	Team     *Team
	Username string
	Password string
}

func (u User) GetID() int {
	return u.ID
}

func (u User) GetTeam() models.Team {
	return u.Team
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetPassword() string {
	return u.Password
}

type TeamStore struct {
	DB *gorm.DB
}

func (store *TeamStore) AllTeams() ([]models.Team, error) {
	var teams []Team
	if err := store.DB.Find(&teams).Error; err != nil {
		return nil, err
	}

	result := make([]models.Team, len(teams))
	for i, _ := range teams {
		result[i] = &teams[i]
	}

	return result, nil
}

func (store *TeamStore) GetTeam(id int) (models.Team, error) {
	var t Team
	if err := store.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (store *TeamStore) SaveTeam(t models.Team) error {
	team, ok := t.(*Team)
	if !ok {
		return models.ErrModelMismatched
	}

	return store.DB.Create(team).Error
}
