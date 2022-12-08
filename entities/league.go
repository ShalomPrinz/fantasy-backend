package entities

type League struct {
	Entity  `mapstructure:",squash"`
	Members []string `json:"members"`
	Name    string   `json:"name"`
}

type AddLeague struct {
	Members []string `json:"members"`
	Name    string   `json:"name"`
}
