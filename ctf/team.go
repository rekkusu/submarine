package ctf

type Team interface {
	GetID() int
	GetName() string
}

type User interface {
	GetID() int
	GetTeam() Team
	GetUsername() string
	GetPassword() string
	GetRole() string
}

type TeamStore interface {
	AllTeams() ([]Team, error)
	GetTeam(id int) (Team, error)
	SaveTeam(t Team) error
}
