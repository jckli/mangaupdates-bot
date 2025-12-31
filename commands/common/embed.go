package common

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
)

const (
	ColorPrimary = 0x3083e3
	ColorError   = 0xff4f4f
)

func GenerateListEmbed(
	title string,
	iconURL string,
	description string,
	totalItems int,
	botIconURL string,
) discord.Embed {
	return discord.NewEmbedBuilder().
		SetAuthor("MangaUpdates", "", botIconURL).
		SetTitle(title).
		SetThumbnail(iconURL).
		SetDescription(description).
		SetColor(ColorPrimary).
		SetFooterText(fmt.Sprintf("Total: %d", totalItems)).
		Build()
}

func ErrorEmbed(content string) discord.Embed {
	return discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription(content).
		SetColor(ColorError).
		Build()
}
