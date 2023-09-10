package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MCausc78/cgorithm"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	Logger           *logrus.Logger            = nil
	Handlers         map[string]CommandHandler = map[string]CommandHandler{}
	TDeepLTranslator *DeepLTranslator          = nil
	Database         *sql.DB                   = nil
	StartedAt        int64                     = 0
)

func main() {
	logger := new(logrus.Logger)
	logger.Out = os.Stderr
	logger.Formatter = new(logrus.TextFormatter)
	logger.Hooks = make(logrus.LevelHooks)
	logger.Level = logrus.InfoLevel
	logger.SetReportCaller(true)
	Logger = logger
	if err := godotenv.Load(); err != nil {
		Logger.Fatal(err)
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"),
	))
	if err != nil {
		Logger.Fatal(err)
	}
	Database = db
	if err != nil {
		Logger.Fatal(err)
	}
	InitializeRandom()
	translator, err := NewDeepLTranslator(os.Getenv("DEEPL_API_KEY"))
	if err != nil {
		Logger.Error(err)
	} else {
		TDeepLTranslator = translator
	}
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		Logger.Fatal(err)
	}
	session.Identify.Intents = discordgo.IntentGuilds |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildWebhooks
	session.Identify.Presence.Status = string(discordgo.StatusIdle)
	session.Identify.Presence.Game = discordgo.Activity{
		Type: discordgo.ActivityTypeCompeting,
		Name: "/help | Try /imagine, /gpt & /coinprice",
	}
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if h, ok := Handlers[interaction.ApplicationCommandData().Name]; ok {
			h(session, interaction)
		}
	})
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		StartedAt = time.Now().Unix()
		Logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	session.AddHandler(OnMemberJoined)
	InitializeI18n()
	if err = session.Open(); err != nil {
		Logger.Fatal(err)
	}
	Handlers = map[string]CommandHandler{}
	commands := []*Command{}
	groups := [][]*Command{
		BaseCommands,
		TranslateCommands,
	}
	for _, group := range groups {
		commands = append(commands, group...)
	}
	if _, err := session.ApplicationCommandBulkOverwrite(
		session.State.Application.ID,
		"",
		cgorithm.Transform(commands, func(_ int, command *Command) *discordgo.ApplicationCommand {
			Handlers[command.Command.Name] = command.Handler
			return command.Command
		}),
	); err != nil {
		Logger.Fatal(err)
	}
	fmt.Println("Bot ran. Press CTRL-C to stop.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	if err = session.Close(); err != nil {
		Logger.Error(err)
	}
	if err = db.Close(); err != nil {
		Logger.Error(err)
	}
}
