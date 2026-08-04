package common

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaSearchSummary(result utils.MangaSearchResult) string {
	parts := make([]string, 0, 3)
	if result.Type != "" {
		parts = append(parts, strings.ToUpper(result.Type[:1])+result.Type[1:])
	}
	if result.Year != "" {
		parts = append(parts, result.Year)
	}
	if result.Rating > 0 {
		parts = append(parts, fmt.Sprintf("Rating: %.2f", result.Rating))
	}
	return strings.Join(parts, " · ")
}

func MangaAutocompleteSummary(result utils.MangaSearchResult) string {
	parts := make([]string, 0, 2)
	if result.Type != "" {
		parts = append(parts, strings.ToUpper(result.Type[:1])+result.Type[1:])
	}
	if result.Year != "" {
		parts = append(parts, result.Year)
	}
	return strings.Join(parts, ", ")
}

func TruncateText(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max-3]) + "..."
}

func GenerateSearchDropdown(
	customID string,
	placeholder string,
	results []utils.MangaSearchResult,
) []discord.LayoutComponent {

	var options []discord.StringSelectMenuOption

	max := 25
	if len(results) < max {
		max = len(results)
	}

	for _, res := range results[0:max] {
		option := discord.StringSelectMenuOption{
			Label: TruncateText(res.Title, 100),
			Value: fmt.Sprintf("%d", res.ID),
		}
		if summary := MangaSearchSummary(res); summary != "" {
			option.Description = TruncateText(summary, 100)
		}
		options = append(options, option)
	}

	return []discord.LayoutComponent{
		discord.NewActionRow(
			discord.NewStringSelectMenu(customID, placeholder, options...),
		),
	}
}
