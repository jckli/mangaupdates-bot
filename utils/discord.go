package utils

import (
	"github.com/disgoorg/disgo/discord"
)

func DcErrorTechnicalErrorEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("Something went wrong. Please ask for assistance in the support server.").
		SetColor(0xff4f4f).
		Build()
	return embed
}
