package server

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
)

func RoleSetHandler(e *handler.CommandEvent, b *mubot.Bot) error {
	responder := &common.CommandResponder{Event: e}
	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return responder.Error(err.Error())
	}

	data := e.SlashCommandInteractionData()
	role := data.Role("role")
	roleType := data.String("type")

	desc := fmt.Sprintf("Are you sure you want to set %s as the **%s** role?\n\nUsers with this role will be able to manage the bot.", role.Mention(), roleType)
	if roleType != "admin" {
		desc = fmt.Sprintf("Are you sure you want to set %s as the **%s** role?", role.Mention(), roleType)
	}

	embed := common.StandardEmbed("Confirm Role Update", desc)

	confirmPath := fmt.Sprintf("/server_role_confirm/%s/%s/yes", roleType, role.ID)
	cancelPath := fmt.Sprintf("/server_role_confirm/%s/%s/no", roleType, role.ID)

	buttons := common.CreateConfirmButtons(confirmPath, cancelPath)

	return responder.Respond(embed, buttons)
}

func HandleRoleConfirmation(e *handler.ComponentEvent, b *mubot.Bot) error {
	e.DeferUpdateMessage()

	if err := common.GuardServerAdmin(b, e.GuildID().String(), e.Member()); err != nil {
		return err
	}

	action := e.Vars["action"]
	roleType := e.Vars["type"]
	roleID := e.Vars["role_id"]

	if action == "no" {
		_, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					common.StandardEmbed("Cancelled", "Role configuration was not changed."),
				},
				Components: &[]discord.ContainerComponent{},
			})
		return err
	}

	err := b.ApiClient.SetServerRole(e.GuildID().String(), roleID, roleType)
	if err != nil {
		errEmbed := common.StandardEmbed("Error", err.Error())
		errEmbed.Color = common.ColorError
		_, _ = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.MessageUpdate{Embeds: &[]discord.Embed{errEmbed}})
		return err
	}

	displayType := roleType
	if len(roleType) > 0 {
		displayType = strings.ToUpper(roleType[:1]) + roleType[1:]
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				common.StandardEmbed("Success", fmt.Sprintf("**%s** role has been updated.", displayType)),
			},
			Components: &[]discord.ContainerComponent{},
		})
	return err
}
