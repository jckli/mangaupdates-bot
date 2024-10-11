package manga

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func MangaRemoveHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	_, inGuild := e.Guild()
	if !inGuild {
		adapter := &CommandEventAdapter{Event: e}
		return MangaRemoveUserHandler(adapter, b)
	} else {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.MessageCreate{
				Embeds: []discord.Embed{
					selectServerOrUserEmbed(
						"Remove Manga",
						"Would you like to remove a manga from this server or your DMs?",
					),
				},
				Components: selectServerOrUserComponents("manga", "remove", ""),
			},
		)
	}
}

func MangaRemoveUserHandler(
	e EventHandler,
	b *mubot.Bot,
) error {
	userId := int64(e.User().ID)

	exists, err := utils.DbUserCheckExists(b, userId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to check if user exists (MangaRemoveUserHandler): %s", err.Error()),
		)
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(utils.DcErrorTechnicalErrorEmbed()).Build(),
		)
	}
	if !exists {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().SetEmbeds(errorMangaSetupNeededEmbed()).Build(),
		)
	}

	return nil
}

func MangaRemoveServerHandler(
	e *handler.ComponentEvent,
	b *mubot.Bot,
) error {
	serverId := int64(*e.GuildID())

	exists, err := utils.DbServerCheckExists(b, serverId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf(
				"Failed to check if server exists (MangaRemoveServerHandler): %s",
				err.Error(),
			),
		)
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{utils.DcErrorTechnicalErrorEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}
	if !exists {
		return e.UpdateMessage(
			discord.MessageUpdate{
				Embeds:     &[]discord.Embed{errorMangaSetupNeededEmbed()},
				Components: &[]discord.ContainerComponent{},
			},
		)
	}

	return nil
}
