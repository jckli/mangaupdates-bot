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
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a manga to this server's tracking list",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "title",
					Description:  "The title of the manga",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	},
}
