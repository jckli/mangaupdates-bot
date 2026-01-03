package user

import (
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func SetupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}

	err := b.ApiClient.SetupUser(e.User().ID.String(), e.User().EffectiveName())
	if err != nil {
		return responder.Error("Failed to setup: " + err.Error())
	}

	desc := "Your personal manga list has been created!\nYou will receive DM notifications for manga you add."
	return responder.Respond(common.StandardEmbed("Setup Complete", desc), nil)
}
