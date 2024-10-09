package bot

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	instance *bot
)

// create and start the session, wait for an interrupt signal, and shut down
func Run(token string) int {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("failed to create discordgo session", slog.Any("error", err))
		return 1
	}

	session.AddHandler(newMessage)
	session.AddHandler(interactionCreate)

	err = session.Open()
	if err != nil {
		slog.Error("failed to open session", slog.Any("error", err))
		return 1
	}
	defer session.Close()

	if session.State == nil || session.State.User == nil {
		slog.Error("no valid user in session")
		return 1
	}
	instance = New(session.State.User.ID)
	instance.registerCommands(session)

	slog.Info("bot is running")
	waitForInterrupt()
	slog.Info("interrupt receieved, bot shutting down")
	return 0
}

func waitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	instance.newMessage(session, message)
}

func interactionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	instance.interactionCreate(session, interaction)
}
