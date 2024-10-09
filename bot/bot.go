package bot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	id string
}

func New(id string) *bot {
	return &bot{
		id: id,
	}
}

func (b *bot) registerCommands(session *discordgo.Session) {
	command := &discordgo.ApplicationCommand{
		Name:        "cuddle-bot",
		Description: "Information about cuddle-bot",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "Get information about cuddle-bot",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "help",
						Value: "help",
					},
				},
			},
		},
	}
	_, err := session.ApplicationCommandCreate(instance.id, "", command)
	if err != nil {
		slog.Error("failed to create slash command", slog.Any("error", err))
	}
}

func (b *bot) newMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore own messages
	if b.id == message.Author.ID {
		return
	}
	slog.Debug("newMessage", slog.Any("message", message.Content))

	// respond to user message if it starts with an @me
	parts := strings.Split(message.Content, " ")
	if parts[0] != fmt.Sprintf("<@%s>", b.id) {
		return
	}
	if len(parts) == 1 {
		session.ChannelMessageSend(message.ChannelID, "you have to tell me what you want :weary:")
		return
	}
	parts = parts[1:]
	switch {
	case parts[0] == "hi":
		session.ChannelMessageSend(message.ChannelID, "sup sup :sunglasses:")
	default:
		session.ChannelMessageSend(message.ChannelID, "sorry, i don't follow :sweat_smile:")
	}
}

func (b *bot) interactionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	slog.Debug("created", slog.Any("interaction", interaction))

	// ignore interactions that aren't application commands for now
	if interaction.Type != discordgo.InteractionApplicationCommand {
		slog.Info("ignoring non-application command interaction")
		return
	}

	// see if we can get a user from the interaction, and make sure it isn't us
	var userID string
	if interaction.User != nil {
		userID = interaction.User.ID
	} else if interaction.Member != nil {
		userID = interaction.Member.User.ID
	} else {
		slog.Error("no user available in interaction")
		return
	}
	slog.Debug("interaction", slog.Any("user", userID))
	if b.id == userID {
		slog.Error("somehow got command from self")
		return
	}

	// read the command data
	data := interaction.ApplicationCommandData()
	json, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		slog.Error("failed to marshal interaction application command data", slog.Any("error", err))
	}
	slog.Debug("interaction application command data=" + string(json))

	// respond to the command
	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":question: you know as much as I do, dawg...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("failed to respond to interaction", slog.Any("error", err))
	}
}
