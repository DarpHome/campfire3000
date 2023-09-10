package main

import "github.com/bwmarrin/discordgo"

func Truncate[T any](arr []T, maxLen int) []T {
	if len(arr) < maxLen {
		return arr
	}
	return arr[:maxLen]
}

func RespondWithContent(session *discordgo.Session, interaction *discordgo.Interaction, flags discordgo.MessageFlags, content string) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flags,
		},
	})
}

func RespondWithError(session *discordgo.Session, interaction *discordgo.Interaction, flags discordgo.MessageFlags, guild *Guild, description string) error {
	if guild == nil {
		guild = &Guild{
			Locale: "en-US",
			Color:  GuildDefaultColor,
		}
	}
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Color:       int(guild.Color),
					Title:       GetStringI18nValue(guild.Locale, "cv3000.Error"),
					Description: description,
				},
			},
			Flags: flags,
		},
	})
}
