package entities

type AddLeagueInvitation struct {
	To       string `json:"to" binding:"required"`
	LeagueId string `json:"leagueId" binding:"required"`
}

type Message struct {
	From     string `json:"from"`
	LeagueId string `json:"leagueId"`
}
