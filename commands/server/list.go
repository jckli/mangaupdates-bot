package server

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func ListHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if e.GuildID() == nil {
		return nil
	}

	guild, _ := e.Guild()

	icon := ""
	if i := guild.IconURL(); i != nil {
		icon = *i
	}

	responder := &common.CommandResponder{Event: e}

	return manga.RunMangaList(
		responder,
		b,
		"server",
		e.GuildID().String(),
		guild.Name,
		icon,
		1,
	)
}
