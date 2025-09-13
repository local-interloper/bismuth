package bismuth

import "github.com/bwmarrin/discordgo"

type CommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}
