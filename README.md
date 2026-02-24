# bismuth

[![Go](https://github.com/local-interloper/bismuth/actions/workflows/go.yml/badge.svg)](https://github.com/local-interloper/bismuth/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/badge/go-1.24-blue?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/local-interloper/bismuth.svg)](https://pkg.go.dev/github.com/local-interloper/bismuth)
[![Go Report Card](https://goreportcard.com/badge/github.com/local-interloper/bismuth)](https://goreportcard.com/report/github.com/local-interloper/bismuth)
[![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A lightweight Go library for building Discord bots, built on top of [discordgo](https://github.com/bwmarrin/discordgo).

## Installation

```bash
go get github.com/local-interloper/bismuth
```

## Usage

### Creating a bot

```go
bot, err := bismuth.NewBot("your-bot-token")
if err != nil {
    log.Fatal(err)
}
```

### Registering slash commands

```go
bot.RegisterCommand(bismuth.Command{
    Command: &discordgo.ApplicationCommand{
        Name:        "ping",
        Description: "Replies with Pong!",
    },
    Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Pong!",
            },
        })
    },
})
```

### Registering message processors

Message processors receive every message and return a `bool` indicating whether the next processor in the chain should run.

```go
bot.RegisterMessageProcessor(func(s *discordgo.Session, m *discordgo.MessageCreate) bool {
    if m.Content == "hello" {
        s.ChannelMessageSend(m.ChannelID, "Hello!")
        return false // stop processing
    }
    return true // continue to next processor
})
```

### Starting the bot

```go
if err := bot.Start(); err != nil {
    log.Fatal(err)
}
```

`Start` blocks until an interrupt signal is received, then gracefully closes the session.

## API

| Function | Description |
|---|---|
| `NewBot(token string) (*Bot, error)` | Creates a new bot with the given token |
| `(*Bot).RegisterCommand(Command)` | Registers a single slash command |
| `(*Bot).RegisterCommands([]Command)` | Registers multiple slash commands |
| `(*Bot).RegisterMessageProcessor(MessageProcessor)` | Registers a single message processor |
| `(*Bot).RegisterMessageProcessors([]MessageProcessor)` | Registers multiple message processors |
| `(*Bot).Start() error` | Opens the session and blocks until interrupted |

## License

[MIT](LICENSE)
