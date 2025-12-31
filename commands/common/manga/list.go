package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunMangaList(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	id string,
	displayName string,
	displayIcon string,
	page int,
) error {
	list, err := b.ApiClient.GetWatchlist(endpoint, id)
	if err != nil {
		b.Logger.Error("Failed to fetch manga list", "error", err)
		return r.Error("Technical error fetching list.")
	}
	if list == nil {
		return r.Error(fmt.Sprintf("The **%s** manga list is not set up yet.\nUse `/%s setup` to create it.", endpoint, endpoint))
	}

	formattedItems := make([]string, len(*list))
	for i, m := range *list {
		line := fmt.Sprintf("• %s", m.Title)
		if m.GroupName != "" {
			line += fmt.Sprintf(" (*%s*)", m.GroupName)
		}
		formattedItems[i] = line
	}

	slicedItems, totalPages := common.GetPageSlice(formattedItems, page)

	var description string
	if len(slicedItems) == 0 {
		description = "No items found."
	} else {
		for _, item := range slicedItems {
			description += item + "\n"
		}
	}

	botIcon := ""
	if self, ok := b.Client.Caches().SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}
	embed := common.GenerateListEmbed(
		fmt.Sprintf("%s's Manga List", displayName),
		displayIcon,
		description,
		len(formattedItems),
		botIcon,
	)

	components := common.GeneratePaginationButtons(
		fmt.Sprintf("/list_nav/%s", endpoint),
		page,
		totalPages,
	)

	return r.Respond(embed, components)
}

func HandleMangaListPagination(e *handler.ComponentEvent, b *mubot.Bot) error {
	mode := e.Vars["mode"]
	page, _ := strconv.Atoi(e.Vars["page"])

	var targetID, name, icon string

	if mode == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
		if g, ok := e.Guild(); ok {
			name = g.Name
			if i := g.IconURL(); i != nil {
				icon = *i
			}
		}
	} else {
		targetID = e.User().ID.String()
		name = e.User().EffectiveName()
		icon = e.User().EffectiveAvatarURL()
	}

	responder := &common.ComponentResponder{Event: e}

	return RunMangaList(responder, b, mode, targetID, name, icon, page)
}
