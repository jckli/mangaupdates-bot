package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/commands/search"
	"github.com/jckli/mangaupdates-bot/commands/server"
	"github.com/jckli/mangaupdates-bot/commands/user"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var CommandList = []discord.ApplicationCommandCreate{
	helpCommand,
	pingCommand,
	infoCommand,
	alertCommand,
	server.ServerCommand,
	user.UserCommand,
	search.SearchCommand,
}

func CommandHandlers(b *mubot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/help", func(e *handler.CommandEvent) error {
		return HelpHandler(e, b)
	})
	h.Command("/mangaupdates", func(e *handler.CommandEvent) error {
		return InfoHandler(e, b)
	})
	h.Command("/alert", func(e *handler.CommandEvent) error {
		return AlertHandler(e, b)
	})
	h.Command("/ping", PingHandler)
	h.Route("/search", func(h handler.Router) {
		h.Command("/", func(e *handler.CommandEvent) error {
			return search.SearchHandler(e, b)
		})
		h.Autocomplete("/", func(e *handler.AutocompleteEvent) error {
			// we can reuse the add autocomplete cuz its the same functionally
			return manga.HandleAddAutocomplete(e, b, "title")
		})
	})
	h.Route("/server", func(h handler.Router) {
		h.Command("/setup", func(e *handler.CommandEvent) error {
			return server.SetupHandler(e, b)
		})
		h.Command("/delete", func(e *handler.CommandEvent) error {
			return server.RunDelete(e, b)
		})
		h.Command("/list", func(e *handler.CommandEvent) error {
			return server.ListHandler(e, b)
		})
		h.Command("/add", func(e *handler.CommandEvent) error {
			return server.AddHandler(e, b)
		})
		h.Autocomplete("/add", func(e *handler.AutocompleteEvent) error {
			return manga.HandleAddAutocomplete(e, b, "title")
		})
		h.Command("/remove", func(e *handler.CommandEvent) error {
			return server.RemoveHandler(e, b)
		})
		h.Autocomplete("/remove", func(e *handler.AutocompleteEvent) error {
			if e.GuildID() == nil {
				return e.AutocompleteResult(nil)
			}
			return manga.HandleWatchlistAutocomplete(e, b, "server", e.GuildID().String(), "title")
		})
		h.Route("/group", func(h handler.Router) {
			h.Command("/set", func(e *handler.CommandEvent) error {
				return server.SetGroupHandler(e, b)
			})
			h.Autocomplete("/set", func(e *handler.AutocompleteEvent) error {
				if e.GuildID() == nil {
					return e.AutocompleteResult(nil)
				}
				return manga.HandleSetGroupAutocomplete(e, b, "server", e.GuildID().String())
			})
			h.Command("/remove", func(e *handler.CommandEvent) error {
				return server.RemoveGroupHandler(e, b)
			})
			h.Autocomplete("/remove", func(e *handler.AutocompleteEvent) error {
				if e.GuildID() == nil {
					return e.AutocompleteResult(nil)
				}
				return manga.HandleWatchlistAutocomplete(e, b, "server", e.GuildID().String(), "title")
			})
		})
		h.Route("/role", func(h handler.Router) {
			h.Command("/set", func(e *handler.CommandEvent) error {
				return server.RoleSetHandler(e, b)
			})
			h.Command("/remove", func(e *handler.CommandEvent) error {
				return server.RoleRemoveHandler(e, b)
			})
		})
		h.Route("/channel", func(h handler.Router) {
			h.Command("/set", func(e *handler.CommandEvent) error {
				return server.ChannelSetHandler(e, b)
			})
		})
	})
	h.Route("/user", func(h handler.Router) {
		h.Command("/setup", func(e *handler.CommandEvent) error {
			return user.SetupHandler(e, b)
		})
		h.Command("/delete", func(e *handler.CommandEvent) error {
			return user.RunDelete(e, b)
		})
		h.Command("/list", func(e *handler.CommandEvent) error {
			return user.ListHandler(e, b)
		})
		h.Command("/add", func(e *handler.CommandEvent) error {
			return user.AddHandler(e, b)
		})
		h.Autocomplete("/add", func(e *handler.AutocompleteEvent) error {
			return manga.HandleAddAutocomplete(e, b, "title")
		})
		h.Command("/remove", func(e *handler.CommandEvent) error {
			return user.RemoveHandler(e, b)
		})
		h.Autocomplete("/remove", func(e *handler.AutocompleteEvent) error {
			return manga.HandleWatchlistAutocomplete(e, b, "user", e.User().ID.String(), "title")
		})
		h.Route("/group", func(h handler.Router) {
			h.Command("/set", func(e *handler.CommandEvent) error {
				return user.SetGroupHandler(e, b)
			})
			h.Autocomplete("/set", func(e *handler.AutocompleteEvent) error {
				return manga.HandleSetGroupAutocomplete(e, b, "user", e.User().ID.String())
			})
			h.Command("/remove", func(e *handler.CommandEvent) error {
				return user.RemoveGroupHandler(e, b)
			})
			h.Autocomplete("/remove", func(e *handler.AutocompleteEvent) error {
				return manga.HandleWatchlistAutocomplete(e, b, "user", e.User().ID.String(), "title")

			})
		})

	})

	// delete
	h.Component("/server_delete_confirm/{action}", func(e *handler.ComponentEvent) error {
		return server.HandleDeleteConfirmation(e, b)
	})
	h.Component("/user_delete_confirm/{action}", func(e *handler.ComponentEvent) error {
		return user.HandleDeleteConfirmation(e, b)
	})

	// list
	h.Component("/list_nav/{mode}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleMangaListPagination(e, b)
	})

	// add
	h.Component("/manga_add_select/{mode}", func(e *handler.ComponentEvent) error {
		return manga.HandleAddSelection(e, b)
	})
	h.Component("/manga_add_confirm/{mode}/{manga_id}/{action}", func(e *handler.ComponentEvent) error {
		return manga.HandleAddConfirmation(e, b)
	})

	// remove
	h.Component("/manga_remove_select/{mode}", func(e *handler.ComponentEvent) error {
		return manga.HandleRemoveSelection(e, b)
	})
	h.Component("/manga_remove_confirm/{mode}/{manga_id}/{action}", func(e *handler.ComponentEvent) error {
		return manga.HandleRemoveConfirmation(e, b)
	})
	h.Component("/remove_nav/{mode}/{query}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleRemovePagination(e, b)
	})

	// search
	h.Component("/manga_search_select", func(e *handler.ComponentEvent) error {
		return manga.HandleSearchSelection(e, b)
	})

	// set group
	h.Component("/setgroup_manga_select/{mode}", func(e *handler.ComponentEvent) error {
		return manga.HandleSetGroupMangaSelection(e, b)
	})
	h.Component("/setgroup_manga_nav/{mode}/{query}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleSetGroupMangaPagination(e, b)
	})
	h.Component("/setgroup_group_select/{mode}/{manga_id}", func(e *handler.ComponentEvent) error {
		return manga.HandleSetGroupGroupSelection(e, b)
	})
	h.Component("/setgroup_group_nav/{mode}/{manga_id}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleSetGroupGroupPagination(e, b)
	})
	h.Component("/setgroup_confirm/{mode}/{manga_id}/{group_id}/{action}", func(e *handler.ComponentEvent) error {
		return manga.HandleSetGroupConfirmation(e, b)
	})

	// set group, but also remove (cuz its in setgroup file)
	h.Component("/group_remove_manga_select/{mode}", func(e *handler.ComponentEvent) error {
		return manga.HandleGroupRemoveMangaSelection(e, b)
	})
	h.Component("/group_remove_manga_nav/{mode}/{query}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleGroupRemoveMangaPagination(e, b)
	})

	// role
	h.Component("/server_role_confirm/{type}/{role_id}/{action}", func(e *handler.ComponentEvent) error {
		return server.HandleRoleConfirmation(e, b)
	})
	h.Component("/server_role_remove_confirm/{type}/{action}", func(e *handler.ComponentEvent) error {
		return server.HandleRoleRemoveConfirmation(e, b)
	})

	// channel
	h.Component("/server_channel_confirm/{channel_id}/{action}", func(e *handler.ComponentEvent) error {
		return server.HandleChannelConfirmation(e, b)
	})

	return h
}
