package commands

import (
	"github.com/disgoorg/disgo/discord"
)

func errorTechnicalErrorEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("Something went wrong. Please ask for assistance in the support server.").
		SetColor(0xff4f4f).
		Build()
	return embed
}
