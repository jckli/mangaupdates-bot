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
func CreateConfirmButtons(confirmID, cancelID string) []discord.LayoutComponent {
	return []discord.LayoutComponent{
		discord.NewActionRow(
			discord.NewDangerButton("Cancel", cancelID),
			discord.NewSuccessButton("Confirm", confirmID),
		),
	}
}

// actual embeds
func StandardEmbed(title, description string) discord.Embed {
	return discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithColor(ColorPrimary)
}

func GenerateListEmbed(
	title string,
	iconURL string,
	description string,
	totalItems int,
	botIconURL string,
) discord.Embed {
	return discord.NewEmbed().
		WithAuthor("MangaUpdates", "", botIconURL).
		WithTitle(title).
		WithThumbnail(iconURL).
		WithDescription(description).
		WithColor(ColorPrimary).
		WithFooterText(fmt.Sprintf("Total: %d", totalItems))
}

func GenerateConfirmationEmbed(details utils.MangaDetails) discord.Embed {
	authorStr, artistStr := formatAuthorsAndArtists(details.Authors)

	embed := discord.NewEmbed().
		WithTitlef("Is `%s` correct?", details.Title).
		WithDescription(details.Description).
		WithColor(ColorPrimary).
		AddField("Year", details.Year, true).
		AddField("Type", details.Type, true).
		AddField("Latest Chapter", fmt.Sprintf("%d", details.LatestChapter), true).
		AddField("Authors", authorStr, true).
		AddField("Artists", artistStr, true).
		AddField("Rating", fmt.Sprintf("%.2f", details.BayesianRating), true)

	if details.Image != nil {
		embed = embed.WithImage(details.Image.URL.Original)
	}

	return embed
}

func ErrorEmbed(content string) discord.Embed {
	return discord.NewEmbed().
		WithTitle("Error").
		WithDescription(content).
		WithColor(ColorError)
}

func GenerateDetailEmbed(details utils.MangaDetails, botIconURL string) discord.Embed {
	authorStr, artistStr := formatAuthorsAndArtists(details.Authors)

	embed := discord.NewEmbed().
		WithAuthor("MangaUpdates", "", botIconURL).
		WithTitlef("%s (%s)", details.Title, formatStatus(details.Completed)).
		WithURL(details.URL).
		WithDescription(details.Description).
		WithColor(ColorPrimary).
		AddField("Year", details.Year, true).
		AddField("Type", details.Type, true).
		AddField("Latest Chapter", fmt.Sprintf("%d", details.LatestChapter), true).
		AddField("Authors", authorStr, true).
		AddField("Artists", artistStr, true).
		AddField("Rating", fmt.Sprintf("%.2f", details.BayesianRating), true)

	if details.Image != nil {
		embed = embed.WithImage(details.Image.URL.Original)
	}

	return embed
}

func GenerateGroupConfirmationEmbed(group *utils.GroupDetails, mangaTitle string) discord.Embed {
	if group == nil || group.GroupID == 0 {
		return discord.NewEmbed().
			WithTitlef("Clear filter for `%s`?", mangaTitle).
			WithDescription("You are **clearing** the scanlation group filter.\n\nYou will receive notifications for **all** releases.").
			WithColor(ColorPrimary)
	}

	embed := discord.NewEmbed().
		WithTitlef("Set filter for `%s`?", mangaTitle).
		WithDescriptionf("You are limiting notifications to **%s**.", group.Name).
		WithColor(ColorPrimary).
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
		embed = embed.AddField("Links", strings.Join(links, " | "), false)
	} else {
		embed = embed.AddField("Links", "N/A", false)
	}

	return embed
}
