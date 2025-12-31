package common

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Responder interface {
	Respond(embed discord.Embed, components []discord.ContainerComponent) error
	Error(content string) error
}

type CommandResponder struct {
	Event *handler.CommandEvent
}

func (c *CommandResponder) Respond(embed discord.Embed, components []discord.ContainerComponent) error {
	return c.Event.Respond(discord.InteractionResponseTypeCreateMessage, discord.MessageCreate{
		Embeds: []discord.Embed{embed}, Components: components,
	})
}

func (c *CommandResponder) Error(content string) error {
	return c.Event.Respond(discord.InteractionResponseTypeCreateMessage, discord.MessageCreate{
		Embeds: []discord.Embed{ErrorEmbed(content)},
		Flags:  discord.MessageFlagEphemeral,
	})
}

type ComponentResponder struct {
	Event *handler.ComponentEvent
}

func (c *ComponentResponder) Respond(embed discord.Embed, components []discord.ContainerComponent) error {
	if components == nil {
		components = []discord.ContainerComponent{}
	}

	return c.Event.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed}, Components: &components,
	})
}

func (c *ComponentResponder) Error(content string) error {
	embed := ErrorEmbed(content)
	return c.Event.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed}, Components: &[]discord.ContainerComponent{},
	})
}
