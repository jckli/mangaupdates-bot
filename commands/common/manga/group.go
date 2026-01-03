package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func RunSetGroupEntry(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	targetID string,
	query string,
	groupQuery string,
) error {
	mangaID, errM := strconv.ParseInt(query, 10, 64)

	if errM == nil && groupQuery != "" {
		groupID, errG := strconv.ParseInt(groupQuery, 10, 64)
		if errG == nil {
			return sendSetGroupConfirmation(r, b, endpoint, mangaID, groupID)
		}
	}

	if errM == nil && groupQuery != "0" {
		return RunSetGroupGroupMenu(r, b, endpoint, mangaID, 1)
	}

	prefix := "setgroup_manga_select"
	nav := "/setgroup_manga_nav"
	title := "Select Manga"
	placeholder := "Select manga..."

	if groupQuery == "0" {
		prefix = "group_remove_manga_select"
		nav = "/group_remove_manga_nav"
		title = "Select Manga to Clear Filter"
		placeholder = "Select manga to clear filter..."
	}

	embed, components, err := GenerateWatchlistMenu(b, WatchlistMenuConfig{
		Endpoint:            endpoint,
		TargetID:            targetID,
		Query:               query,
		Page:                1,
		SelectIDPrefix:      prefix,
		NavIDPrefix:         nav,
		Title:               title,
		DropdownPlaceholder: placeholder,
	})
	if err != nil {
		return r.Error(err.Error())
	}
	return r.Respond(embed, components)
}

func RunSetGroupGroupMenu(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	mangaID int64,
	page int,
) error {
	groups, err := b.ApiClient.GetMangaGroups(mangaID)
	if err != nil {
		return r.Error("Failed to fetch groups.")
	}

	var options []utils.MangaSearchResult
	var displayLines []string
	options = append(options, utils.MangaSearchResult{ID: 0, Title: "All Groups (Clear Filter)"})
	displayLines = append(displayLines, "• **All Groups (Clear Filter)**")

	for _, g := range groups {
		options = append(options, utils.MangaSearchResult{ID: g.ID, Title: g.Name})
		displayLines = append(displayLines, fmt.Sprintf("• %s", g.Name))
	}

	slicedOptions, totalPages := common.GetPageSlice(options, page)
	slicedLines, _ := common.GetPageSlice(displayLines, page)

	customSelectID := fmt.Sprintf("setgroup_group_select/%s/%d", endpoint, mangaID)
	dropdown := common.GenerateSearchDropdown(customSelectID, "Select scanlation group...", slicedOptions)

	btnPrefix := fmt.Sprintf("/setgroup_group_nav/%s/%d", endpoint, mangaID)
	buttons := common.GeneratePaginationButtons(btnPrefix, page, totalPages)

	allComponents := append(dropdown, buttons...)

	var description string
	for _, line := range slicedLines {
		description += line + "\n"
	}

	header := fmt.Sprintf("Found %d active groups for this series.", len(groups))
	fullDescription := header + "\nPlease select one from the dropdown below:\n\n" + description

	embed := common.StandardEmbed("Select Group", fullDescription)
	embed.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Page %d/%d", page, totalPages)}

	return r.Respond(embed, allComponents)
}

func sendSetGroupConfirmation(
	r common.Responder,
	b *mubot.Bot,
	endpoint string,
	mangaID int64,
	groupID int64,
) error {
	mangaDetails, _ := b.ApiClient.GetMangaDetails(mangaID)
	mangaTitle := "Unknown Manga"
	if mangaDetails != nil {
		mangaTitle = mangaDetails.Title
	}

	var groupDetails *utils.GroupDetails
	if groupID != 0 {
		var err error
		groupDetails, err = b.ApiClient.GetGroupDetails(groupID)
		if err != nil {
			groupDetails = &utils.GroupDetails{Name: "Unknown Group", GroupID: groupID}
		}
	}

	embed := common.GenerateGroupConfirmationEmbed(groupDetails, mangaTitle)

	prefix := fmt.Sprintf("/setgroup_confirm/%s/%d/%d", endpoint, mangaID, groupID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	return r.Respond(embed, buttons)
}

func HandleSetGroupMangaSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	groups, err := b.ApiClient.GetMangaGroups(mangaID)
	if err != nil {
		return err
	}

	var options []utils.MangaSearchResult
	var displayLines []string
	options = append(options, utils.MangaSearchResult{ID: 0, Title: "All Groups (Clear Filter)"})
	displayLines = append(displayLines, "• **All Groups (Clear Filter)**")

	for _, g := range groups {
		options = append(options, utils.MangaSearchResult{ID: g.ID, Title: g.Name})
		displayLines = append(displayLines, fmt.Sprintf("• %s", g.Name))
	}

	slicedOptions, totalPages := common.GetPageSlice(options, 1)
	slicedLines, _ := common.GetPageSlice(displayLines, 1)

	customSelectID := fmt.Sprintf("setgroup_group_select/%s/%d", endpoint, mangaID)
	dropdown := common.GenerateSearchDropdown(customSelectID, "Select scanlation group...", slicedOptions)

	btnPrefix := fmt.Sprintf("/setgroup_group_nav/%s/%d", endpoint, mangaID)
	buttons := common.GeneratePaginationButtons(btnPrefix, 1, totalPages)

	allComponents := append(dropdown, buttons...)

	var description string
	for _, line := range slicedLines {
		description += line + "\n"
	}

	header := fmt.Sprintf("Found %d active groups for this series.", len(groups))
	fullDescription := header + "\nPlease select one from the dropdown below:\n\n" + description

	embed := common.StandardEmbed("Select Group", fullDescription)
	embed.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Page %d/%d", 1, totalPages)}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &allComponents,
		})
	return err
}

func HandleSetGroupGroupSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)

	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	groupID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	mangaDetails, _ := b.ApiClient.GetMangaDetails(mangaID)
	mangaTitle := "Unknown Manga"
	if mangaDetails != nil {
		mangaTitle = mangaDetails.Title
	}

	var groupDetails *utils.GroupDetails
	if groupID != 0 {
		var err error
		groupDetails, err = b.ApiClient.GetGroupDetails(groupID)
		if err != nil {
			groupDetails = &utils.GroupDetails{Name: "Unknown Group", GroupID: groupID}
		}
	}

	embed := common.GenerateGroupConfirmationEmbed(groupDetails, mangaTitle)

	prefix := fmt.Sprintf("/setgroup_confirm/%s/%d/%d", endpoint, mangaID, groupID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &buttons,
		})
	return err
}

func HandleSetGroupConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	groupID, _ := strconv.ParseInt(e.Vars["group_id"], 10, 64)
	action := e.Vars["action"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Group filter was not changed."),
				},
				Components: &[]discord.ContainerComponent{},
			},
		)
		return err
	}

	targetID := e.User().ID.String()
	if endpoint == "server" {
		if e.GuildID() == nil {
			return nil
		}
		targetID = e.GuildID().String()
	}

	groupName := "All"
	if groupID != 0 {
		found := false
		if len(e.Message.Embeds) > 0 {
			for _, field := range e.Message.Embeds[0].Fields {
				if field.Name == "Group Name" {
					groupName = field.Value
					found = true
					break
				}
			}
		}
		if !found {
			g, err := b.ApiClient.GetGroupDetails(groupID)
			if err == nil {
				groupName = g.Name
			}
		}
	}

	err := b.ApiClient.UpdateMangaGroup(endpoint, targetID, mangaID, groupName, groupID)

	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = 0xFF0000

		_, _ = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{Embeds: &[]discord.Embed{errEmbed}})
		return err
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("Success", fmt.Sprintf("Group filter set to **%s**.", groupName)),
			},
			Components: &[]discord.ContainerComponent{},
		},
	)
	return err
}

func HandleSetGroupMangaPagination(e *handler.ComponentEvent, b *mubot.Bot) error {
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
		SelectIDPrefix:      "setgroup_manga_select",
		NavIDPrefix:         "/setgroup_manga_nav",
		Title:               "Select Manga",
		DropdownPlaceholder: "Select manga...",
	})
}

func HandleSetGroupGroupPagination(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	mangaID, _ := strconv.ParseInt(e.Vars["manga_id"], 10, 64)
	page, _ := strconv.Atoi(e.Vars["page"])

	groups, err := b.ApiClient.GetMangaGroups(mangaID)
	if err != nil {
		return err
	}

	var options []utils.MangaSearchResult
	var displayLines []string
	options = append(options, utils.MangaSearchResult{ID: 0, Title: "All Groups (Clear Filter)"})
	displayLines = append(displayLines, "• **All Groups (Clear Filter)**")

	for _, g := range groups {
		options = append(options, utils.MangaSearchResult{ID: g.ID, Title: g.Name})
		displayLines = append(displayLines, fmt.Sprintf("• %s", g.Name))
	}

	slicedOptions, totalPages := common.GetPageSlice(options, page)
	slicedLines, _ := common.GetPageSlice(displayLines, page)

	customSelectID := fmt.Sprintf("setgroup_group_select/%s/%d", endpoint, mangaID)
	dropdown := common.GenerateSearchDropdown(customSelectID, "Select scanlation group...", slicedOptions)

	btnPrefix := fmt.Sprintf("/setgroup_group_nav/%s/%d", endpoint, mangaID)
	buttons := common.GeneratePaginationButtons(btnPrefix, page, totalPages)

	allComponents := append(dropdown, buttons...)

	var description string
	for _, line := range slicedLines {
		description += line + "\n"
	}

	header := fmt.Sprintf("Found %d active groups for this series.", len(groups))
	fullDescription := header + "\nPlease select one from the dropdown below:\n\n" + description

	embed := common.StandardEmbed("Select Group", fullDescription)
	embed.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Page %d/%d", page, totalPages)}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &allComponents,
		})
	return err
}

func HandleGroupRemoveMangaSelection(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardWidget(e, b, true); err != nil {
		return err
	}

	endpoint := e.Vars["mode"]
	if len(e.StringSelectMenuInteractionData().Values) == 0 {
		return nil
	}
	mangaID, _ := strconv.ParseInt(e.StringSelectMenuInteractionData().Values[0], 10, 64)

	mangaDetails, _ := b.ApiClient.GetMangaDetails(mangaID)
	mangaTitle := "Unknown Manga"
	if mangaDetails != nil {
		mangaTitle = mangaDetails.Title
	}

	embed := common.GenerateGroupConfirmationEmbed(nil, mangaTitle)

	prefix := fmt.Sprintf("/setgroup_confirm/%s/%d/0", endpoint, mangaID)
	buttons := common.CreateConfirmButtons(prefix+"/yes", prefix+"/no")

	_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &buttons,
		})
	return err
}

func HandleGroupRemoveMangaPagination(e *handler.ComponentEvent, b *mubot.Bot) error {
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
		SelectIDPrefix:      "group_remove_manga_select",
		NavIDPrefix:         "/group_remove_manga_nav",
		Title:               "Select Manga to Clear Filter",
		DropdownPlaceholder: "Select manga to clear filter...",
	})
}
