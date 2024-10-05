package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

type EventHandler interface {
	User() discord.User
	Respond(
		responseType discord.InteractionResponseType,
		data discord.InteractionResponseData,
		opts ...rest.RequestOpt,
	) error
	UpdateInteractionResponse(
		messageUpdate discord.MessageUpdate,
		opts ...rest.RequestOpt,
	) (*discord.Message, error)
}

type CommandEventAdapter struct {
	Event *handler.CommandEvent
}

func (c *CommandEventAdapter) User() discord.User {
	return c.Event.User()
}

func (c *CommandEventAdapter) Respond(
	responseType discord.InteractionResponseType,
	data discord.InteractionResponseData,
	opts ...rest.RequestOpt,
) error {
	return c.Event.Respond(responseType, data, opts...)
}

func (c *CommandEventAdapter) UpdateInteractionResponse(
	messageUpdate discord.MessageUpdate,
	opts ...rest.RequestOpt,
) (*discord.Message, error) {
	return c.Event.UpdateInteractionResponse(messageUpdate, opts...)
}

type ComponentEventAdapter struct {
	Event *handler.ComponentEvent
}

func (c *ComponentEventAdapter) User() discord.User {
	return c.Event.User()
}

func (c *ComponentEventAdapter) Respond(
	responseType discord.InteractionResponseType,
	data discord.InteractionResponseData,
	opts ...rest.RequestOpt,
) error {
	return c.Event.Respond(responseType, data, opts...)
}

func (c *ComponentEventAdapter) UpdateInteractionResponse(
	messageUpdate discord.MessageUpdate,
	opts ...rest.RequestOpt,
) (*discord.Message, error) {
	return c.Event.UpdateInteractionResponse(messageUpdate, opts...)
}

type searchResultsFormatted struct {
	Title  string
	Year   string
	Rating float64
}
