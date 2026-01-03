package user

import (
	"github.com/disgoorg/disgo/discord"
)

var UserCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Manage your personal manga tracking list",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "setup",
			Description: "Initialize your personal manga watchlist (DM notifications)",
		},
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
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Remove a manga from your personal tracking list",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "title",
					Description:  "The title of the manga",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "group",
			Description: "Manage scanlation group filters",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "set",
					Description: "Set a specific scanlation group filter for a manga",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:         "title",
							Description:  "The manga to modify",
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
				{
					Name:        "remove",
					Description: "Clear the scanlation group filter (Allow all groups)",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:         "title",
							Description:  "The manga to modify",
							Required:     false,
							Autocomplete: true,
						},
					},
				},
			},
		},
	},
}
