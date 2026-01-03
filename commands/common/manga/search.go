package manga

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
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

	embed, components, err := GenerateGlobalSearchMenu(b, GlobalSearchConfig{
		Query:          query,
		SelectIDPrefix: "manga_search_select",
		Title:          "Search Results",
		Placeholder:    "Select a manga...",
	})
	if err != nil {
		return r.Error(err.Error())
	}
	return r.Respond(embed, components)
}

func HandleSearchSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()

	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	details, err := b.ApiClient.GetMangaDetails(mangaID)
	if err != nil {
		return err
	}

	botIcon := ""
	if self, ok := b.Client.Caches().SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	embed := common.GenerateDetailEmbed(*details, botIcon)

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &[]discord.ContainerComponent{},
		})
	return err
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
