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

		// manga set
		h.Route("/scanlator", func(h handler.Router) {

			// manga scanlator add
			h.Command("/add", func(e *handler.CommandEvent) error {
				return manga.MangaScanlatorAddHandler(e, b)
			})
			h.Component("/add/mode/{mode}", func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				if mode == "user" {
					adapter := &manga.ComponentEventAdapter{Event: e}
					return manga.MangaScanlatorAddUserHandler(adapter, b)
				} else {
					return manga.MangaScanlatorAddServerHandler(e, b)
				}
			})
			h.Component(
				"/add/search/mode/{mode}/{page}",
				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					page := e.Vars["page"]
					pageInt, _ := strconv.Atoi(page)
					if mode == "user" {
						adapter := &manga.ComponentEventAdapter{Event: e}
						return manga.MangaScanlatorAddUserSearchHandler(adapter, b, pageInt)
					} else {
						return manga.MangaScanlatorAddServerSearchHandler(e, b, pageInt)
					}
				},
			)
			h.Component("/add/manga/select/{mode}", func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				return manga.MangaScanlatorAddMangaSelectHandler(e, b, mode)
			})
			h.Component(
				"/add/groups/p/{mangaId}/mode/{mode}/{page}",
				func(e *handler.ComponentEvent) error {
					id := e.Vars["mangaId"]
					mode := e.Vars["mode"]
					page := e.Vars["page"]
					pageInt, _ := strconv.Atoi(page)
					return manga.MangaScanlatorAddGroupHandler(e, b, pageInt, mode, id)
				},
			)
			h.Component(
				"/add/groups/select/{mode}/{mangaId}",
				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					id := e.Vars["mangaId"]
					return manga.MangaScanlatorAddGroupSelectHandler(e, b, mode, id)
				},
			)
			h.Component(
				"/add/confirm/select/{mode}/{mangaId}/{groupId}/{decision}",

				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					mangaId := e.Vars["mangaId"]
					groupId := e.Vars["groupId"]
					decision := e.Vars["decision"]
					if decision == "cancel" {
						return manga.MangaScanlatorAddGroupCancelHandler(e)
					} else {
						return manga.MangaScanlatorAddGroupConfirmHandler(e, b, mode, mangaId, groupId)
					}
				},
			)

			// manga scanlator remove
			h.Command("/remove", func(e *handler.CommandEvent) error {
				return manga.MangaScanlatorRemoveHandler(e, b)
			})
			h.Component("/remove/mode/{mode}", func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				if mode == "user" {
					adapter := &manga.ComponentEventAdapter{Event: e}
					return manga.MangaScanlatorRemoveUserHandler(adapter, b)
				} else {
					return manga.MangaScanlatorRemoveServerHandler(e, b)
				}
			})
			h.Component(
				"/remove/search/mode/{mode}/{page}",
				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					page := e.Vars["page"]
					pageInt, _ := strconv.Atoi(page)
					if mode == "user" {
						adapter := &manga.ComponentEventAdapter{Event: e}
						return manga.MangaScanlatorRemoveUserSearchHandler(adapter, b, pageInt)
					} else {
						return manga.MangaScanlatorRemoveServerSearchHandler(e, b, pageInt)
					}
				},
			)
			h.Component("/remove/manga/select/{mode}", func(e *handler.ComponentEvent) error {
				mode := e.Vars["mode"]
				mangaId := e.StringSelectMenuInteractionData().Values[0]
				if mode == "user" {
					return manga.MangaScanlatorRemoveUserGroupHandler(e, b, 1, mangaId)
				} else {
					return manga.MangaScanlatorRemoveServerGroupHandler(e, b, 1, mangaId)
				}
			})
			h.Component(
				"/remove/groups/p/{mangaId}/mode/{mode}/{page}",
				func(e *handler.ComponentEvent) error {
					id := e.Vars["mangaId"]
					mode := e.Vars["mode"]
					page := e.Vars["page"]
					pageInt, _ := strconv.Atoi(page)
					if mode == "user" {
						return manga.MangaScanlatorRemoveUserGroupHandler(e, b, pageInt, id)
					} else {
						return manga.MangaScanlatorRemoveServerGroupHandler(e, b, pageInt, id)
					}
				},
			)
			h.Component(
				"/remove/groups/select/{mode}/{mangaId}",
				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					id := e.Vars["mangaId"]
					return manga.MangaScanlatorRemoveGroupSelectHandler(e, b, mode, id)
				},
			)
			h.Component(
				"/remove/confirm/select/{mode}/{mangaId}/{groupId}/{decision}",

				func(e *handler.ComponentEvent) error {
					mode := e.Vars["mode"]
					mangaId := e.Vars["mangaId"]
					groupId := e.Vars["groupId"]
					decision := e.Vars["decision"]
					if decision == "cancel" {
						return manga.MangaScanlatorRemoveGroupCancelHandler(e)
					} else {
						if mode == "user" {
							return manga.MangaScanlatorRemoveUserGroupConfirmHandler(e, b, mangaId, groupId)
						} else {
							return manga.MangaScanlatorRemoveServerGroupConfirmHandler(e, b, mangaId, groupId)
						}
					}
				},
			)

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
