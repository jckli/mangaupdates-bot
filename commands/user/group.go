package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func SetGroupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	query := data.String("title")
	group := data.String("group")

	if e.GuildID() == nil {
		return nil
	}
	targetID := e.GuildID().String()
	return manga.RunSetGroupEntry(responder, b, "user", targetID, query, group)
}

func RemoveGroupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	query := data.String("title")

	if e.GuildID() == nil {
		return nil
	}
	targetID := e.GuildID().String()
	return manga.RunSetGroupEntry(responder, b, "user", targetID, query, "0")
}
