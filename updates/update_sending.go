package update_sending

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/jckli/mangaupdates-bot/utils"
)

var (
	oldEntries map[string]utils.MangaEntry
)

func getMangaEntryKey(entry utils.MangaEntry) string {
	return entry.Title + "|" + entry.Chapter + "|" + entry.ScanGroup + "|" + entry.Link
}

func StartRssCheck(b *mubot.Bot) {
	newFullEntries, err := utils.RssParseFeed()
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Failed to parse RSS feed: %s", err.Error()))
		return
	}

	for _, entry := range newFullEntries {
		key := getMangaEntryKey(entry)
		oldEntries = make(map[string]utils.MangaEntry)
		oldEntries[key] = entry
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				checkRssForUpdates(b)
			}
		}
	}()
}

func checkRssForUpdates(b *mubot.Bot) {
	b.Logger.Info("Checking RSS for new updates!")

	newFullEntries, err := utils.RssParseFeed()
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Failed to parse RSS feed: %s", err.Error()))
		return
	}

	newOldEntries := map[string]utils.MangaEntry{}
	newEntries := []utils.MangaEntry{}

	for _, entry := range newFullEntries {
		key := getMangaEntryKey(entry)
		if _, ok := oldEntries[key]; !ok {
			newEntries = append(newEntries, entry)
		}

		newOldEntries[key] = entry
	}

	if len(newEntries) > 0 {
		b.Logger.Info("New entries found.")
		var wg sync.WaitGroup

		for _, entry := range newEntries {
			b.Logger.Info(fmt.Sprintf("Notifying entry: %s", entry.Title))
			wg.Add(1)
			go func(entry utils.MangaEntry) {
				defer wg.Done()
				notify(b, entry)
			}(entry)
		}
		wg.Wait()
		oldEntries = newOldEntries
		return
	}
}

func notify(b *mubot.Bot, entry utils.MangaEntry) {
	errorChannel := snowflake.ID(990005048408936529)

	var (
		image      string
		err        error
		scanGroups []utils.MuSearchGroupsGroup
	)

	if entry.Link != "" {
		idRegex := regexp.MustCompile(`id=([0-9]+)`)
		pathRegex := regexp.MustCompile(`series/([^/]+)`)
		var mangaId int64
		if matches := idRegex.FindStringSubmatch(entry.Link); len(matches) > 1 {
			strMangaId := matches[1]
			mangaId, err = strconv.ParseInt(strMangaId, 10, 64)
			if err != nil {
				b.Logger.Error(fmt.Sprintf("Failed to convert new ID: %s", err.Error()))
				return
			}
		} else if matches := pathRegex.FindStringSubmatch(entry.Link); len(matches) > 1 {
			strMangaId := matches[1]
			mangaId, err = utils.MuConvertNewId(strMangaId)
			if err != nil {
				b.Logger.Error(fmt.Sprintf("Failed to convert new ID: %s", err.Error()))
				return
			}
		} else {
			b.Logger.Error(fmt.Sprintf("Failed to get manga ID from URL: %s", entry.Link))
			return
		}
		entry.NewId = mangaId
		var seriesInfo *utils.MuSeriesInfoResponse
		seriesInfo, err = utils.MuGetSeriesInfo(b, mangaId)
		if err != nil {
			b.Logger.Error(fmt.Sprintf("Failed to get series info: %s", err.Error()))
			return
		}
		image = seriesInfo.Image.URL.Original
	}

	if entry.ScanGroup != "" {
		scanGroups, err = getScanGroups(b, entry.ScanGroup)
		if err != nil {
			b.Logger.Error(fmt.Sprintf("Failed to get scan groups: %s", err.Error()))
			return
		}
	}

	serverWant, userWant, err := getWantLists(b, entry, scanGroups)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Failed to get want lists: %s", err.Error()))
	}
	if serverWant == nil && userWant == nil {
		return
	}

	if serverWant != nil {
		for _, server := range serverWant {
			sendServerUpdate(b, entry, server, image, scanGroups, errorChannel)
		}
	}
	if userWant != nil {
		for _, user := range userWant {
			sendUserUpdate(b, entry, user, image, scanGroups, errorChannel)
		}
	}

	b.Logger.Info(fmt.Sprintf("Finished notifying for %s", entry.Title))
	_, _ = b.Client.Rest().
		CreateMessage(errorChannel, discord.MessageCreate{
			Content: fmt.Sprintf("Finished notifying for %s", entry.Title),
		})

	return
}

func getScanGroups(b *mubot.Bot, scanGroup string) ([]utils.MuSearchGroupsGroup, error) {
	groups := strings.Split(scanGroup, ",")
	var scanGroups []utils.MuSearchGroupsGroup
	for i, group := range groups {
		groups[i] = strings.TrimSpace(group)
		results, err := utils.MuPostSearchGroups(b, group)
		if err != nil {
			return nil, err
		}

		scanGroups = append(scanGroups, results.Results[0])
	}

	return scanGroups, nil
}

