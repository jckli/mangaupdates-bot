package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/manga"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var CommandList = []discord.ApplicationCommandCreate{
	pingCommand,
	manga.MangaCommand,
	serverCommand,
	userCommand,
}

func CommandHandlers(b *mubot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)

	h.Route("/manga", func(h handler.Router) {

		// manga add
		h.Command("/add", func(e *handler.CommandEvent) error {
			return manga.MangaAddHandler(e, b)
		})
		h.Component("/add/mode/{mode}/{title}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			if mode == "user" {
				adapter := &manga.ComponentEventAdapter{Event: e}
				return manga.MangaAddUserHandler(adapter, b, e.Vars["title"])
			} else {
				return manga.MangaAddServerHandler(e, b, e.Vars["title"])
			}
		})
		h.Component("/add/select/{mode}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			return manga.MangaAddSearchHandler(e, b, mode)
		})
		h.Component(
			"/add/confirm/select/{mode}/{id}/{decision}",
			func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				id := e.Vars["id"]
				decision := e.Vars["decision"]
				if decision == "cancel" {
					return manga.MangaAddCancelHandler(e)
				} else {
					return manga.MangaAddConfirmHandler(e, b, mode, id)
				}
			},
		)

		// manga remove
		h.Command("/remove", func(e *handler.CommandEvent) error {
			return manga.MangaRemoveHandler(e, b)
		})
		h.Component("/remove/mode/{mode}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			if mode == "user" {
				adapter := &manga.ComponentEventAdapter{Event: e}
				return manga.MangaRemoveUserHandler(adapter, b)
			} else {
				return manga.MangaRemoveServerHandler(e, b)
			}
		})
	})

	h.Route("/server", func(h handler.Router) {
		h.Command("/setup", func(e *handler.CommandEvent) error {
			return serverSetupHandler(e, b)
		})
	})

	h.Route("/user", func(h handler.Router) {
		h.Command("/setup", func(e *handler.CommandEvent) error {
			return userSetupHandler(e, b)
		})
	})

	return h
}
