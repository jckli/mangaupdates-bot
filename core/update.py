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

# s = requests.Session()

# Get image from mangaupdates
def getImage(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser") 
    for img in soup.find_all("img"):
        if img["src"].startswith("https://www.mangaupdates.com/image/"):
            image = img["src"]
    return image

def getTitle(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    text = soup.find("span", {"class": "releasestitle tabletitle"})
    title = text.get_text()
    return title

def getAllTitles(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.findAll('div')
    mainTitle = soup.find("span", {"class": "releasestitle tabletitle"})
    title = mainTitle.get_text()
    i = 0
    for div in table:
        if str(div) == '<div class="sCat"><b>Associated Names</b></div>':
            namesRaw = table[i+1]
            namesRaw = str(namesRaw).replace('<div class="sContent">', "")
            namesRaw = namesRaw.replace('</div>', "")
            namesRaw = namesRaw.replace('<br/>', "")
            namesRaw = namesRaw.replace('</br>', "")
            namesRaw = namesRaw.replace('<br>', ",")
            names = namesRaw.split(",")
            names.append(title)
            for name in names:
                if name == "\n":
                    names.remove(name)
            return names
        i += 1

def getLink(title):
    with requests.Session() as s:
        websiteResult = s.post("https://mangaupdates.com/search.html", params={"search": title}, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all('a', {"alt": "Series Info"})
    for manga in table:
        if title in str(manga):
            link = str(manga).replace('<a alt="Series Info" href="', "")
            link = link.replace(f'"><i>{title}</i></a>', "")
            return link
