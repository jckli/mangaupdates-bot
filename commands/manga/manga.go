package manga

import (
	"github.com/disgoorg/disgo/discord"
)

var MangaCommand = discord.SlashCommandCreate{
	Name:        "manga",
	Description: "Interact with your manga list",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add",
			Description: "Add a manga to your list",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "title",
					Description: "The title of the manga",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Remove a manga from your list",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "list",
			Description: "List all manga you have added",
		},
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "scanlator",
			Description: "Modify a manga's scanlator",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Add a manga's scanlator",
				},
				{
					Name:        "remove",
					Description: "Remove a manga's scanlator",
				},
			},
		},
	},
}
