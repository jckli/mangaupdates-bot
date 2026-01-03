package search

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var SearchCommand = discord.SlashCommandCreate{
	Name:        "search",
	Description: "Search for a manga on MangaUpdates",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "title",
			Description:  "The title of the manga",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func SearchHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	query := e.SlashCommandInteractionData().String("title")
	responder := &common.CommandResponder{Event: e}

	return manga.RunSearchEntry(responder, b, query)
}
