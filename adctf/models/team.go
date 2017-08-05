package models

import (
	"github.com/activedefense/submarine/ctf"
	"github.com/jinzhu/gorm"
)

type Team struct {
	ID       int    `json:"id"`
	Username string `json:"username" gorm:"unique_index"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

func (t Team) GetID() int {
	return t.ID
}

func (t Team) GetName() string {
	return t.Username
}

func (u Team) GetUsername() string {
	return u.Username
}

func (u Team) GetPassword() string {
	return u.Password
}

func (u Team) GetTeam() ctf.Team {
	return u
}

type TeamStore struct {
	DB *gorm.DB
}

func (store *TeamStore) AllTeams() ([]ctf.Team, error) {
	var teams []Team
	if err := store.DB.Find(&teams).Error; err != nil {
		return nil, err
	}

	result := make([]ctf.Team, len(teams))
	for i, _ := range teams {
		result[i] = &teams[i]
	}

	return result, nil
}

func (store *TeamStore) GetTeam(id int) (ctf.Team, error) {
	var t Team
	if err := store.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (store *TeamStore) GetTeamByName(name string) (ctf.Team, error) {
	var t Team
	if err := store.DB.First(&t, "username=?", name).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (store *TeamStore) SaveTeam(t ctf.Team) error {
	team, ok := t.(*Team)
	if !ok {
		return ctf.ErrModelMismatched
	}

	return store.DB.Create(team).Error
}
