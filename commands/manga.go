package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

	dbUser, err := utils.DbGetUser(b, userId)
	if dbUser == nil || err != nil {
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

	dbServer, err := utils.DbGetServer(b, int64(serverId))
	if dbServer == nil || err != nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
		)
	}

	searchResults, searchResultsFormatted := searchResultsEmbed(b, "Add Manga", title)
	if searchResultsFormatted == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(searchResults).Build(),
		)
	}
	dropdownSearchResults := dropdownSearchResultsComponents("manga", "add", searchResultsFormatted)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.MessageCreate{
			Embeds:     []discord.Embed{searchResults},
			Components: dropdownSearchResults,
		},
	)
}
