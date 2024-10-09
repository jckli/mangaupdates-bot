package commands

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

var mangaCommand = discord.SlashCommandCreate{
	Name:        "manga",
	Description: "Interact with manga",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a manga to your list",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "title",
					Description: "The title of the manga",
					Required:    true,
				},
			},
		},
	},
}

func mangaAddHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	title := e.SlashCommandInteractionData().String("title")

	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return userMangaAddHandler(adapter, b, title)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Add Manga",
						"Would you like this manga's chapter updates to send to this server or your DMs?",
					),
				},
				Components: selectServerOrUserComponents("manga", "add", title),
			},
		)
	}
}

func userMangaAddHandler(
	e EventHandler,
	b *mubot.Bot,
	title string,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to check if user exists (userMangaAddHandler): %s", err.Error()),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorTechnicalErrorEmbed()).Build(),
		)
	}
	if !exists {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
	)
}

func serverMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot, title string) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (serverMangaAddHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
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

	searchResults, searchResultsFormatted := searchResultsEmbed(b, "Add Manga", title)
	if searchResultsFormatted == nil {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{searchResults},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	dropdownSearchResults := dropdownSearchResultsComponents(
		"manga",
		"add",
		"server",
		searchResultsFormatted,
	)

	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{searchResults},
			Components: &dropdownSearchResults,
		},
	)

	/* Example of handling a custom error
	var customErr rest.Error
	if errors.As(err, &customErr) {
		fmt.Println(string(customErr.RsBody))
	}
	*/

	return err
}

func searchMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot, mode string) error {
	mangaId := e.StringSelectMenuInteractionData().Values[0]

	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Failed to parse manga ID (searchMangaAddHandler): %s", mangaId))
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	components := selectConfirmMangaComponents("manga", "add", mode, mangaId)
	err = e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{confirmMangaEmbed(b, "Add Manga", intMangaId)},
			Components: &components,
		},
	)

	var customErr rest.Error
	if errors.As(err, &customErr) {
		fmt.Println(string(customErr.RsBody))
	}

	return err
}

func cancelMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot) error {
	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{cancelMangaEmbed("Add Manga")},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func confirmMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot, mode, mangaId string) error {
	intMangaId, err := strconv.ParseInt(mangaId, 10, 64)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to parse manga ID (confirmMangaAddHandler): %s", mangaId),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{errorTechnicalErrorEmbed()},
			},
		)
	}

	if mode == "user" {
		return userConfirmMangaAddHandler(e, b, intMangaId)
	} else {
		return serverConfirmMangaAddHandler(e, b, intMangaId)
	}
}

func userConfirmMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot, mangaId int64) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckMangaExists(b, userId, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to check if manga exists in user: %s", err.Error()),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	if exists {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{mangaExistsEmbed("Add Manga")},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	seriesinfo, err := utils.MuGetSeriesInfo(b, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to get series info (userConfirmMangaAddHandler): %s", err.Error()),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	mangaEntry := utils.MDbManga{
		Id:    mangaId,
		Title: seriesinfo.Title,
	}

	err = utils.DbUserAddManga(b, userId, mangaEntry)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to add manga to user: %s", err.Error()),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{successMangaAddEmbed("Add Manga", seriesinfo.Title)},
			Components: &[]discord.ContainerComponent{},
		},
	)
}

func serverConfirmMangaAddHandler(e *handler.ComponentEvent, b *mubot.Bot, mangaId int64) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckMangaExists(b, serverId, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if manga exists in server: %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	if exists {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{mangaExistsEmbed("Add Manga")},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	seriesinfo, err := utils.MuGetSeriesInfo(b, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to get series info (serverConfirmMangaAddHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	mangaEntry := utils.MDbManga{
		Id:    mangaId,
		Title: seriesinfo.Title,
	}

	err = utils.DbServerAddManga(b, serverId, mangaEntry)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to add manga to server: %s", err.Error()),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return e.UpdateMessage(
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{successMangaAddEmbed("Add Manga", seriesinfo.Title)},
			Components: &[]discord.ContainerComponent{},
		},
	)
}
