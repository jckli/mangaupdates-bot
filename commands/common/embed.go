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

func formatActive(active bool) string {
	if active {
		return "Yes"
	}
	return "No"
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

func GenerateGroupConfirmationEmbed(group *utils.GroupDetails, mangaTitle string) discord.Embed {
	if group == nil || group.GroupID == 0 {
		return discord.NewEmbedBuilder().
			SetTitle(fmt.Sprintf("Clear filter for `%s`?", mangaTitle)).
			SetDescription("You are **clearing** the scanlation group filter.\n\nYou will receive notifications for **all** releases.").
			SetColor(ColorPrimary).
			Build()
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(fmt.Sprintf("Set filter for `%s`?", mangaTitle)).
		SetDescription(fmt.Sprintf("You are limiting notifications to **%s**.", group.Name)).
		SetColor(ColorPrimary).
		AddField("Group Name", group.Name, true).
		AddField("Active", formatActive(group.Active), true)

	var links []string
	if group.Social.Site != "" {
		links = append(links, fmt.Sprintf("[Website](%s)", group.Social.Site))
	}
	if group.Social.Discord != "" {
		val := group.Social.Discord
		if strings.HasPrefix(val, "http") {
			links = append(links, fmt.Sprintf("[Discord](%s)", val))
		} else {
			links = append(links, fmt.Sprintf("Discord: `%s`", val))
		}
	}
	if group.Social.Twitter != "" {
		val := group.Social.Twitter
		if strings.HasPrefix(val, "http") {
			links = append(links, fmt.Sprintf("[Twitter](%s)", val))
		} else {
			links = append(links, fmt.Sprintf("Twitter: `%s`", val))
		}
	}

	if len(links) > 0 {
		embed.AddField("Links", strings.Join(links, " | "), false)
	} else {
		embed.AddField("Links", "N/A", false)
	}

	return embed.Build()
}
