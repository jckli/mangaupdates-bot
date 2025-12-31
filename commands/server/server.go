package server

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
)

var ServerCommand = discord.SlashCommandCreate{
	Name:                     "server",
	Description:              "Manage the server's manga tracking list",
	DMPermission:             json.Ptr(false),
	DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageGuild),
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "Show all manga tracked by this server",
		},
	},
}
