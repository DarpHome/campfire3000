package main

import (
	"strconv"

	"github.com/MCausc78/cgorithm"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

func OnMemberJoined(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	rows, err := Database.Query("SELECT role_id FROM autoroles WHERE guild_id = $1", AlwaysValidU64(event.GuildID))
	if err != nil {
		Logger.Error(err)
		return
	}
	roles := []string{}
	for rows.Next() {
		var roleId uint64
		rows.Scan(&roleId)
		roles = append(roles, strconv.FormatUint(roleId, 10))
	}
	if len(roles) == 0 {
		return
	}
	member, err := session.GuildMemberEdit(event.GuildID, event.User.ID, &discordgo.GuildMemberParams{
		Roles: &roles,
	})
	if err != nil {
		Logger.Error(err)
		return
	}
	nonExistingRoles := cgorithm.Filter(roles, func(_ int, roleId string) bool {
		return slices.Contains(member.Roles, roleId)
	})
	if len(nonExistingRoles) == 0 {
		return
	}
	Database.Exec("DELETE FROM autoroles WHERE guild_id = $1 AND role_id IN $2", event.GuildID, nonExistingRoles)
}
