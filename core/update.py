import feedparser
import re
import requests
import os
from bs4 import BeautifulSoup as bs
from bs4 import Comment
import aiohttp
import asyncio

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

s = requests.Session()

def getAllData(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")

    # Get title
    text = soup.find("span", {"class": "releasestitle tabletitle"})
    title = text.get_text()

    # Get type
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[1]
    type = contents.get_text().replace("\n", "")

    # Get year
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[20]
    year = contents.get_text().replace("\n", "")

    # Get finished/unfinished status
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[7]
    statusRaw = contents.get_text().replace("\n", "")
    if statusRaw == "No":
        status = "Ongoing"
    else:
        status = "Completed"

    # Get image
    for img in soup.find_all("img"):
        if img["src"].startswith("https://www.mangaupdates.com/image/"):
            image = img["src"]
            break

    # Get description
    table = soup.find('div', {"class": "col-6 p-2 text"})
    div = table.find('div', {"style": "text-align:justify"})
    if div.find('div', {"style": "display:none"}) != None:
        div = div.find('div', {"style": "display:none"})
    for element in div(text=lambda it: isinstance(it, Comment)):
        element.extract()
    for a in div.find_all("a"):
        a.replace_with("")
    for b in div.find_all("b"):
        b.replace_with("")
    for br in div.find_all("br"):
        br.replace_with("")
    list = [ele for ele in div.contents if ele.strip()]
    newList = []
    for i in (line for line in list if not line.startswith(',')):
        newList.append(i)
    oldDescription = "".join(newList)
    description = oldDescription.replace("\n", "")

    # Get associated names
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
            associatedNames = names
            break
        i += 1

    # Get authors
    table = soup.find_all("div", {"class": "sContent"})
    authors = []
    mainDiv = table[18]
    for i in mainDiv.contents:
        if '\xa0[' in str(i):
            author = i.replace('\xa0[',"")
            authors.append({"name": author, "id": None})
    div = mainDiv.find_all("a")
    for i in div:
        if "<u>Add</u>" not in str(i):
            if "<u>Anthology</u>" not in str(i):
                author = i.get_text()
                link = i.get("href")
                authorid = link.partition("https://www.mangaupdates.com/authors.html?id=")[2]
                authors.append({"name": author, "id": authorid})

    # Get artists
    table = soup.find_all("div", {"class": "sContent"})
    artists = []
    mainDiv = table[19]
    for i in mainDiv.contents:
        if '\xa0[' in str(i):
            artist = i.replace('\xa0[',"")
            artists.append({"name": artist, "id": None})
    div = mainDiv.find_all("a")
    for i in div:
        if "<u>Add</u>" not in str(i):
            if "<u>Anthology</u>" not in str(i):
                artist = i.get_text()
                link = i.get("href")
                artistid = link.partition("https://www.mangaupdates.com/authors.html?id=")[2]
                artists.append({"name": artist, "id": artistid})

    # Get rating
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[11].contents
    average = re.search('Average: (.*?) ', str(contents[0])).group(1)
    bayesianRating = re.search('<b>(.*?)</b>', str(contents[5])).group(1)
    rating = {"average": average, "bayesianRating": bayesianRating}

    # Get latest chapters
    # I just want to say that coding this gave me brain damage
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[5].contents[:-2]
    datesRaw = table[5].find_all("span")
    groupsRaw = []
    k = 0
    v = []
    c = []
    for i in contents:
        if "v." in str(i):
            v.append(k+1)
        if "c." in str(i):
            c.append(k+1)
            if "v." not in contents[k-2]:
                v.append(None)
        if "Group Info" in str(i):
            groupsRaw.append(i)
        if " & " in str(i):
            groupsRaw.append(i)
        k += 1
    volumes = []
    chapters = []
    for index, i in enumerate(c):
        if v[index] is None:
            volumes.append(None)
        else:
            volumes.append(contents[v[index]].get_text())
        chapters.append(contents[i].get_text())
    groups = []
    index = 0
    for i in groupsRaw:
        if "Group Info" in str(i):
            group1 = re.search('<a href="(.*?)" title="Group Info">(.*?)</a>', str(i))
            groupid = group1.group(1).partition("https://www.mangaupdates.com/groups.html?id=")[2]
            group = {"group1": {"name": group1.group(2), "id": groupid}}
            if index == len(groupsRaw)-1:
                groups.append(group)
            else:
                if " & " in str(groupsRaw[index+1]):
                    l = True
                    h = 2
                    while l == True:
                        if index == len(groupsRaw)-1:
                            break
                        else:
                            if " & " not in str(groupsRaw[index+1]):
                                l = False
                            else:
                                newGroup = re.search('<a href="(.*?)" title="Group Info">(.*?)</a>', str(groupsRaw[index+2]))
                                newgroupid = newGroup.group(1).partition("https://www.mangaupdates.com/groups.html?id=")[2]
                                group["group"+str(h)] = {"name": newGroup.group(2), "id": newgroupid}
                                groupsRaw.remove(groupsRaw[index+1])
                                groupsRaw.remove(groupsRaw[index+1])
                                h += 1
                groups.append(group)
        index += 1
    dates = []
    for index, i in enumerate(datesRaw):
        dates.append(datesRaw[index]["title"])
    latestChapters = []
    for index, i in enumerate(chapters):
        latestChapters.append({"volume": volumes[index], "chapter": chapters[index], "groups": groups[index], "date": dates[index]})

    # Return
    return {"title": title, "type": type, "status": status, "year": year, "image": image, "description": description, "associatedNames": associatedNames, "authors": authors, "artists": artists, "rating": rating, "latestChapters": latestChapters}

def getImage(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser") 
    for img in soup.find_all("img"):
        if img["src"].startswith("https://www.mangaupdates.com/image/"):
            image = img["src"]
            break
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
    if "&" in title:
        title = title.replace("&", "&amp;")
    with requests.Session() as s:
        websiteResult = s.post("https://mangaupdates.com/search.html", params={"search": title}, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all('a', {"alt": "Series Info"})
    for manga in table:
        if title in str(manga):
            link = str(manga).replace('<a alt="Series Info" href="', "")
            link = link.partition('">')[0]
            return link

def getGroups(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.findAll('div')
    i = 0
    for div in table:
        if str(div) == '<div class="sCat"><b>Groups Scanlating</b></div>':
            groupsRaw = table[i+1]
            groupsRaw = str(groupsRaw).replace('<div class="sContent">', "")
            groupsRaw = re.sub('<script type="text/javascript">(.*)</script>', "", groupsRaw, flags=re.DOTALL)
            groupsRaw = groupsRaw.replace('<a href="javascript:dispgroups()" id="div_groups_link"><u><b>M</b>ore...</u></a>', "")
            groupsRaw = groupsRaw.replace('<div id="div_groups_more" style="display:none">', "")
            groupsRaw = groupsRaw.replace('</div>', "")
            groupsRaw = groupsRaw.replace('<a href="javascript:dispLessgroups()"><u><b>L</b>ess...</u></a>', "")
            groupsRaw = groupsRaw.replace('<br/>', "\n")
            groupsRaw = os.linesep.join([s for s in groupsRaw.splitlines() if s])
            groupsList = groupsRaw.splitlines()
            groups = []
            for rawGroups in groupsList:
                search = re.search('<a href="(.*?)" title="Group Info"><u>(.*?)</u></a>', rawGroups)
                groupid = search.group(1).partition("https://www.mangaupdates.com/groups.html?id=")[2]
                groups.append({"groupName": search.group(2), "groupid": groupid})
            return groups
        i += 1

def getDescription(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find('div', {"class": "col-6 p-2 text"})
    div = table.find('div', {"style": "text-align:justify"})
    if div.find('div', {"style": "display:none"}) != None:
        div = div.find('div', {"style": "display:none"})
    for element in div(text=lambda it: isinstance(it, Comment)):
        element.extract()
    for a in div.find_all("a"):
        a.replace_with("")
    for b in div.find_all("b"):
        b.replace_with("")
    for br in div.find_all("br"):
        br.replace_with("")
    list = [ele for ele in div.contents if ele.strip()]
    newList = []
    for i in (line for line in list if not line.startswith(',')):
        newList.append(i)
    result = "".join(newList)
    result = result.replace("\n", "")
    return result

def getAuthors(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    authors = []
    mainDiv = table[18]
    for i in mainDiv.contents:
        if '\xa0[' in str(i):
            author = i.replace('\xa0[',"")
            authors.append({"name": author, "id": None})
    div = mainDiv.find_all("a")
    for i in div:
        if "<u>Add</u>" not in str(i):
            if "<u>Anthology</u>" not in str(i):
                author = i.get_text()
                link = i.get("href")
                authorid = link.partition("https://www.mangaupdates.com/authors.html?id=")[2]
                authors.append({"name": author, "id": authorid})
    return authors

def getArtists(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    artists = []
    mainDiv = table[19]
    for i in mainDiv.contents:
        if '\xa0[' in str(i):
            artist = i.replace('\xa0[',"")
            artists.append({"name": artist, "id": None})
    div = mainDiv.find_all("a")
    for i in div:
        if "<u>Add</u>" not in str(i):
            if "<u>Anthology</u>" not in str(i):
                artist = i.get_text()
                link = i.get("href")
                artistid = link.partition("https://www.mangaupdates.com/authors.html?id=")[2]
                artists.append({"name": artist, "id": artistid})
    return artists

def getRating(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[11].contents
    average = re.search('Average: (.*?) ', str(contents[0])).group(1)
    bayesianRating = re.search('<b>(.*?)</b>', str(contents[5])).group(1)
    rating = {"average": average, "bayesianRating": bayesianRating}
    return rating

def getYear(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[20]
    year = contents.get_text().replace("\n", "")
    return year

def getType(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[1]
    type = contents.get_text().replace("\n", "")
    return type

def getStatus(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[7]
    status = contents.get_text().replace("\n", "")
    if status == "No":
        return "Completed"
    else:
        return "Ongoing"

# I just want to say that coding this gave me brain damage
def getLatestChapter(link):
    with requests.Session() as s:
        websiteResult = s.get(link, headers={"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"})
    htmlData = websiteResult.text
    soup = bs(htmlData, "html.parser")
    table = soup.find_all("div", {"class": "sContent"})
    contents = table[5].contents[:-2]
    datesRaw = table[5].find_all("span")
    groupsRaw = []
    k = 0
    v = []
    c = []
    for i in contents:
        if "v." in str(i):
            v.append(k+1)
        if "c." in str(i):
            c.append(k+1)
            if "v." not in contents[k-2]:
                v.append(None)
        if "Group Info" in str(i):
            groupsRaw.append(i)
        if " & " in str(i):
            groupsRaw.append(i)
        k += 1
    volumes = []
    chapters = []
    for index, i in enumerate(c):
        if v[index] is None:
            volumes.append(None)
        else:
            volumes.append(contents[v[index]].get_text())
        chapters.append(contents[i].get_text())

    groups = []
    index = 0
    for i in groupsRaw:
        if "Group Info" in str(i):
            group1 = re.search('<a href="(.*?)" title="Group Info">(.*?)</a>', str(i))
            groupid = group1.group(1).partition("https://www.mangaupdates.com/groups.html?id=")[2]
            group = {"group1": {"name": group1.group(2), "id": groupid}}
            if index == len(groupsRaw)-1:
                groups.append(group)
            else:
                if " & " in str(groupsRaw[index+1]):
                    l = True
                    h = 2
                    while l == True:
                        if index == len(groupsRaw)-1:
                            break
                        else:
                            if " & " not in str(groupsRaw[index+1]):
                                l = False
                            else:
                                newGroup = re.search('<a href="(.*?)" title="Group Info">(.*?)</a>', str(groupsRaw[index+2]))
                                newgroupid = newGroup.group(1).partition("https://www.mangaupdates.com/groups.html?id=")[2]
                                group["group"+str(h)] = {"name": newGroup.group(2), "id": newgroupid}
                                groupsRaw.remove(groupsRaw[index+1])
                                groupsRaw.remove(groupsRaw[index+1])
                                h += 1
                groups.append(group)
        index += 1

    dates = []
    for index, i in enumerate(datesRaw):
        dates.append(datesRaw[index]["title"])

    latestChapters = []
    for index, i in enumerate(chapters):
        latestChapters.append({"volume": volumes[index], "chapter": chapters[index], "groups": groups[index], "date": dates[index]})

    return latestChapters