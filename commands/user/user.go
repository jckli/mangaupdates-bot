package user

import (
	"github.com/disgoorg/disgo/discord"
)

var UserCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Manage your personal manga tracking list",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "Show your personal tracked manga",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a manga to your personal tracking list",
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
