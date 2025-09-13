package bismuth

import (
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestCreateBot(t *testing.T) {
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		t.Error("BOT_TOKEN not provided")
		return
	}

	bismuth, err := NewBismuth(botToken)

	if err != nil {
		t.Error("Failed to create Bismuth instance")
		return
	}

	bismuth.RegisterCommand(Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "foo",
			Description: "Replies with 'bar'",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "LOL",
				},
			})
		},
	})

	bismuth.RegisterMessageProcessor(func(s *discordgo.Session, m *discordgo.MessageCreate) bool {
		msg := strings.ToLower(m.Message.Content)

		if !strings.Contains(msg, "fuck") {
			return true
		}

		s.ChannelMessageSend(m.ChannelID, "Bad language is not allowed")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return false
	})

	bismuth.RegisterMessageProcessor(func(s *discordgo.Session, m *discordgo.MessageCreate) bool {
		msg := strings.ToLower(m.Message.Content)

		if !strings.Contains(msg, "fizz") {
			return true
		}

		s.ChannelMessageSendReply(m.Message.ChannelID, "buzz", m.Message.Reference())

		return true
	})

	if err := bismuth.Start(); err != nil {
		t.Error("Failed to init Bismuth")
		return
	}
}
