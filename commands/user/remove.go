package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RemoveHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return responder.Error(err.Error())
	}

	query := e.SlashCommandInteractionData().String("title")
	targetID := e.User().ID.String()
	return manga.RunRemoveEntry(responder, b, "user", targetID, query)
}
