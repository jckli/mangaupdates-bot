package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var (
	startTime = time.Now()

	infoCommand = discord.SlashCommandCreate{
		Name:        "mangaupdates",
		Description: "Displays basic information about the bot",
	}
)

func InfoHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	responder := &common.CommandResponder{Event: e}

	var (
		guildCount  int
		memberCount int
	)

	e.Client().Caches().GuildsForEach(func(guild discord.Guild) {
		guildCount++
		memberCount += guild.MemberCount
	})

	uptime := time.Since(startTime)
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	var uptimeStr string
	if days > 0 {
		uptimeStr = fmt.Sprintf("%d days, ", days)
	}
	uptimeStr += fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	botUser, _ := e.Client().Caches().SelfUser()
	botIcon := ""
	if botUser.AvatarURL() != nil {
		botIcon = *botUser.AvatarURL()
	}

	description := fmt.Sprintf(
		"Thanks for using MangaUpdates Bot! MangaUpdates is a simple but powerful bot that sends every new manga, manhwa, or doujin chapter update to either your direct messages or a server channel.\n\n"+
			"Any questions can be brought up in the support server. This bot is open-source! All code can be found on GitHub (Please leave a star ⭐ if you enjoy the bot).\n\n"+
			"**Server Count:** %d\n"+
			"**User Count:** %d\n"+
			"**Bot Uptime**: %s",
		guildCount,
		memberCount,
		uptimeStr,
	)

	embed := discord.NewEmbedBuilder().
		SetTitle("MangaUpdates Bot").
		SetAuthor("MangaUpdates", "", botIcon).
		SetColor(common.ColorPrimary).
		SetDescription(description).
		SetTimestamp(time.Now()).
		Build()

	actionRow := discord.NewActionRow(
		discord.NewLinkButton("Support Server", "https://jackli.dev/discord"),
		discord.NewLinkButton("GitHub", "https://github.com/jckli/mangaupdates-bot"),
		discord.NewLinkButton("Invite Bot", "https://jackli.dev/mangaupdates"),
	)

	return responder.Respond(embed, []discord.ContainerComponent{actionRow})
}
