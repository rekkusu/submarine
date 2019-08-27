package ctf

type Team interface {
	GetID() int
	GetName() string
}

type User interface {
	GetID() int
	GetTeamID() int
	GetUsername() string
	GetPassword() string
	GetRole() string
}
