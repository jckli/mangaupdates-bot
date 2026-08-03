package common

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/mangaupdates-bot/mubot"
)

var (
	ErrServerNotSetup = fmt.Errorf("This server is not set up.\nPlease run `/server setup` first.")
	ErrUserNotSetup   = fmt.Errorf("Your user is not set up.\nPlease run `/user setup` first.")
	ErrNotAdmin       = fmt.Errorf("You do not have permission.\nRequires Server Owner, `Administrator`, `Manage Server`, or the Configured Admin Role.")
	ErrRoleAdmin      = fmt.Errorf("You do not have permission.\nRequires Server Owner, `Administrator`, or `Manage Server`.")
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
	if hasGuildAuthority(b, guildID, member) {
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

	if hasConfiguredAdminRole(config.Roles.Admin, member) {
		return nil
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

	if hasGuildAuthority(b, guildID, member) {
		return nil
	}
	if hasConfiguredAdminRole(config.Roles.Admin, member) {
		return nil
	}
	return ErrNotAdmin
}

func GuardServerRole(b *mubot.Bot, guildID string, member *discord.ResolvedMember, roleType string) error {
	config, err := b.ApiClient.GetServerConfig(guildID)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}
	if config == nil {
		return ErrServerNotSetup
	}
	if hasGuildAuthority(b, guildID, member) || (roleType == "ping" && hasConfiguredAdminRole(config.Roles.Admin, member)) {
		return nil
	}
	if roleType == "admin" {
		return ErrRoleAdmin
	}
	return ErrNotAdmin
}

func hasConfiguredAdminRole(adminRoleID int64, member *discord.ResolvedMember) bool {
	for _, roleID := range member.RoleIDs {
		if adminRoleID > 0 && roleID == snowflake.ID(adminRoleID) {
			return true
		}
	}
	return false
}

func hasGuildAuthority(b *mubot.Bot, guildID string, member *discord.ResolvedMember) bool {
	if member.Permissions.Has(discord.PermissionAdministrator) || member.Permissions.Has(discord.PermissionManageGuild) {
		return true
	}
	id, err := snowflake.Parse(guildID)
	if err != nil {
		return false
	}
	guild, ok := b.Client.Caches.Guild(id)
	return ok && guild.OwnerID == member.User.ID
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
