package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func AddHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	query := e.SlashCommandInteractionData().String("title")
	responder := &common.CommandResponder{Event: e}
	return manga.RunAddEntry(responder, b, "user", query)
}
