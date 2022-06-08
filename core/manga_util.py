import discord
from core.utils import Util

util = Util()

class SearchData:
    def __init__(self, series_info):
        self.title = series_info["title"]
        desc = series_info["description"]
        self.description = util.format_mu_description(desc)
        if series_info["completed"] is True:
            self.status = "Completed"
        else:
            self.status = "Ongoing"
        self.image = series_info["image"]["url"]["original"]
        self.url = series_info["url"]
        self.mangatype = series_info["type"]
        self.year = series_info["year"]
        self.latest_chapter = series_info["latest_chapter"]
        self.rating = series_info["bayesian_rating"]
        authorsList = []
        artistsList = []
        for author in series_info["authors"]:
            if author["type"] == "Author":
                authorsList.append(author["name"])
            if author["type"] == "Artist":
                artistsList.append(author["name"])
        self.authors = ", ".join(authorsList)
        self.artists = ", ".join(artistsList)