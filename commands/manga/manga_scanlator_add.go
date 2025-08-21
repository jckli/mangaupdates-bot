package manga

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaScanlatorAddHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaScanlatorAddUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Add Scanlator",
						"Would you like to add a manga's scanlator for this server or your DMs?",
					),
				},
				Components: selectServerOrUserNestedComponents("manga", "scanlator", "add", ""),
			},
		)
	}
}

func MangaScanlatorAddUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if user exists (MangaScanlatorAddUserHandler): %s",
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

	return MangaScanlatorAddUserSearchHandler(e, b, 1)
}

func MangaScanlatorAddServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaScanlatorAddServerHandler): %s",
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

	return MangaScanlatorAddServerSearchHandler(e, b, 1)
}

func MangaScanlatorAddUserSearchHandler(
	e EventHandler,
	b *mubot.Bot,
	page int,
) error {
	userId := int64(e.User().ID)

	user, err := utils.DbGetUser(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get user (MangaScanlatorAddUserHandler): %s",
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
		"Add Scanlator",
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
		"add",
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
		"add",
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

func MangaScanlatorAddServerSearchHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
) error {
	serverId := int64(*e.GuildID())

	server, err := utils.DbGetServer(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get server (MangaScanlatorAddServerHandler): %s",
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
		"Add Scanlator",
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
		"add",
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
		"add",
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

func MangaScanlatorAddMangaSelectHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode string,
) error {
	mangaId := e.StringSelectMenuInteractionData().Values[0]

	return MangaScanlatorAddGroupHandler(e, b, 1, mode, mangaId)
}

func MangaScanlatorAddGroupHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	page int,
	mode,
	mangaId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga ID (MangaScanlatorAddSelectGroupHandler): %s",
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
	groups, err := utils.MuGetSeriesGroups(b, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get manga (MangaScanlatorAddSelectGroupHandler): %s",
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

	parsed := parsePaginationSeriesGroups(groups.GroupList, page)

	searchResults, searchResultsFormatted := selectSeriesGroupEmbed(
		"Add Scanlator",
		parsed.GroupList,
	)
	if searchResultsFormatted == nil {
		fmt.Println("no search results formatted")
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
		"add",
		mode,
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

	pagination := paginationGroupListNestedComponents(
		"manga",
		"scanlator",
		"add",
		mode,
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

func MangaScanlatorAddGroupSelectHandler(
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
				"Failed to parse group ID (MangaScanlatorAddGroupSelectHandler): %s",
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

	components := selectConfirmGroupComponents("manga", "scanlator", "add", mode, mangaId, groupId)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{confirmGroupEmbed(b, "Add Scanlator", intGroupId)},
			Components: &components,
		},
	)

	return err
}

func MangaScanlatorAddGroupCancelHandler(
	e *handler.ComponentEvent,
) error {
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{cancelEmbed("Add Scanlator")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func MangaScanlatorAddGroupConfirmHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	mode,
	mangaId,
	groupId string,
) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to parse manga id (MangaScanlatorAddGroupConfirmHandler): %s",
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
				"Failed to parse group id (MangaScanlatorAddGroupConfirmHandler): %s",
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
				"Failed to get group info (MangaScanlatorAddGroupConfirmHandler): %s",
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

	if mode == "user" {
		userId := int64(e.User().ID)
		exists, err := utils.DbUserCheckGroupExists(b, userId, intMangaId, intGroupId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to check group exists from user (MangaScanlatorAddGroupConfirmHandler): %s",
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
		if exists {
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{groupExistsEmbed("Add Scanlator")},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}
		err = utils.DbUserAddGroup(b, userId, intMangaId, intGroupId, groupInfo.Name)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to set scanlator from user (MangaScanlatorAddGroupConfirmHandler): %s",
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
	} else {
		serverId := int64(*e.GuildID())
		exists, err := utils.DbServerCheckGroupExists(b, serverId, intMangaId, intGroupId)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to check group exists from server (MangaScanlatorAddGroupConfirmHandler): %s",
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
		if exists {
			return e.UpdateMessage(
				discord.MessageUpdate{
					Embeds:     &[]discord.Embed{groupExistsEmbed("Add Scanlator")},
					Components: &[]discord.ContainerComponent{},
				},
			)
		}

		err = utils.DbServerAddGroup(b, serverId, intMangaId, intGroupId, groupInfo.Name)
		if err != nil {
			b.Logger.Error(
				fmt.Sprintf(
					"Failed to set scanlator from server (MangaScanlatorAddGroupConfirmHandler): %s",
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
	}

	seriesInfo, err := utils.MuGetSeriesInfo(b, intMangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get series info (MangaScanlatorAddGroupConfirmHandler): %s",
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

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				successMangaSetScanlatorEmbed("Add Scanlator", seriesInfo.Title, groupInfo.Name),
			},
			Components: &[]discord.ContainerComponent{},
		},
	)
}