func getWantLists(
	b *mubot.Bot,
	entry utils.MangaEntry,
	scanGroups []utils.MuSearchGroupsGroup,
) ([]utils.MDbServer, []utils.MDbUser, error) {
	var serverErr error
	var userErr error
	serverWant, serverErr := utils.DbServersWanted(b, &scanGroups, &entry)
	userWant, userErr := utils.DbUsersWanted(b, &scanGroups, &entry)

	if userErr != nil && serverErr != nil {
		return nil, nil, fmt.Errorf(
			"Both lists errored. userErr: %s, serverErr: %s",
			userErr.Error(),
			serverErr.Error(),
		)
	}
	if userErr != nil {
		return serverWant, nil, userErr
	}
	if serverErr != nil {
		return nil, userWant, serverErr
	}

	return serverWant, userWant, nil
}

func sendServerUpdate(
	b *mubot.Bot,
	entry utils.MangaEntry,
	server utils.MDbServer,
	image string,
	scanGroups []utils.MuSearchGroupsGroup,
	errorChannel snowflake.ID,
) {
	bu, ok := b.Client.Caches().SelfUser()
	embed := discord.NewEmbedBuilder().
		SetTitlef("New %s Chapter!", entry.Title).
		SetDescriptionf("Chapter `%s` has been released!", entry.Chapter).
		SetColor(0x3083e3)
	if ok {
		embed.SetAuthor(bu.Username, "", *bu.AvatarURL())
	}
	if entry.Link != "" {
		embed.SetURL(entry.Link)
	}
	if image != "" {
		embed.SetImage(image)
	}
	if entry.Chapter != "" {
		embed.AddField("Chapter", entry.Chapter, true)
	}
	if scanGroups != nil {
		scanGroupNames := []string{}
		scanGroupLinks := []string{}
		for _, group := range scanGroups {
			scanGroupNames = append(scanGroupNames, group.Record.Name)
			scanGroupLinks = append(scanGroupLinks, group.Record.URL)
		}
		embed.AddField("Scanlator(s)", strings.Join(scanGroupNames, ", "), true)
		embed.AddField("Scanlator Link(s)", strings.Join(scanGroupLinks, ", "), true)
	}

	_, err := b.Client.Rest().
		CreateMessage(snowflake.MustParse(strconv.FormatInt(server.ChannelId, 10)), discord.MessageCreate{
			Embeds: []discord.Embed{embed.Build()},
		})
	if err != nil {
		sendError := fmt.Sprintf("Failed to send message: %s", err.Error())
		b.Logger.Error(sendError)
		_, _ = b.Client.Rest().
			CreateMessage(errorChannel, discord.MessageCreate{
				Content: sendError,
			})
	} else {
		_, _ = b.Client.Rest().CreateMessage(errorChannel, discord.MessageCreate{
			Content: fmt.Sprintf("**SERVER**: Sent message to ID %s\nTitle: %s\nScanlator: %s\nLink: %s", server.ChannelId, entry.Title, entry.ScanGroup, entry.Link),
		})
	}
}

func sendUserUpdate(
	b *mubot.Bot,
	entry utils.MangaEntry,
	user utils.MDbUser,
	image string,
	scanGroups []utils.MuSearchGroupsGroup,
	errorChannel snowflake.ID,
) {
	bu, ok := b.Client.Caches().SelfUser()
	embed := discord.NewEmbedBuilder().
		SetTitlef("New %s Chapter!", entry.Title).
		SetDescriptionf("Chapter `%s` has been released!", entry.Chapter).
		SetColor(0x3083e3)
	if ok {
		embed.SetAuthor(bu.Username, "", *bu.AvatarURL())
	}
	if entry.Link != "" {
		embed.SetURL(entry.Link)
	}
	if image != "" {
		embed.SetImage(image)
	}
	if entry.Chapter != "" {
		embed.AddField("Chapter", entry.Chapter, true)
	}
	if scanGroups != nil {
		scanGroupNames := []string{}
		scanGroupLinks := []string{}
		for _, group := range scanGroups {
			scanGroupNames = append(scanGroupNames, group.Record.Name)
			scanGroupLinks = append(scanGroupLinks, group.Record.URL)
		}
		embed.AddField("Scanlator(s)", strings.Join(scanGroupNames, ", "), true)
		embed.AddField("Scanlator Link(s)", strings.Join(scanGroupLinks, ", "), true)
	}

	userChannel, err := b.Client.Rest().
		CreateDMChannel((snowflake.MustParse(strconv.FormatInt(user.UserId, 10))))
	if err != nil {
		sendError := fmt.Sprintf("Failed to create DM channel: %s", err.Error())
		b.Logger.Error(sendError)
	}

	_, err = b.Client.Rest().
		CreateMessage(userChannel.ID(), discord.MessageCreate{
			Embeds: []discord.Embed{embed.Build()},
		})
	if err != nil {
		sendError := fmt.Sprintf("Failed to send message: %s", err.Error())
		b.Logger.Error(sendError)
		_, _ = b.Client.Rest().
			CreateMessage(errorChannel, discord.MessageCreate{
				Content: sendError,
			})
	} else {
		_, _ = b.Client.Rest().CreateMessage(errorChannel, discord.MessageCreate{
			Content: fmt.Sprintf("**USER**: Sent message to ID %s\nTitle: %s\nScanlator: %s\nLink: %s", user.UserId, entry.Title, entry.ScanGroup, entry.Link),
		})
	}
}
