package update_sending

import (
	"fmt"
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

	if entry.Link != "" {

	}

}
