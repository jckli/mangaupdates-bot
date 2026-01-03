package manga

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

type GlobalSearchConfig struct {
	Query          string
	SelectIDPrefix string
	EndpointSuffix string
	Title          string
	Placeholder    string
}

func GenerateGlobalSearchMenu(
	b *mubot.Bot,
	cfg GlobalSearchConfig,
) (discord.Embed, []discord.ContainerComponent, error) {
	results, err := b.ApiClient.SearchManga(cfg.Query)
	if err != nil {
		return discord.Embed{}, nil, fmt.Errorf("failed to search manga")
	}
	if len(results) == 0 {
		return discord.Embed{}, nil, fmt.Errorf("no results found for `%s`", cfg.Query)
	}

	max := 25
	if len(results) < max {
		max = len(results)
	}

	description := fmt.Sprintf("Found %d results for `%s`.\nPlease select one from the dropdown below:\n\n", len(results), cfg.Query)
	for i, res := range results[0:max] {
		line := fmt.Sprintf("`%d.` %s", i+1, res.Title)
		if res.Year != "" {
			line += fmt.Sprintf(" (%s)", res.Year)
		}
		if res.Rating > 0 {
			line += fmt.Sprintf(" • Rating: %.2f", res.Rating)
		}
		description += line + "\n"
	}

	customID := cfg.SelectIDPrefix
	if cfg.EndpointSuffix != "" {
		customID += "/" + cfg.EndpointSuffix
	}

	components := common.GenerateSearchDropdown(customID, cfg.Placeholder, results[0:max])

	embed := common.StandardEmbed(cfg.Title, description)
	return embed, components, nil
}
