package server

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func ChannelSetHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	channel := data.Channel("channel")

	desc := fmt.Sprintf("Are you sure you want to change the notification channel to <#%s>?\n\nAll future manga updates will be posted there.", channel.ID.String())
	embed := common.StandardEmbed("Confirm Channel Update", desc)

	confirmPath := fmt.Sprintf("/server_channel_confirm/%s/yes", channel.ID)
	cancelPath := fmt.Sprintf("/server_channel_confirm/%s/no", channel.ID)

	buttons := common.CreateConfirmButtons(confirmPath, cancelPath)

	return responder.Respond(embed, buttons)
}

func HandleChannelConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return err
	}

	action := e.Vars["action"]
	channelID := e.Vars["channel_id"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Notification channel was not changed."),
				},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	err := b.ApiClient.UpdateServerChannel(e.GuildID().String(), channelID)
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
				common.StandardEmbed("Success", fmt.Sprintf("Manga updates will now be posted in <#%s>.", channelID)),
			},
			Components: &[]discord.ContainerComponent{},
		})
	return err
}
