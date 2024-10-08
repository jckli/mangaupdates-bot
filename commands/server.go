package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

var serverCommand = discord.SlashCommandCreate{
	Name:        "server",
	Description: "Interact with the server",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "setup",
			Description: "Sets up the server for manga updates",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionChannel{
					Name:        "channel",
					Description: "The channel to send manga updates to",
					Required:    true,
				},
			},
		},
	},
}

func serverSetupHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	channel := e.SlashCommandInteractionData().Channel("channel")
	serverId := int64(*e.GuildID())
	server, notInGuild := e.Guild()
	if notInGuild {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorTechnicalErrorEmbed()).Build(),
		)
	}

	dbServer, err := utils.DbGetServer(b, int64(serverId))
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error getting server in serverSetupHandler: %s", err.Error()))
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorTechnicalErrorEmbed()).Build(),
		)
	}
	if !e.Member().Permissions.Has(discord.PermissionAdministrator) {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorNoPermissionsEmbed()).Build(),
		)
	}
	if dbServer != nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorServerAlreadySetupEmbed()).Build(),
		)
	}

	err = utils.DbAddServer(b, server.Name, int64(serverId), int64(channel.ID))
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error adding server in serverSetupHandler: %s", err.Error()))
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorTechnicalErrorEmbed()).Build(),
		)
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("Server Setup").
		SetDescription("Great! This server is now setup for manga updates. You can add manga now using the `/manga add` command.").
		SetColor(0x3083e3).
		Build()
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(embed).Build(),
	)
}
