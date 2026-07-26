package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
)

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Pong!",
}

func PingHandler(e *handler.CommandEvent) error {
	var ping string
	if e.Client().HasGateway() {
		ping = e.Client().Gateway.Latency().String()
	}

	embed := discord.NewEmbed().
		WithTitle("Pong! 🏓").
		WithDescription("My ping is " + ping).
		WithColor(common.ColorPrimary).
		WithTimestamp(e.CreatedAt())

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreate().WithEmbeds(embed),
	)
}
