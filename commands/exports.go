package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common/manga"
	"github.com/jckli/mangaupdates-bot/commands/server"
	"github.com/jckli/mangaupdates-bot/commands/user"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var CommandList = []discord.ApplicationCommandCreate{
	pingCommand,
	server.ServerCommand,
	user.UserCommand,
}

func CommandHandlers(b *mubot.Bot) *handler.Mux {
	h := handler.New()

	h.Command("/ping", PingHandler)
	h.Route("/server", func(h handler.Router) {
		h.Command("/list", func(e *handler.CommandEvent) error {
			return server.ListHandler(e, b)
		})
	})
	h.Route("/user", func(h handler.Router) {
		h.Command("/list", func(e *handler.CommandEvent) error {
			return user.ListHandler(e, b)
		})
	})
	h.Component("/list_nav/{mode}/{page}", func(e *handler.ComponentEvent) error {
		return manga.HandleMangaListPagination(e, b)
	})

	return h
}
