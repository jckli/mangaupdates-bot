import feedparser
import re
import requests
from bs4 import BeautifulSoup as bs

def getLatest():
    # Parses the rss feed from mangaupdates
    try:
        feed = feedparser.parse("https://www.mangaupdates.com/rss.php")
    except:
        feed = None
        if feed == None:
            try:
                feed = feedparser.parse("https://www.mangaupdates.com/rss.php")
            except:
                print("Error connecting to mangaupdates.com")
    
    mangas = []
    for entry in feed["entries"]:
        title = entry["title"]
        try:
            chapter = re.search(r"(v.\d{1,} )?c.\d{1,}(\.\d)?(-\d{1,}(\.\d)?)?", title).group() 
        except:
            chapter = None
        scanGroup = re.search("(?<=\[).+?(?=\])", title).group()
        title = title[len(scanGroup) + 2:]
        if chapter != None:
            title = title[0: len(title) - len(chapter) - 1]
        else:
            title = title[0:]
        link = entry["links"][0]["href"]
        
        # Get image from mangaupdates
        websiteResult = requests.get(link)
        htmlData = websiteResult.text
        soup = bs(htmlData, 'html.parser') 
        for img in soup.find_all('img'):
            if img['src'].startswith("https://www.mangaupdates.com/image/"):
                image = img['src']
                break

        mangas.append({"title": title, "chapter": chapter, "scanGroup": scanGroup, "link": link, "image": image})
    return mangas