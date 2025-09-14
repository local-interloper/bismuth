package bismuth

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session           *discordgo.Session
	commands          map[string]Command
	messageProcessors []MessageProcessor
}

func NewBot(token string) (*Bot, error) {
	app := &Bot{}

	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	app.session = session

	app.commands = map[string]Command{}

	return app, nil
}

func (b *Bot) RegisterCommand(command Command) {
	b.commands[command.Command.Name] = command
}

func (b *Bot) RegisterCommands(commands []Command) {
	for _, command := range commands {
		b.RegisterCommand(command)
	}
}

func (b *Bot) RegisterMessageProcessor(processor MessageProcessor) {
	b.messageProcessors = append(b.messageProcessors, processor)
}

func (b *Bot) RegisterMessageProcessors(processors []MessageProcessor) {
	for _, processor := range processors {
		b.RegisterMessageProcessor(processor)
	}
}

func (b *Bot) initCommands() {
	commands, err := b.session.ApplicationCommands(b.session.State.Application.ID, "")

	if err != nil {
		log.Panic("Failed to cleanup commands")
	}

	for _, command := range commands {
		b.session.ApplicationCommandDelete(b.session.State.Application.ID, "", command.ID)
	}

	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		command, ok := b.commands[i.ApplicationCommandData().Name]

		if !ok {
			return
		}

		command.Handler(s, i)
	})

	for _, command := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.Application.ID, "", command.Command)
		if err != nil {
			log.Println("Failed to register command", command.Command.Name, "\nError:", err.Error())
		}
	}
}

func (b *Bot) initMessageProcessing() {
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Message.Author.ID == s.State.User.ID {
			return
		}

		for _, processor := range b.messageProcessors {
			shouldContinue := processor(s, m)

			if !shouldContinue {
				break
			}
		}
	})
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}

	b.initCommands()
	b.initMessageProcessing()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	b.session.Close()

	return nil
}
