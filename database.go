package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type GuildPunishment int

const (
	GuildPunishmentNone    GuildPunishment = 0
	GuildPunishmentKick    GuildPunishment = 1
	GuildPunishmentBan     GuildPunishment = 2
	GuildPunishmentTimeout GuildPunishment = 3
	GuildPunishmentWarn    GuildPunishment = 4
)

type Guild struct {
	ID                          uint64
	Locale                      string
	Color                       int
	Flags                       uint64
	MaxWarns                    int
	FinalWarnPunishment         GuildPunishment
	FinalWarnPunishmentDuration int64
}

const (
	GuildDefaultColor int = 0x39b440
)

func FindGuild(id uint64) (*Guild, error) {
	rows, err := Database.Query(
		"SELECT id, locale, color, flags, max_warns, final_warn_punishment, final_warn_punishment_duration FROM guilds WHERE id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		rows, err = Database.Query(
			"INSERT INTO guilds(id, locale, color, flags, max_warns, final_warn_punishment, final_warn_punishment_duration) "+
				"VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
			id,
			"en-US",
			GuildDefaultColor,
			0,
			3,
			GuildPunishmentNone,
			0,
		)
		if err != nil {
			return nil, err
		}
		rows.Next()
	}
	guild := &Guild{}
	rows.Scan(
		&guild.ID,
		&guild.Locale,
		&guild.Color,
		&guild.Flags,
		&guild.MaxWarns,
		&guild.FinalWarnPunishment,
		&guild.FinalWarnPunishmentDuration,
	)
	return guild, nil
}

func AlwaysValidU64(target string) uint64 {
	r, _ := strconv.ParseUint(target, 10, 64)
	return r
}

func GetScopeBasedInfo(interaction *discordgo.Interaction) (*Guild, error) {
	if interaction.GuildID == "" {
		return &Guild{
			ID:     0,
			Locale: "en-US",
			Color:  GuildDefaultColor,
			Flags:  0,
		}, nil
	}
	return FindGuild(AlwaysValidU64(interaction.GuildID))
}
