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
	pingCommand,
	server.ServerCommand,
	user.UserCommand,
	search.SearchCommand,
}

func CommandHandlers(b *mubot.Bot) *handler.Mux {
	h := handler.New()

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
			return manga.HandleRemoveAutocomplete(e, b, "server", e.GuildID().String(), "title")
		})
	})
	h.Route("/user", func(h handler.Router) {
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
			return manga.HandleRemoveAutocomplete(e, b, "user", e.User().ID.String(), "title")
		})
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

	return h
}
