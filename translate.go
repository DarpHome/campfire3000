package main

import (
	"fmt"
	"strings"

	"github.com/EYERCORD/deepl-sdk-go/types"
	"github.com/MCausc78/cgorithm"
	"github.com/bwmarrin/discordgo"
)

func DeepLTranslate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.Interaction.ApplicationCommandData()
	switch interaction.Interaction.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		if data.Options[0].Name == "text" {
			data.Options = data.Options[1:]
		}
		input := ""
		switch {
		case data.Options[0].Focused:
			input = data.Options[0].StringValue()
		case data.Options[1].Focused:
			input = data.Options[1].StringValue()
		}
		input = strings.ToUpper(input)
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: Truncate(cgorithm.Transform(cgorithm.Filter(TDeepLLanguages, func(_ int, choice *discordgo.ApplicationCommandOptionChoice) bool {
					return input == "" || strings.Contains(strings.ToUpper(choice.Name), input) || strings.Contains(strings.ToUpper(choice.Value.(string)), input)
				}), func(_ int, choice *discordgo.ApplicationCommandOptionChoice) *discordgo.ApplicationCommandOptionChoice {
					return &discordgo.ApplicationCommandOptionChoice{
						Name:  fmt.Sprintf("%s: %s", choice.Value, choice.Name),
						Value: choice.Value,
					}
				}), 25),
			},
		})
	case discordgo.InteractionApplicationCommand:
		guild, err := GetScopeBasedInfo(interaction.Interaction)
		if err != nil {
			RespondWithContent(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, "Database error: "+err.Error())
			return
		}
		if TDeepLTranslator == nil {
			RespondWithError(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, guild, GetStringI18nValue(guild.Locale, "cv3000.UnavailableTranslator"))
			return
		}
		sourceLanguage := ""
		if len(data.Options) > 2 {
			sourceLanguage = data.Options[2].StringValue()
		}
		tr, err := TDeepLTranslator.Translate(
			data.Options[0].StringValue(),
			types.TargetLangCode(data.Options[1].StringValue()),
			types.SourceLangCode(sourceLanguage),
		)
		if err != nil {
			RespondWithError(session, interaction.Interaction, discordgo.MessageFlagsEphemeral, guild, err.Error())
			return
		}

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:  discordgo.EmbedTypeRich,
						Title: GetStringI18nValue(guild.Locale, "cv3000.TranslationResult"),
						Color: int(guild.Color),
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   GetStringI18nValue(guild.Locale, "cv3000.DetectedLanguage"),
								Value:  tr.DetectedLanguage,
								Inline: false,
							},
							{
								Name:   GetStringI18nValue(guild.Locale, "cv3000.TranslatedText"),
								Value:  tr.Result,
								Inline: false,
							},
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

var (
	TDeepLLanguages []*discordgo.ApplicationCommandOptionChoice = []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "Bulgarian",
			Value: "BG",
		},
		{
			Name:  "Czech",
			Value: "CS",
		},
		{
			Name:  "Danish",
			Value: "DA",
		},
		{
			Name:  "German",
			Value: "DE",
		},
		{
			Name:  "Greek",
			Value: "EL",
		},
		{
			Name:  "English (British)",
			Value: "EN-GB",
		},
		{
			Name:  "English (American)",
			Value: "EN-US",
		},
		{
			Name:  "Spanish",
			Value: "ES",
		},
		{
			Name:  "Estonian",
			Value: "ET",
		},
		{
			Name:  "Finnish",
			Value: "FI",
		},
		{
			Name:  "French",
			Value: "FR",
		},
		{
			Name:  "Hungarian",
			Value: "HU",
		},
		{
			Name:  "Indonesian",
			Value: "ID",
		},
		{
			Name:  "Italian",
			Value: "IT",
		},
		{
			Name:  "Japanese",
			Value: "JA",
		},
		{
			Name:  "Korean",
			Value: "KO",
		},
		{
			Name:  "Lithuanian",
			Value: "LT",
		},
		{
			Name:  "Latvian",
			Value: "LV",
		},
		{
			Name:  "Norwegian (Bokmål)",
			Value: "NB",
		},
		{
			Name:  "Dutch",
			Value: "NL",
		},
		{
			Name:  "Polish",
			Value: "PL",
		},
		{
			Name:  "Portuguese (Brazilian)",
			Value: "PT-BR",
		},
		{
			Name:  "Portuguese (all Portuguese varieties excluding Brazilian Portuguese)",
			Value: "PT-PT",
		},
		{
			Name:  "Romanian",
			Value: "RO",
		},
		{
			Name:  "Russian",
			Value: "RU",
		},
		{
			Name:  "Slovak",
			Value: "SK",
		},
		{
			Name:  "Slovenian",
			Value: "SL",
		},
		{
			Name:  "Swedish",
			Value: "SV",
		},
		{
			Name:  "Turkish",
			Value: "TR",
		},
		{
			Name:  "Ukrainian",
			Value: "UK",
		},
		{
			Name:  "Chinese (simplified)",
			Value: "ZH",
		},
	}
	TMinimalTextLength  int        = 1
	TLanguageCodeLength int        = 2
	TranslateCommands   []*Command = []*Command{
		{
			Command: &discordgo.ApplicationCommand{
				Type:         discordgo.ChatApplicationCommand,
				Name:         "deepl-translate",
				DMPermission: &BTrue,
				Description:  "Translate text using DeepL",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "text",
						Description:  "Text to be translated",
						Required:     true,
						Autocomplete: false,
						MinLength:    &TMinimalTextLength,
					},
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "target-language",
						Description:  "Language translate to",
						Required:     true,
						Autocomplete: true,
						MinLength:    &TLanguageCodeLength,
					},
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "source-language",
						Description:  "Language translate from",
						Required:     false,
						Autocomplete: true,
						MinLength:    &TLanguageCodeLength,
					},
				},
			},
			Handler: DeepLTranslate,
		},
	}
)
