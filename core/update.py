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
        try:
            link = entry["links"][0]["href"]
        except:
            link = None
        mangas.append({"title": title.rstrip(), "chapter": chapter, "scanGroup": scanGroup, "link": link})
    return mangas

# Get image from mangaupdates
def getImage(link):
    websiteResult = requests.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser") 
    for img in soup.find_all("img"):
        if img["src"].startswith("https://www.mangaupdates.com/image/"):
            image = img["src"]
    return image

def getTitle(link):
    websiteResult = requests.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    text = soup.find("span", {"class": "releasestitle tabletitle"})
    title = text.get_text()
    return title