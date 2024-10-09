package commands

import (
	"fmt"
	"strconv"
	"strings"

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
				"/"+command+"/"+subcommand+"/mode/server/"+title,
			),
			discord.NewSecondaryButton(
				"User (DMs)",
				"/"+command+"/"+subcommand+"/mode/user/"+title,
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
		if i >= 25 {
			break
		}
		description += fmt.Sprintf(
			"%d. %s (%s, Rating: %.2f)\n",
			i+1,
			result.Record.Title,
			result.Record.Year,
			result.Record.BayesianRating,
		)

		allResults = append(allResults, searchResultsFormatted{
			Title:  result.Record.Title,
			Year:   result.Record.Year,
			Rating: result.Record.BayesianRating,
			Id:     result.Record.SeriesID,
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
	command, subcommand, mode string,
	results []searchResultsFormatted,
) []discord.ContainerComponent {
	options := []discord.StringSelectMenuOption{}
	for i, result := range results {
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%d. %s", i+1, result.Title),
			Description: fmt.Sprintf("%s, Rating: %.2f", result.Year, result.Rating),
			Value:       strconv.Itoa(result.Id),
		})

	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:    "/" + command + "/" + subcommand + "/select/" + mode,
				Placeholder: "Select a Manga",
				Options:     options,
			},
		},
	}
}

func confirmMangaEmbed(b *mubot.Bot, embedTitle string, mangaId int64) discord.Embed {
	seriesInfo, err := utils.MuGetSeriesInfo(b, mangaId)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to get series info (searchMangaAddHandler): %s", err.Error()),
		)
		return discord.NewEmbedBuilder().
			SetTitle("Error").
			SetDescription("Could not get series info, please try again later.").
			SetColor(0xff4f4f).
			Build()
	}

	description, err := utils.MuCleanupDescription(seriesInfo.Description)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("Failed to cleanup description (searchMangaAddHandler): %s", err.Error()),
		)
		return errorTechnicalErrorEmbed()
	}

	authorArray := []string{}
	for _, author := range seriesInfo.Authors {
		authorArray = append(authorArray, author.Name)
	}
	authorString := strings.Join(authorArray, ", ")

	embed := discord.NewEmbedBuilder().
		SetColor(0x3083e3).
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("**Did you mean to add `%s`?**", seriesInfo.Title)).
		AddField("Description", description, false).
		AddField("Author(s)", authorString, true).
		AddField("Year", seriesInfo.Year, true).
		SetImage(seriesInfo.Image.URL.Original).
		Build()

	return embed
}

func selectConfirmMangaComponents(
	command, subcommand, mode, mangaId string,
) []discord.ContainerComponent {

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewDangerButton(
				"Cancel",
				"/"+command+"/"+subcommand+"/confirm/select/"+mode+"/"+mangaId+"/cancel",
			),
			discord.NewSuccessButton(
				"Confirm",
				"/"+command+"/"+subcommand+"/confirm/select/"+mode+"/"+mangaId+"/confirm",
			),
		},
	}
}

func cancelMangaEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("Successfully cancelled.").
		SetColor(0xff4f4f).
		Build()
	return embed
}

func successMangaAddEmbed(embedTitle, mangaTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription(fmt.Sprintf("Successfully added `%s` to your manga list.", mangaTitle)).
		SetColor(0x3083e3).
		Build()
	return embed
}

func mangaExistsEmbed(embedTitle string) discord.Embed {
	embed := discord.NewEmbedBuilder().
		SetTitle(embedTitle).
		SetDescription("This manga is already in your list.").
		SetColor(0xff4f4f).
		Build()
	return embed
}
