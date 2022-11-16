package entities

type Team struct {
	Entity `mapstructure:",squash"`
}

type AddTeam struct {
	Entity `mapstructure:",squash"`
}
