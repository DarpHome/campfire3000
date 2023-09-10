package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	DiscordSupportLink string = "https://discord.gg/EVtSwKttEH"
)

func Ping(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := GetScopeBasedInfo(interaction.Interaction)
	if err != nil {
		Logger.Error(err)
		RespondWithContent(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, "Database error: "+err.Error())
		return
	}
	sent := time.Now().UnixMilli()
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🏓 " + GetStringI18nValue(guild.Locale, "cv3000.Pong"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	spent := time.Now().UnixMilli() - sent
	content := "🏓 " + fmt.Sprintf(GetStringI18nValue(guild.Locale, "cv3000.PongLatency"), float64(spent)/1000)
	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})

}

func Dice(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := GetScopeBasedInfo(interaction.Interaction)
	if err != nil {
		Logger.Error(err)
		RespondWithContent(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, "Database error: "+err.Error())
		return
	}
	var sides int64 = 6
	options := interaction.Interaction.ApplicationCommandData().Options
	if len(options) != 0 {
		sides = options[0].IntValue()
	}
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf(GetStringI18nValue(guild.Locale, "cv3000.DiceResult"), 1+RNG.Int63n(sides-1)),
					Color: guild.Color,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func Bot(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := GetScopeBasedInfo(interaction.Interaction)
	if err != nil {
		Logger.Error(err)
		RespondWithContent(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, "Database error: "+err.Error())
		return
	}
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "GitHub",
							Style: discordgo.LinkButton,
							Emoji: discordgo.ComponentEmoji{
								Name:     "dh_github_octocat",
								ID:       "1147138555202773062",
								Animated: false,
							},
							URL: "https://github.com/DarpHome/campfire3000",
						},
						discordgo.Button{
							Style: discordgo.LinkButton,
							Emoji: discordgo.ComponentEmoji{
								Name:     "dh_discord_clyde_blurple",
								ID:       "1145500476150927420",
								Animated: false,
							},
							URL: DiscordSupportLink,
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "<:dh_campfire3000:1150029889395765268> " + GetStringI18nValue(guild.Locale, "cv3000.Bot"),
					Color: int(guild.Color),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   GetStringI18nValue(guild.Locale, "cv3000.Uptime"),
							Value:  fmt.Sprintf("<t:%d:R>", StartedAt),
							Inline: false,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func Config(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild, err := GetScopeBasedInfo(interaction.Interaction)
	if err != nil {
		Logger.Error(err)
		RespondWithContent(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, "Database error: "+err.Error())
		return
	}
	data := interaction.ApplicationCommandData()
	subOptions := data.Options[0].Options
	switch data.Options[0].Name {
	case "color":
		if len(subOptions) != 1 {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Color:       guild.Color,
							Description: fmt.Sprintf(GetStringI18nValue(guild.Locale, "cv3000.CurrentColor"), guild.Color),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		newColor := subOptions[0].IntValue()
		if newColor != int64(guild.Color) {
			_, err = Database.Exec("UPDATE guilds SET color = $1 WHERE id = $2", newColor, guild.ID)
			if err != nil {
				Logger.Error(err)
				RespondWithError(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, guild, err.Error())
				return
			}
		}
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       int(newColor),
						Description: fmt.Sprintf(GetStringI18nValue(guild.Locale, "cv3000.EmbedNewColorDescription"), guild.Color, newColor),
					},
				},
			},
		})

	case "locale":
		if len(subOptions) != 1 {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(GetStringI18nValue(guild.Locale, "cv3000.CurrentLocale"), GetStringI18nValue(guild.Locale, "cv3000.Language")),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		newLocale := data.Options[0].Options[0].StringValue()
		if newLocale != guild.Locale {
			_, err = Database.Exec("UPDATE guilds SET locale = $1 WHERE id = $2", newLocale, guild.ID)
			if err != nil {
				Logger.Error(err)
				RespondWithError(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, guild, err.Error())
				return
			}
		}
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "✅ " + fmt.Sprintf(GetStringI18nValue(newLocale, "cv3000.SwitchLocale"), GetStringI18nValue(newLocale, "cv3000.Language")),
			},
		})
	}
}

var (
	BFalse                    bool       = false
	BTrue                     bool       = true
	DMinDices                 float64    = 1.0
	CMinColor                 float64    = 0.0
	CDefaultMemberPermissions int64      = discordgo.PermissionManageServer
	BaseCommands              []*Command = []*Command{
		{
			Command: &discordgo.ApplicationCommand{
				Type:         discordgo.ChatApplicationCommand,
				Name:         "bot",
				DMPermission: &BTrue,
				Description:  "Get bot information.",
			},
			Handler: Bot,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Type:                     discordgo.ChatApplicationCommand,
				Name:                     "config",
				DefaultMemberPermissions: &CDefaultMemberPermissions,
				DMPermission:             &BFalse,
				Description:              "Configure",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "color",
						Description: "Change the embed color",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "new-color",
								Description: "New embed color",
								Required:    false,
								MinValue:    &CMinColor,
								MaxValue:    0xFFFFFF,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "locale",
						Description: "Change the locale",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "new-locale",
								Description: "New locale",
								Required:    false,
								Choices: []*discordgo.ApplicationCommandOptionChoice{
									{
										Name:  "English (American)",
										Value: "en-US",
									},
									{
										Name:  "Russian",
										Value: "ru",
									},
									{
										Name:  "Ukrainian",
										Value: "ua",
									},
								},
							},
						},
					},
				},
			},
			Handler: Config,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Type:         discordgo.ChatApplicationCommand,
				Name:         "dice",
				DMPermission: &BTrue,
				Description:  "Roll a dice.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "sides",
						Description: "Sides of dice",
						MinValue:    &DMinDices,
						Required:    false,
					},
				},
			},
			Handler: Dice,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Type:         discordgo.ChatApplicationCommand,
				Name:         "ping",
				DMPermission: &BTrue,
				Description:  "Pong!",
			},
			Handler: Ping,
		},
	}
)
