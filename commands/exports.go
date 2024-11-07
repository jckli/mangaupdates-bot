package commands

import (
	"strconv"

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
			return manga.MangaAddSelectHandler(e, b, mode)
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
		h.Component("/remove/search/mode/{mode}/{page}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			page := e.Vars["page"]
			pageInt, _ := strconv.Atoi(page)
			if mode == "user" {
				adapter := &manga.ComponentEventAdapter{Event: e}
				return manga.MangaRemoveUserSearchHandler(adapter, b, pageInt)
			} else {
				return manga.MangaRemoveServerSearchHandler(e, b, pageInt)
			}
		})
		h.Component(
			"/remove/select/{mode}",
			func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				return manga.MangaRemoveSelectHandler(e, b, mode)
			},
		)
		h.Component(
			"/remove/confirm/select/{mode}/{id}/{decision}",
			func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				id := e.Vars["id"]
				decision := e.Vars["decision"]
				if decision == "cancel" {
					return manga.MangaRemoveCancelHandler(e)
				} else {
					return manga.MangaRemoveConfirmHandler(e, b, mode, id)
				}
			},
		)

		// manga list
		h.Command("/list", func(e *handler.CommandEvent) error {
			return manga.MangaListHandler(e, b)
		})
		h.Component("/list/mode/{mode}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			if mode == "user" {
				adapter := &manga.ComponentEventAdapter{Event: e}
				return manga.MangaListUserHandler(adapter, b)
			} else {
				return manga.MangaListServerHandler(e, b)
			}
		})
		h.Component("/list/p/mode/{mode}/{page}", func(e *handler.ComponentEvent) error {

			mode := e.Vars["mode"]
			page := e.Vars["page"]
			pageInt, _ := strconv.Atoi(page)
			if mode == "user" {
				adapter := &manga.ComponentEventAdapter{Event: e}
				return manga.MangaListUserListHandler(adapter, b, pageInt)
			} else {
				return manga.MangaListServerListHandler(e, b, pageInt)
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
