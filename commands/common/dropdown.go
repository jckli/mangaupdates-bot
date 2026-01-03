package common

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/utils"
)

func GenerateSearchDropdown(
	customID string,
	placeholder string,
	results []utils.MangaSearchResult,
) []discord.ContainerComponent {

	var options []discord.StringSelectMenuOption

	max := 25
	if len(results) < max {
		max = len(results)
	}

	for _, res := range results[0:max] {
		label := res.Title
		if len(label) > 95 {
			label = label[:95] + "..."
		}

		var desc string
		var parts []string

		if res.Year != "" {
			parts = append(parts, fmt.Sprintf("Year: %s", res.Year))
		}
		if res.Rating > 0 {
			parts = append(parts, fmt.Sprintf("Rating: %.2f", res.Rating))
		}

		if len(parts) > 0 {
			desc = strings.Join(parts, " | ")
		}

		option := discord.StringSelectMenuOption{
			Label: label,
			Value: fmt.Sprintf("%d", res.ID),
		}
		if desc != "" {
			option.Description = desc
		}
		options = append(options, option)
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    customID,
				Placeholder: placeholder,
				Options:     options,
			},
		},
	}
}
