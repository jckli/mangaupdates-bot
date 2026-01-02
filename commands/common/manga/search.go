package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunSearchEntry(
	r common.Responder,
	b *mubot.Bot,
	query string,
) error {
	if mangaID, err := strconv.ParseInt(query, 10, 64); err == nil {
		return sendMangaDetails(r, b, mangaID)
	}

	return RunSearchMenu(r, b, query)
}

func RunSearchMenu(
	r common.Responder,
	b *mubot.Bot,
	query string,
) error {
	results, err := b.ApiClient.SearchManga(query)
	if err != nil {
		return r.Error("Failed to search manga.")
	}
	if len(results) == 0 {
		return r.Error("No results found for `" + query + "`")
	}

	max := 25
	if len(results) < max {
		max = len(results)
	}

	description := fmt.Sprintf("Found %d results for `%s`.\nPlease select one from the dropdown below:\n\n", len(results), query)
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

	components := common.GenerateSearchDropdown("manga_search_select", "Select a manga...", results)

	embed := common.StandardEmbed("Search Results", description)
	return r.Respond(embed, components)
}

func HandleSearchSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}

	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)
	responder := &common.ComponentResponder{Event: e}

	return sendMangaDetails(responder, b, mangaID)
}

func sendMangaDetails(r common.Responder, b *mubot.Bot, mangaID int64) error {
	details, err := b.ApiClient.GetMangaDetails(mangaID)
	if err != nil {
		return r.Error("Failed to fetch details.")
	}

	botIcon := ""
	if self, ok := b.Client.Caches().SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	embed := common.GenerateDetailEmbed(*details, botIcon)

	return r.Respond(embed, nil)
}
