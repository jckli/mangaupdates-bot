package manga

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"strconv"
	"strings"
)

func HandleAddAutocomplete(e *handler.AutocompleteEvent, b *mubot.Bot, queryName string) error {
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

func HandleWatchlistAutocomplete(
	e *handler.AutocompleteEvent,
	b *mubot.Bot,
	endpoint string,
	targetID string,
	queryName string,
) error {
	query := e.Data.String(queryName)

	list, err := b.ApiClient.GetWatchlist(endpoint, targetID)
	if err != nil || list == nil {
		return e.AutocompleteResult(nil)
	}

	var choices []discord.AutocompleteChoice
	queryLower := strings.ToLower(query)
	count := 0

	for _, item := range *list {
		if count >= 25 {
			break
		}

		if query == "" || strings.Contains(strings.ToLower(item.Title), queryLower) {

			label := item.Title
			if item.GroupName != "" {
				label += fmt.Sprintf(" [%s]", item.GroupName)
			}
			if len(label) > 100 {
				label = label[:97] + "..."
			}

			val := fmt.Sprintf("%d", item.ID)

			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  label,
				Value: val,
			})
			count++
		}
	}

	return e.AutocompleteResult(choices)
}

func HandleSetGroupAutocomplete(e *handler.AutocompleteEvent, b *mubot.Bot, endpoint, targetID string) error {
	focused := e.Data.Focused()

	switch focused.Name {
	case "title":
		return HandleWatchlistAutocomplete(e, b, endpoint, targetID, "title")

	case "group":
		mangaIDStr := e.Data.String("title")
		mangaID, err := strconv.ParseInt(mangaIDStr, 10, 64)
		if err != nil {
			return e.AutocompleteResult(nil)
		}
		groups, err := b.ApiClient.GetMangaGroups(mangaID)
		if err != nil {
			return e.AutocompleteResult(nil)
		}
		query := strings.ToLower(e.Data.String("group"))
		var choices []discord.AutocompleteChoice
		if strings.Contains(strings.ToLower("All Groups (Clear Filter)"), query) {
			choices = append(choices, discord.AutocompleteChoiceString{
				Name:  "All Groups (Clear Filter)",
				Value: "0",
			})
		}
		count := 0
		for _, g := range groups {
			if count >= 24 {
				break
			}

			if query == "" || strings.Contains(strings.ToLower(g.Name), query) {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  g.Name,
					Value: fmt.Sprintf("%d", g.ID),
				})
				count++
			}
		}
		return e.AutocompleteResult(choices)
	}

	return e.AutocompleteResult(nil)
}
