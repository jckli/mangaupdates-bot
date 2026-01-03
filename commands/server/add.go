package server

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func AddHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	query := e.SlashCommandInteractionData().String("title")
	return manga.RunAddEntry(responder, b, "server", query)
}
