package bismuth

import "github.com/bwmarrin/discordgo"

type MessageProcessor = func(s *discordgo.Session, m *discordgo.MessageCreate) bool
