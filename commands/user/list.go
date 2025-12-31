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
	},
}
