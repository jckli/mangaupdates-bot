package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func ListHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return responder.Error(err.Error())
	}

	return manga.RunMangaList(
		responder,
		b,
		"user",
		e.User().ID.String(),
		e.User().EffectiveName(),
		e.User().EffectiveAvatarURL(),
		1,
	)
}
