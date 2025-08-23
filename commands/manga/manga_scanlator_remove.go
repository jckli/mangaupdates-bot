package manga

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaScanlatorRemoveHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaScanlatorRemoveUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Remove Scanlator",
						"Would you like to remove a manga's scanlator from this server or your DMs?",
					),
				},
				Components: selectServerOrUserNestedComponents("manga", "scanlator", "remove", ""),
			},
		)
	}
}

func MangaScanlatorRemoveUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if user exists (MangaScanlatorRemoveUserHandler): %s",
				err.Error(),
			),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}
	if !exists {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
		)
	}

	return MangaScanlatorRemoveUserSearchHandler(e, b, 1)
}

func MangaScanlatorRemoveServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaScanlatorRemoveServerHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	if !exists {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorMangaSetupNeededEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return MangaScanlatorRemoveServerSearchHandler(e, b, 1)
}

func MangaScanlatorRemoveUserSearchHandler(
	e EventHandler,
	b *mubot.Bot,
	page int,
) error {
	userId := int64(e.User().ID)

	user, err := utils.DbGetUser(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get user (MangaScanlatorRemoveUserHandler): %s",
				err.Error(),
			),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}

	parsed := parsePaginationMangaList(user.Manga, page)

	searchResults, searchResultsFormatted := dbMangaSearchResultsEmbed(
		"Remove Scanlator",
		parsed.MangaList,
		page,
	)
	if searchResultsFormatted == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(searchResults).Build(),
		)
	}
	dropdownSearchResults := dropdownDbMangaSearchResultsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"user",
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds:     []discord.Embed{searchResults},
				Components: dropdownSearchResults,
			},
		)
	}

	pagination := paginationMangaSearchResultsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"user",
		parsed,
	)
	components := append(dropdownSearchResults, pagination...)
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.MessageCreate{
			Embeds:     []discord.Embed{searchResults},
			Components: components,
		},
	)
}

func MangaScanlatorRemoveServerSearchHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
) error {
	serverId := int64(*e.GuildID())

	server, err := utils.DbGetServer(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get server (MangaScanlatorRemoveServerHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	parsed := parsePaginationMangaList(server.Manga, page)

	searchResults, searchResultsFormatted := dbMangaSearchResultsEmbed(
		"Remove Scanlator",
		parsed.MangaList,
		page,
	)
	if searchResultsFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	dropdownSearchResults := dropdownDbMangaSearchResultsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"server",
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults},
				Components: &dropdownSearchResults,
			},
		)
	}

	pagination := paginationMangaSearchResultsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"server",
		parsed,
	)
	components := append(dropdownSearchResults, pagination...)
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults},
			Components: &components,
		},
	)
}

func MangaScanlatorRemoveUserGroupHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
	mangaId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga ID (MangaScanlatorRemoveUserGroupHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	userId := int64(e.User().ID)

	manga, err := utils.DbUserGetManga(b, userId, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga (MangaScanlatorRemoveUserGroupHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	parsed := parsePaginationDbScanlators(manga.Scanlators, page)

	searchResults, searchResultsFormatted := selectDbScanlatorsEmbed(
		"Remove Scanlator",
		parsed.GroupList,
	)
	if searchResultsFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	dropdownSearchResults := dropdownSeriesGroupsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"user",
		mangaId,
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &dropdownSearchResults,
			},
		)
	}

	pagination := paginationDbScanlatorsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"user",
		mangaId,
		parsed,
	)

	components := append(dropdownSearchResults, pagination...)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults.Build()},
			Components: &components,
		},
	)

	return err
}

func MangaScanlatorRemoveServerGroupHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
	mangaId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga ID (MangaScanlatorRemoveServerGroupHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	serverId := int64(*e.GuildID())

	manga, err := utils.DbServerGetManga(b, serverId, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga (MangaScanlatorRemoveServerGroupHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	parsed := parsePaginationDbScanlators(manga.Scanlators, page)

	searchResults, searchResultsFormatted := selectDbScanlatorsEmbed(
		"Remove Scanlator",
		parsed.GroupList,
	)
	if searchResultsFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	dropdownSearchResults := dropdownSeriesGroupsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"server",
		mangaId,
		searchResultsFormatted,
	)

	if !parsed.Pagination {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults.Build()},
				Components: &dropdownSearchResults,
			},
		)
	}

	pagination := paginationDbScanlatorsNestedComponents(
		"manga",
		"scanlator",
		"remove",
		"server",
		mangaId,
		parsed,
	)

	components := append(dropdownSearchResults, pagination...)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults.Build()},
			Components: &components,
		},
	)

	var customErr rest.Error
	if errors.As(err, &customErr) {
		fmt.Println(string(customErr.RsBody))
	}

	return err
}

func MangaScanlatorRemoveGroupSelectHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode,
	mangaId string,
) error {
	groupId := e.StringSelectMenuInteractionData().Values[0]

	intGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse group ID (MangaScanlatorRemoveGroupSelectHandler): %s",
				groupId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	components := selectConfirmGroupComponents(
		"manga",
		"scanlator",
		"remove",
		mode,
		mangaId,
		groupId,
	)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{confirmGroupEmbed(b, "Remove Scanlator", intGroupId)},
			Components: &components,
		},
	)

	return err
}

func MangaScanlatorRemoveGroupCancelHandler(
	e *handler.ComponentEvent,
) error {
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{cancelEmbed("Remove Scanlator")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func MangaScanlatorRemoveServerGroupConfirmHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mangaId,
	groupId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga id (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	intGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse group id (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	groupInfo, err := utils.MuGetGroupInfo(b, intGroupId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get group info (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	serverId := int64(*e.GuildID())

	ok, err := utils.DbServerRemoveGroup(b, serverId, intMangaId, intGroupId)
	if !ok || err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to remove scanlator from server (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	manga, err := utils.DbServerGetManga(b, serverId, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga after removing scanlator from server (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				successMangaRemoveScanlatorEmbed("Remove Scanlator", manga.Title, groupInfo.Name),
			},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func MangaScanlatorRemoveUserGroupConfirmHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mangaId,
	groupId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga id (MangaScanlatorRemoveUserGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	intGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse group id (MangaScanlatorRemoveUserGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	groupInfo, err := utils.MuGetGroupInfo(b, intGroupId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get group info (MangaScanlatorRemoveUserGroupConfirmHandler): %s",
				mangaId,
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	userId := int64(e.User().ID)

	ok, err := utils.DbUserRemoveGroup(b, userId, intMangaId, intGroupId)
	if !ok || err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to remove scanlator from server (MangaScanlatorRemoveUserGroupConfirmHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	manga, err := utils.DbUserGetManga(b, userId, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga after removing scanlator from server (MangaScanlatorRemoveServerGroupConfirmHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				successMangaRemoveScanlatorEmbed("Remove Scanlator", manga.Title, groupInfo.Name),
			},
			Components: &[]discord.ContainerComponent{},
		},
	)
}
