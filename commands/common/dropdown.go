package common

import (
	"fmt"

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

		desc := fmt.Sprintf("Year: %s | Rating: %.2f", res.Year, res.Rating)

		options = append(options, discord.StringSelectMenuOption{
			Label:       label,
			Value:       fmt.Sprintf("%d", res.ID),
			Description: desc,
		})
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
