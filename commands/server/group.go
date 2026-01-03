package server

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func SetGroupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	query := data.String("title")
	group := data.String("group")

	if e.GuildID() == nil {
		return nil
	}
	targetID := e.GuildID().String()
	return manga.RunSetGroupEntry(responder, b, "server", targetID, query, group)
}

func RemoveGroupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	query := data.String("title")

	if e.GuildID() == nil {
		return nil
	}
	targetID := e.GuildID().String()
	return manga.RunSetGroupEntry(responder, b, "server", targetID, query, "0")
}
