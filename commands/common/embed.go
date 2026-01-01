package common

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/utils"
)

const (
	ColorPrimary = 0x3083e3
	ColorError   = 0xff4f4f
)

func StandardEmbed(title, description string) discord.Embed {
	return discord.NewEmbedBuilder().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorPrimary).
		Build()
}

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

func GenerateConfirmationEmbed(details utils.MangaDetails) discord.Embed {
	var authorStr string
	for _, a := range details.Authors {
		authorStr += fmt.Sprintf("%s, ", a.Name)
	}
	if len(authorStr) > 2 {
		authorStr = authorStr[:len(authorStr)-2]
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(fmt.Sprintf("Is `%s` correct?", details.Title)).
		SetDescription(details.Description).
		SetColor(ColorPrimary).
		AddField("Year", details.Year, true).
		AddField("Rating", fmt.Sprintf("%.2f", details.BayesianRating), true).
		AddField("Authors", authorStr, false)

	if details.Image != nil {
		embed.SetImage(details.Image.URL.Original)
	}

	return embed.Build()
}

func ErrorEmbed(content string) discord.Embed {
	return discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription(content).
		SetColor(ColorError).
		Build()
}

func CreateConfirmButtons(confirmID, cancelID string) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton("Cancel", cancelID),
			discord.NewSuccessButton("Confirm", confirmID),
		},
	}
}
