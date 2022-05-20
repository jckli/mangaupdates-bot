import feedparser
import re

class RSSParser:
    def __get_latest():
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

    def parse_feed():
        feed = RSSParser.__get_latest()
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
                link = entry.link
            except:
                link = None
            
            manga_list.append({"title": title, "chapter": chapter, "scan_group": scan_group, "link": link})
        return manga_list