package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunRemoveEntry(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	targetID string,
	query string,
) error {
	if mangaID, err := strconv.ParseInt(query, 10, 64); err == nil {
		return sendRemoveConfirmation(r, b, endpoint, mangaID)
	}

	embed, components, err := GenerateWatchlistMenu(b, WatchlistMenuConfig{
		Endpoint:            endpoint,
		TargetID:            targetID,
		Query:               query,
		Page:                1,
		SelectIDPrefix:      "manga_remove_select",
		NavIDPrefix:         "/remove_nav",
		Title:               "Select Manga to Remove",
		DropdownPlaceholder: "Select manga to remove...",
	})
	if err != nil {
		return r.Error(err.Error())
	}
	return r.Respond(embed, components)
}

func HandleRemovePagination(e *handler.ComponentEvent, b *mubot.Bot) error {
	endpoint := e.Vars["mode"]
	query := e.Vars["query"]
	if query == "-" {
		query = ""
	}

	targetID := e.User().ID.String()
	if endpoint == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	}

	return HandleGenericWatchlistPagination(e, b, WatchlistMenuConfig{
		Endpoint:            endpoint,
		TargetID:            targetID,
		Query:               query,
		SelectIDPrefix:      "manga_remove_select",
		NavIDPrefix:         "/remove_nav",
		Title:               "Select Manga to Remove",
		DropdownPlaceholder: "Select manga to remove...",
	})
}

func HandleRemoveSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	mode := e.Vars["mode"]
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}

	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)
	responder := &common.ComponentResponder{Event: e}

	return sendRemoveConfirmation(responder, b, mode, mangaID)
}

func HandleRemoveConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	mode := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	action := e.Vars["action"]

	responder := &common.ComponentResponder{Event: e}

	if action == "no" {
		return responder.Respond(common.StandardEmbed("Cancelled", "No manga was removed."), nil)
	}

	targetID := e.User().ID.String()
	if mode == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	}

	err := b.ApiClient.RemoveMangaFromWatchlist(mode, targetID, mangaID)
	if err != nil {
		return responder.Error(err.Error())
	}

	return responder.Respond(common.StandardEmbed("Success", "Manga removed from watchlist."), nil)
}

func sendRemoveConfirmation(r common.Responder, b *mubot.Bot, endpoint string, mangaID int64) error {
	details, err := b.ApiClient.GetMangaDetails(mangaID)
	if err != nil {
		return r.Error("Failed to fetch details.")
	}

	embed := common.GenerateConfirmationEmbed(*details)

	prefix := fmt.Sprintf("/manga_remove_confirm/%s/%d", endpoint, mangaID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	return r.Respond(embed, buttons)
}
