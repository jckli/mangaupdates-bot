package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RemoveHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	query := e.SlashCommandInteractionData().String("title")
	targetID := e.User().ID.String()
	responder := &common.CommandResponder{Event: e}
	return manga.RunRemoveEntry(responder, b, "user", targetID, query)
}
