package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

func errorMangaSetupNeededEmbed() discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle("Error").
		SetDescription("Please run the setup command first.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func selectServerOrUserComponents(command, subcommand, title string) []discord.ContainerComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewSecondaryButton(
				"Server",
				"/"+command+"/"+subcommand+"/server/"+title,
			),
			discord.NewSecondaryButton(
				"User (DMs)",
				"/"+command+"/"+subcommand+"/user"+title,
			),
		},
	}
}

func selectServerOrUserEmbed(embedTitle, embedDescription string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(embedDescription).
		SetColor(0x3083e3).
		Build()

	return embed

}

func searchResultsEmbed(
	b *mubot.Bot,
	embedTitle, mangaTitle string,
) (discord.Embed, []searchResultsFormatted) {
	searchResults, err := utils.MuPostSearchSeries(b, mangaTitle)
	if err != nil {
		embed := discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("Failed to search for series. Try again later").
			SetColor(0xff4f4f).
			Build()
		return embed, nil
	}

	description := "Select a manga from the search results:\n"
	if len(searchResults.Results) == 0 {
		description = "No results found for: `" + mangaTitle + "`. Try again or input a full https://mangaupdates.com link."
		return discord.NewEmbedBuilder().
			SetTitle(embedTitle).
			SetDescription(description).
			SetColor(0x3083e3).
			Build(), nil
	}

	allResults := []searchResultsFormatted{}
	for i, result := range searchResults.Results {
		description += fmt.Sprintf(
			"%d. %s (%d, Rating: %.2f)\n",
			i+1,
			result.Record.Title,
			result.Record.Year,
			result.Record.BayesianRating,
		)

		allResults = append(allResults, searchResultsFormatted{
			Title:  result.Record.Title,
			Year:   result.Record.Year,
			Rating: result.Record.BayesianRating,
		})
	}

	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(description).
		SetColor(0x3083e3).
		Build()
	return embed, allResults
}

func dropdownSearchResultsComponents(
	command, subcommand string,
	results []searchResultsFormatted,
) []discord.SelectMenuComponent {
	return []discord.ContainerComponent{
		discord.ActionRowComponent{

		},
	}
	options := []discord.SelectMenuOption{}
	for i, result := range results {
			Label: fmt.Sprintf("%d. %s", i, result.Title),
			Value: fmt.Sprintf("%d, Rating: %.2f", result.Year, result.Rating),
		})
	}

	return []discord.SelectMenuComponent{
		discord.NewSelectMenu().
			SetCustomID(command + "/" + subcommand).
			SetPlaceholder("Select a manga").
			SetOptions(options...),
	}
}
