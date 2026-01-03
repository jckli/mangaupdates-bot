package common

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var (
	ErrServerNotSetup = fmt.Errorf("This server is not set up.\nPlease run `/server setup` first.")
	ErrUserNotSetup   = fmt.Errorf("Your user is not set up.\nPlease run `/user setup` first.")
	ErrNotAdmin       = fmt.Errorf("You do not have permission.\nRequires `Manage Server` or the Configured Admin Role.")
)

func GuardWidget(e *handler.ComponentEvent, b *mubot.Bot, requireAdmin bool) error {
	mode := e.Vars["mode"]

	if mode == "server" {
		if requireAdmin {
			return GuardServerAdmin(b, e.GuildID().String(), e.Member())
		} else {
			return GuardServerExists(b, e.GuildID().String())
		}
	} else {
		return GuardUser(b, e.User().ID.String())
	}
}

func GuardServerExists(b *mubot.Bot, guildID string) error {
	config, err := b.ApiClient.GetServerConfig(guildID)
	if err != nil {
		return fmt.Errorf("failed to check server status: %w", err)
	}
	if config == nil {
		return ErrServerNotSetup
	}
	return nil
}

func GuardAdminOnly(b *mubot.Bot, guildID string, member *discord.ResolvedMember) error {
	if member.Permissions.Has(discord.PermissionManageGuild) {
		return nil
	}

	config, err := b.ApiClient.GetServerConfig(guildID)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	// if server is not set up (nil config) and they failed the check above, they simply dont have permission
	if config == nil {
		return ErrNotAdmin
	}

	if config.Roles.Admin == 0 {
		return ErrNotAdmin
	}

	adminRoleID := fmt.Sprintf("%d", config.Roles.Admin)
	for _, roleID := range member.RoleIDs {
		if roleID.String() == adminRoleID {
			return nil
		}
	}

	return ErrNotAdmin
}

func GuardServerAdmin(b *mubot.Bot, guildID string, member *discord.ResolvedMember) error {
	config, err := b.ApiClient.GetServerConfig(guildID)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	if config == nil {
		return ErrServerNotSetup
	}

	if member.Permissions.Has(discord.PermissionManageGuild) {
		return nil
	}

	if config.Roles.Admin == 0 {
		return ErrNotAdmin
	}

	adminRoleID := fmt.Sprintf("%d", config.Roles.Admin)
	for _, roleID := range member.RoleIDs {
		if roleID.String() == adminRoleID {
			return nil
		}
	}

	return ErrNotAdmin
}

func GuardUser(b *mubot.Bot, userID string) error {
	config, err := b.ApiClient.GetUserConfig(userID)
	if err != nil {
		return fmt.Errorf("failed to check user profile: %w", err)
	}
	if config == nil {
		return ErrUserNotSetup
	}
	return nil
}
