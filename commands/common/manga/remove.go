package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
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
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}
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
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	details, err := b.ApiClient.GetMangaDetails(mangaID)
	if err != nil {
		return err
	}

	embed := common.GenerateConfirmationEmbed(*details)
	prefix := fmt.Sprintf("/manga_remove_confirm/%s/%d", endpoint, mangaID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &buttons,
		})
	return err
}

func HandleRemoveConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	action := e.Vars["action"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "No manga was removed."),
				},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	targetID := e.User().ID.String()
	if endpoint == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	}

	err := b.ApiClient.RemoveMangaFromWatchlist(endpoint, targetID, mangaID)
	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = common.ColorError
		_, _ = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errEmbed},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("Success", "Manga removed from watchlist."),
			},
			Components: &[]discord.ContainerComponent{},
		})
	return err
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
