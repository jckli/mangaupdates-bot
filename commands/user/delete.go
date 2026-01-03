package user

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunDelete(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return responder.Error(err.Error())
	}

	embed := common.StandardEmbed("Delete User Profile?", "Are you sure you want to delete your personal profile?\n\n**This action cannot be undone.**\nYou will lose your watchlist and stop receiving notifications.")
	embed.Color = common.ColorError

	buttons := common.CreateConfirmButtons("/user_delete_confirm/yes", "/user_delete_confirm/no")

	return responder.Respond(embed, buttons)
}

func HandleDeleteConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardUser(b, e.User().ID.String()); err != nil {
		return err
	}

	action := e.Vars["action"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Your profile was not deleted."),
				},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	err := b.ApiClient.DeleteUser(e.User().ID.String())
	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = common.ColorError
		_, _ = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errEmbed},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("User Removed", "Your personal profile has been removed."),
			},
			Components: &[]discord.ContainerComponent{},
		})
	return err
}
