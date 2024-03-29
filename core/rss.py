import feedparser
import re

class RSSParser:
    async def __get_latest(self):
        try:
            url = "https://api.mangaupdates.com/v1/releases/rss"
            feed = feedparser.parse(url)
        except:
            try:
                url = "https://api.mangaupdates.com/v1/releases/rss"
                feed = feedparser.parse(url)
            except:
                print("Error: Could not get MangaUpdates RSS feed.")
                return None
        return feed

    async def parse_feed(self):
        feed = await RSSParser.__get_latest(self)
        if feed is None:
            return None
        manga_list = []
        for entry in feed.entries:
            title = entry.title
            try:
                chapter = re.search(r"(v.\d{1,} )?c.\d{1,}(\.\d)?(-\d{1,}(\.\d)?)?", title).group()
                title = title[0: len(title) - len(chapter) - 1]
            except:
                chapter = None
            try:
                scan_group = re.search("(?<=\[).+?(?=\])", title).group()
                title = re.sub("(?<=\[).+?(?=\])", "", title)[3:]
            except:
                scan_group = None
            try:
                if (re.match("^http://", entry.link)):
                    link = re.sub("^http://", "https://", entry.link)
                else:
                    link = entry.link
            except:
                link = None
            
            manga_list.append({"title": title, "chapter": chapter, "scan_group": scan_group, "link": link})
        return manga_list