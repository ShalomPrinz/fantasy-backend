package entities

type AddLeagueInvitation struct {
	To       string `json:"to" binding:"required"`
	LeagueId string `json:"leagueId" binding:"required"`
}

type InsertLeagueInvitation struct {
	From     string `json:"from"`
	LeagueId string `json:"leagueId"`
}

type Message struct {
	Entity   `mapstructure:",squash"`
	From     string `json:"from"`
	LeagueId string `json:"leagueId"`
}
