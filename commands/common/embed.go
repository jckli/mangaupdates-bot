package common

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/utils"
	"strings"
)

const (
	ColorPrimary = 0x3083e3
	ColorError   = 0xff4f4f
)

// helper functions
func formatAuthorsAndArtists(list []utils.MangaAuthor) (string, string) {
	var authors []string
	var artists []string

	for _, person := range list {
		switch person.Type {
		case "Author":
			authors = append(authors, person.Name)
		case "Artist":
			artists = append(artists, person.Name)
		default:
			authors = append(authors, person.Name)
		}
	}

	aStr := "N/A"
	if len(authors) > 0 {
		aStr = strings.Join(authors, ", ")
	}

	artStr := "N/A"
	if len(artists) > 0 {
		artStr = strings.Join(artists, ", ")
	}

	return aStr, artStr
}

func formatStatus(completed bool) string {
	if completed {
		return "Completed"
	}
	return "Ongoing"
}

// buttons
func CreateConfirmButtons(confirmID, cancelID string) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton("Cancel", cancelID),
			discord.NewSuccessButton("Confirm", confirmID),
		},
	}
}

// actual embeds
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
	authorStr, artistStr := formatAuthorsAndArtists(details.Authors)

	embed := discord.NewEmbedBuilder().
		SetTitle(fmt.Sprintf("Is `%s` correct?", details.Title)).
		SetDescription(details.Description).
		SetColor(ColorPrimary).
		AddField("Year", details.Year, true).
		AddField("Type", details.Type, true).
		AddField("Latest Chapter", fmt.Sprintf("%d", details.LatestChapter), true).
		AddField("Authors", authorStr, true).
		AddField("Artists", artistStr, true).
		AddField("Rating", fmt.Sprintf("%.2f", details.BayesianRating), true)

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

func GenerateDetailEmbed(details utils.MangaDetails, botIconURL string) discord.Embed {
	authorStr, artistStr := formatAuthorsAndArtists(details.Authors)

	embed := discord.NewEmbedBuilder().
		SetAuthor("MangaUpdates", "", botIconURL).
		SetTitle(fmt.Sprintf("%s (%s)", details.Title, formatStatus(details.Completed))).
		SetURL(details.URL).
		SetDescription(details.Description).
		SetColor(ColorPrimary).
		AddField("Year", details.Year, true).
		AddField("Type", details.Type, true).
		AddField("Latest Chapter", fmt.Sprintf("%d", details.LatestChapter), true).
		AddField("Authors", authorStr, true).
		AddField("Artists", artistStr, true).
		AddField("Rating", fmt.Sprintf("%.2f", details.BayesianRating), true)

	if details.Image != nil {
		embed.SetImage(details.Image.URL.Original)
	}

	return embed.Build()
}
