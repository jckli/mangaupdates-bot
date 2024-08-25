package update_sending

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

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
	errorChannnel := snowflake.ID(990005048408936529)

	var (
		image         string
		mangaId       string
		err           error
		scanGroupLink string
		scanGroups    []utils.MuSearchGroupsGroup
	)

	if entry.Link != "" {
		urlMangaIdRegex := regexp.MustCompile(`(?<=series/).+?(?=/)`)
		mangaId, err := utils.MuConvertNewId(urlMangaIdRegex.FindString(entry.Link))
		if err != nil {
			b.Logger.Error(fmt.Sprintf("Failed to convert new ID: %s", err.Error()))
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

}

func getScanGroups(b *mubot.Bot, scanGroup string) ([]utils.MuSearchGroupsGroup, error) {
	groups := strings.Split(scanGroup, "&")
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
