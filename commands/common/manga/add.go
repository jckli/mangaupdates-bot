package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
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

	embed, components, err := GenerateGlobalSearchMenu(b, GlobalSearchConfig{
		Query:          query,
		SelectIDPrefix: "manga_add_select",
		EndpointSuffix: endpoint,
		Title:          "Search Results",
		Placeholder:    "Select a manga to add...",
	})
	if err != nil {
		return r.Error(err.Error())
	}
	return r.Respond(embed, components)
}

func HandleAddSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	_ = e.DeferUpdateMessage()

	go func() {
		if err := common.GuardWidget(e, b, true); err != nil {
			common.SendInteractionError(e, err.Error())
			return
		}

		endpoint := e.Vars["mode"]
		if len(e.StringSelectMenuInteractionData().Values) == 0 {
			return
		}
		mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

		metadata, err := b.ApiClient.GetMangaMetadata(mangaID)
		if err != nil {
			common.SendInteractionError(e, err.Error())
			return
		}

		embed := common.GenerateConfirmationEmbed(*metadata)
		prefix := fmt.Sprintf("/manga_add_confirm/%s/%d", endpoint, mangaID)
		buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{embed},
				Components: &buttons,
			})
	}()

	return nil
}

func HandleAddConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	action := e.Vars["action"]

	if action == "no" {
		_, err := e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Manga was not added."),
				},
				Components: &[]discord.LayoutComponent{},
			})
		return err
	}

	var targetID string
	if endpoint == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	} else {
		targetID = e.User().ID.String()
	}

	err := b.ApiClient.AddMangaToWatchlist(endpoint, targetID, mangaID)
	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = common.ColorError
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errEmbed},
				Components: &[]discord.LayoutComponent{},
			})
		return err
	}

	_, err = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("Success", "Manga successfully added to the list."),
			},
			Components: &[]discord.LayoutComponent{},
		})
	return err
}

func sendAddConfirmation(r common.Responder, b *mubot.Bot, endpoint string, mangaID int64) error {
	metadata, err := b.ApiClient.GetMangaMetadata(mangaID)
	if err != nil {
		return r.Error("Failed to fetch manga details.")
	}

	embed := common.GenerateConfirmationEmbed(*metadata)

	prefix := fmt.Sprintf("/manga_add_confirm/%s/%d", endpoint, mangaID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	return r.Respond(embed, buttons)
}
