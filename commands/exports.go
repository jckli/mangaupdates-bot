package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var CommandList = []discord.ApplicationCommandCreate{
	pingCommand,
	mangaCommand,
	serverCommand,
}

func CommandHandlers(b *mubot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)

	h.Route("/manga", func(h handler.Router) {
		h.Command("/add", func(e *handler.CommandEvent) error {
			return mangaAddHandler(e, b)
		})
		h.Component("/manga/add/{mode}/{title}", func(e *handler.ComponentEvent) error {
			mode := e.Vars["mode"]
			if mode == "user" {
				adapter := &ComponentEventAdapter{Event: e}
				return userMangaAddHandler(adapter, b, e.Vars["title"])
			} else {
				return serverMangaAddHandler(e, b, e.Vars["title"])
			}
		})
	})

	h.Route("/server", func(h handler.Router) {
		h.Command("/setup", func(e *handler.CommandEvent) error {
			return serverSetupHandler(e, b)
		})
	})

	return h
}
