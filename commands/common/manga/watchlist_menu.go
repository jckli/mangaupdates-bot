package manga

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/jckli/mangaupdates-bot/commands/common"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

type WatchlistMenuConfig struct {
	Endpoint            string
	TargetID            string
	Query               string
	Page                int
	SelectIDPrefix      string
	NavIDPrefix         string
	Title               string
	DropdownPlaceholder string
}

func GenerateWatchlistMenu(
	b *mubot.Bot,
	cfg WatchlistMenuConfig,
) (discord.Embed, []discord.ContainerComponent, error) {
	list, err := b.ApiClient.GetWatchlist(cfg.Endpoint, cfg.TargetID)
	if err != nil {
		return discord.Embed{}, nil, fmt.Errorf("failed to fetch watchlist")
	}
	if list == nil || len(*list) == 0 {
		return discord.Embed{}, nil, fmt.Errorf("your watchlist is empty")
	}

	var matches []utils.MangaSearchResult
	var displayLines []string
	queryLower := strings.ToLower(cfg.Query)

	for _, item := range *list {
		if cfg.Query == "" || strings.Contains(strings.ToLower(item.Title), queryLower) {
			matches = append(matches, utils.MangaSearchResult{ID: item.ID, Title: item.Title})

			line := fmt.Sprintf("• %s", item.Title)
			if item.GroupName != "" {
				line += fmt.Sprintf(" (*%s*)", item.GroupName)
			}
			displayLines = append(displayLines, line)
		}
	}

	if len(matches) == 0 {
		return discord.Embed{}, nil, fmt.Errorf("no manga found matching `%s`", cfg.Query)
	}

	slicedMatches, totalPages := common.GetPageSlice(matches, cfg.Page)
	slicedLines, _ := common.GetPageSlice(displayLines, cfg.Page)

	customSelectID := fmt.Sprintf("%s/%s", cfg.SelectIDPrefix, cfg.Endpoint)
	dropdown := common.GenerateSearchDropdown(customSelectID, cfg.DropdownPlaceholder, slicedMatches)

	safeQuery := cfg.Query
	if safeQuery == "" {
		safeQuery = "-"
	}
	btnPrefix := fmt.Sprintf("%s/%s/%s", cfg.NavIDPrefix, cfg.Endpoint, safeQuery)
	buttons := common.GeneratePaginationButtons(btnPrefix, cfg.Page, totalPages)

	allComponents := append(dropdown, buttons...)

	var description string
	for _, line := range slicedLines {
		description += line + "\n"
	}

	header := fmt.Sprintf("Found %d results.", len(matches))
	if cfg.Query != "" {
		header = fmt.Sprintf("Found %d results for `%s`.", len(matches), cfg.Query)
	}
	fullDescription := header + "\nPlease select one from the dropdown below:\n\n" + description

	embed := common.StandardEmbed(cfg.Title, fullDescription)
	embed.Footer = &discord.EmbedFooter{Text: fmt.Sprintf("Page %d/%d", cfg.Page, totalPages)}

	return embed, allComponents, nil
}

func HandleGenericWatchlistPagination(
	e *handler.ComponentEvent,
	b *mubot.Bot,
	cfg WatchlistMenuConfig,
) error {
	e.DeferUpdateMessage()

	page, _ := strconv.Atoi(e.Vars["page"])
	cfg.Page = page

	embed, components, err := GenerateWatchlistMenu(b, cfg)
	if err != nil {
		return err
	}

	_, err = e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.MessageUpdate{
			Embeds:     &[]discord.Embed{embed},
			Components: &components,
		})
	return err
}
