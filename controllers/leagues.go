package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewLeague(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.AddLeague
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	insertLeagueDetails := entities.InsertLeague{
		Members: []entities.MemberInfo{
			{
				ID:   UID,
				Role: entities.Admin,
			},
		},
		Name: input.Name,
	}
	leagueId, appError := lib.InsertItem(ctx, "leagues", insertLeagueDetails)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	if appError = signUserToLeague(ctx, UID, leagueId); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leagueId": leagueId})
}

func GetLeagueInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	leagueId, appError := lib.GetRequiredParam(ctx, "id")
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	league, appError := lib.GetSingle[entities.League](ctx, "leagues", leagueId)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	if !entities.LeagueContainsMember(league, UID) {
		appError := lib.Error(http.StatusUnauthorized, "not-league-member")
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	members := make([]entities.Member, len(league.Members))
	for index, member := range league.Members {
		user, team, appError := getUserTeam(ctx, member.ID)
		if appError.HasError() {
			ctx.JSON(appError.Code, appError.Json)
			return
		}

		members[index] = entities.Member{
			Entity:   user.Entity,
			Nickname: user.Nickname,
			Team:     team,
			Role:     member.Role,
		}
	}

	detailedLeague := entities.DetailedLeague{
		Entity:  league.Entity,
		Members: members,
		Name:    league.Name,
	}
	ctx.JSON(http.StatusOK, gin.H{"league": detailedLeague})
}
