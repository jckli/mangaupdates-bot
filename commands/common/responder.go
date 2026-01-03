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
	_, err := c.Event.Client().Rest().UpdateInteractionResponse(
		c.Event.ApplicationID(),
		c.Event.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &components,
		},
	)
	return err
}

func (c *CommandResponder) Error(content string) error {
	embed := ErrorEmbed(content)
	_, err := c.Event.Client().Rest().UpdateInteractionResponse(
		c.Event.ApplicationID(),
		c.Event.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &[]discord.ContainerComponent{},
		},
	)
	return err
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
