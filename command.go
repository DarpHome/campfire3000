package main

import "github.com/bwmarrin/discordgo"

type CommandHandler func(*discordgo.Session, *discordgo.InteractionCreate)
type Command struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}
