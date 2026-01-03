package commands

import (
	"fmt"
	"strings"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var (
	helpCommand = discord.SlashCommandCreate{
		Name:        "help",
		Description: "Show a list of all available commands",
	}

	cachedHelpEmbed discord.Embed
	helpOnce        sync.Once
)

func HelpHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	responder := &common.CommandResponder{Event: e}

	helpOnce.Do(func() {
		cachedHelpEmbed = buildHelpEmbed(b)
	})

	return responder.Respond(cachedHelpEmbed, nil)
}

func buildHelpEmbed(b *mubot.Bot) discord.Embed {
	var sb strings.Builder

	for _, cmd := range CommandList {
		if slashCmd, ok := cmd.(discord.SlashCommandCreate); ok {
			sb.WriteString(fmt.Sprintf("**/%s** - %s\n", slashCmd.Name, slashCmd.Description))
			for _, opt := range slashCmd.Options {
				if sub, ok := opt.(discord.ApplicationCommandOptionSubCommand); ok {
					sb.WriteString(fmt.Sprintf("> `/%s %s` - %s\n", slashCmd.Name, sub.Name, sub.Description))
				}
				if group, ok := opt.(discord.ApplicationCommandOptionSubCommandGroup); ok {
					for _, sub := range group.Options {
						sb.WriteString(fmt.Sprintf("> `/%s %s %s` - %s\n", slashCmd.Name, group.Name, sub.Name, sub.Description))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	botIcon := ""
	if self, ok := b.Client.Caches().SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	return discord.NewEmbedBuilder().
		SetAuthor("MangaUpdates", "", botIcon).
		SetDescription(sb.String()).
		SetColor(common.ColorPrimary).
		SetFooterTextf("Version: %s", b.Version).
		Build()
}
