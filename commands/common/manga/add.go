package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunAddEntry(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	query string,
) error {
	if mangaID, err := strconv.ParseInt(query, 10, 64); err == nil {
		return sendAddConfirmation(r, b, endpoint, mangaID)
	}

	return RunAddSearch(r, b, endpoint, query)
}

func RunAddSearch(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
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

	customID := fmt.Sprintf("manga_add_select/%s", endpoint)
	components := common.GenerateSearchDropdown(customID, "Select a manga to add...", results)

	embed := common.StandardEmbed("Search Results", description)
	return r.Respond(embed, components)
}

func HandleAddSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	mode := e.Vars["mode"]
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	responder := &common.ComponentResponder{Event: e}
	return sendAddConfirmation(responder, b, mode, mangaID)
}

func HandleAddConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	mode := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	action := e.Vars["action"]

	responder := &common.ComponentResponder{Event: e}

	if action == "no" {
		embed := common.StandardEmbed("Cancelled", "Manga was not added.")
		return responder.Respond(embed, nil)
	}

	var targetID string
	if mode == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	} else {
		targetID = e.User().ID.String()
	}

	err := b.ApiClient.AddMangaToWatchlist(mode, targetID, mangaID)
	if err != nil {
		return responder.Error(err.Error())
	}

	embed := common.StandardEmbed("Success", "Manga successfully added to the list.")
	return responder.Respond(embed, nil)
}

func sendAddConfirmation(r common.Responder, b *mubot.Bot, endpoint string, mangaID int64) error {
	details, err := b.ApiClient.GetMangaDetails(mangaID)
	if err != nil {
		return r.Error("Failed to fetch manga details.")
	}

	embed := common.GenerateConfirmationEmbed(*details)

	prefix := fmt.Sprintf("/manga_add_confirm/%s/%d", endpoint, mangaID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	return r.Respond(embed, buttons)
}
