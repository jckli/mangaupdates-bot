package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/utils"
)

const (
	ColorPrimary = 0x3083e3
	ColorError   = 0xff4f4f
)

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

func metadataEmbed(title string, metadata utils.MangaMetadata) discord.Embed {
	embed := discord.NewEmbed().WithTitle(title).WithColor(ColorPrimary)
	if metadata.Description != "" {
		embed = embed.WithDescription(metadata.Description)
	}
	if metadata.Year != "" {
		embed = embed.AddField("Year", metadata.Year, true)
	}
	if metadata.Type != "" {
		embed = embed.AddField("Type", displayName(metadata.Type), true)
	}
	if metadata.LatestChapter != "" {
		chapterLabel := "Latest Chapter"
		if strings.EqualFold(metadata.Status, "completed") {
			chapterLabel = "Total Chapters"
		}
		embed = embed.AddField(chapterLabel, metadata.LatestChapter, true)
	}
	if len(metadata.Authors) > 0 {
		embed = embed.AddField("Authors", strings.Join(metadata.Authors, ", "), true)
	}
	if len(metadata.Artists) > 0 {
		embed = embed.AddField("Artists", strings.Join(metadata.Artists, ", "), true)
	}
	if metadata.Rating != nil {
		embed = embed.AddField("Rating", fmt.Sprintf("%.2f", *metadata.Rating), true)
	}
	if metadata.CoverURL != "" {
		embed = embed.WithImage(metadata.CoverURL)
	}
	return embed
}

func GenerateConfirmationEmbed(metadata utils.MangaMetadata) discord.Embed {
	return metadataEmbed(fmt.Sprintf("Is `%s` correct?", mangaTitle(metadata)), metadata)
}

func ErrorEmbed(content string) discord.Embed {
	return discord.NewEmbed().
		WithTitle("Error").
		WithDescription(content).
		WithColor(ColorError)
}

func GenerateDetailEmbed(metadata utils.MangaMetadata, botIconURL string) discord.Embed {
	return metadataEmbed(mangaTitle(metadata), metadata).
		WithURL("https://www.mangaupdates.com/series/"+strconv.FormatInt(metadata.ID, 36)).
		WithAuthor("MangaUpdates", "", botIconURL)
}

func mangaTitle(metadata utils.MangaMetadata) string {
	if metadata.Status == "" {
		return metadata.Title
	}
	return fmt.Sprintf("%s (%s)", metadata.Title, displayName(metadata.Status))
}

func displayName(value string) string {
	words := strings.Fields(strings.ReplaceAll(value, "_", " "))
	for index, word := range words {
		words[index] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return strings.Join(words, " ")
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
