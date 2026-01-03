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
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Remove a manga from this server's tracking list",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "title",
					Description:  "The title of the manga",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "setgroup",
			Description: "Filter a manga to specific scanlation groups",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "title",
					Description:  "The title of the manga to modify",
					Required:     false,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionString{
					Name:         "group",
					Description:  "The group to filter by",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
	},
}
