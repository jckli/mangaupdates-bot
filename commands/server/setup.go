package server

import (
	"fmt"

	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func SetupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardAdminOnly(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	if e.GuildID() == nil {
		return responder.Error("This command can only be used in a server.")
	}

	data := e.SlashCommandInteractionData()
	channel := data.Channel("channel")

	serverName := "Unknown Server"
	if guild, ok := e.Guild(); ok {
		serverName = guild.Name
	}

	err := b.ApiClient.SetupServer(e.GuildID().String(), serverName, channel.ID.String())
	if err != nil {
		return responder.Error("Failed to setup: " + err.Error())
	}

	desc := fmt.Sprintf("Server **%s** has been successfully initialized.\nManga updates will be posted in <#%s>.", serverName, channel.ID.String())
	return responder.Respond(common.StandardEmbed("Setup Complete", desc), nil)
}
