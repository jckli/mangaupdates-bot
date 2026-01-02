package manga

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
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

	return RunRemoveMenu(r, b, endpoint, targetID, query, 1)
}

func RunRemoveMenu(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	targetID string,
	query string,
	page int,
) error {
	list, err := b.ApiClient.GetWatchlist(endpoint, targetID)
	if err != nil {
		return r.Error("Failed to fetch watchlist.")
	}
	if list == nil || len(*list) == 0 {
		return r.Error("Your watchlist is empty.")
	}

	var matches []utils.MangaSearchResult
	var displayLines []string

	queryLower := strings.ToLower(query)

	for _, item := range *list {
		if query == "" || strings.Contains(strings.ToLower(item.Title), queryLower) {

			matches = append(matches, utils.MangaSearchResult{ID: item.ID, Title: item.Title})

			line := fmt.Sprintf("• %s", item.Title)
			if item.GroupName != "" {
				line += fmt.Sprintf(" (*%s*)", item.GroupName)
			}
			displayLines = append(displayLines, line)
		}
	}

	if len(matches) == 0 {
		return r.Error(fmt.Sprintf("No manga found matching `%s`.", query))
	}

	slicedMatches, totalPages := common.GetPageSlice(matches, page)
	slicedLines, _ := common.GetPageSlice(displayLines, page)

	customSelectID := fmt.Sprintf("manga_remove_select/%s", endpoint)
	dropdown := common.GenerateSearchDropdown(customSelectID, "Select manga to remove...", slicedMatches)

	safeQuery := query
	if safeQuery == "" {
		safeQuery = "-"
	}
	btnPrefix := fmt.Sprintf("/remove_nav/%s/%s", endpoint, safeQuery)
	buttons := common.GeneratePaginationButtons(btnPrefix, page, totalPages)

	allComponents := append(dropdown, buttons...)

	var description string
	if len(slicedLines) == 0 {
		description = "No items found."
	} else {
		description = ""
		for _, line := range slicedLines {
			description += line + "\n"
		}
	}

	header := fmt.Sprintf("Found %d results.", len(matches))
	if query != "" {
		header = fmt.Sprintf("Found %d results for `%s`.", len(matches), query)
	}
	fullDescription := header + "\nPlease select one from the dropdown below:\n\n" + description

	embed := common.StandardEmbed("Select Manga to Remove", fullDescription)
	embed.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Page %d/%d", page, totalPages)}

	return r.Respond(embed, allComponents)
}

func HandleRemovePagination(e *handler.ComponentEvent, b *mubot.Bot) error {
	endpoint := e.Vars["mode"]
	query := e.Vars["query"]
	page, _ := strconv.Atoi(e.Vars["page"])

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

	responder := &common.ComponentResponder{Event: e}
	return RunRemoveMenu(responder, b, endpoint, targetID, query, page)
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
