package commands

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var (
	alertCommand = discord.SlashCommandCreate{
		Name:        "alert",
		Description: "View latest announcements about the bot",
	}
)

func AlertHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	responder := &common.CommandResponder{Event: e}

	botUser, _ := e.Client().Caches().SelfUser()
	botIcon := ""
	if botUser.AvatarURL() != nil {
		botIcon = *botUser.AvatarURL()
	}

	embed := discord.NewEmbedBuilder().
		SetTitle("📢 Bot Announcements").
		SetAuthor("MangaUpdates", "", botIcon).
		SetColor(0xffd700).
		SetDescription(
			"If you have any suggestions for improving the bot, please tell me in the support server.",
		).
		SetTimestamp(time.Now()).
		Build()

	actionRow := discord.NewActionRow(
		discord.NewLinkButton("Report Bugs", "https://jackli.dev/discord"),
	)

	return responder.Respond(embed, []discord.ContainerComponent{actionRow})
}
