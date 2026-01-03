package server

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RunDelete(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	if e.GuildID() == nil {
		return responder.Error("This command can only be used in a server.")
	}

	embed := common.StandardEmbed("Delete Server Data?", "Are you sure you want to delete this server's configuration?\n\n**This action cannot be undone.**\nAll tracked manga and settings will be permanently removed.")
	embed.Color = common.ColorError
	buttons := common.CreateConfirmButtons("/server_delete_confirm/yes", "/server_delete_confirm/no")

	return responder.Respond(embed, buttons)
}

func HandleDeleteConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return err
	}

	action := e.Vars["action"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Server data was not deleted."),
				},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	if e.GuildID() == nil {
		return nil
	}

	err := b.ApiClient.DeleteServer(e.GuildID().String())
	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = common.ColorError
		_, _ = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{Embeds: &[]discord.Embed{errEmbed}})
		return err
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("Server Removed", "This server has been removed from the database."),
			},
			Components: &[]discord.ContainerComponent{},
		})
	return err
}
