package manga

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func HandleAutocomplete(e *handler.AutocompleteEvent, b *mubot.Bot, queryName string) error {
	query := e.Data.String(queryName)

	if len(query) < 3 {
		return e.AutocompleteResult(nil)
	}

	results, err := b.ApiClient.SearchManga(query)
	if err != nil {
		return e.AutocompleteResult(nil)
	}

	var choices []discord.AutocompleteChoice

	max := 25
	if len(results) < max {
		max = len(results)
	}

	for _, res := range results[0:max] {
		label := fmt.Sprintf("%s (%s)", res.Title, res.Year)

		if len(label) > 100 {
			label = label[:97] + "..."
		}

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  label,
			Value: fmt.Sprintf("%d", res.ID),
		})
	}

	return e.AutocompleteResult(choices)
}
