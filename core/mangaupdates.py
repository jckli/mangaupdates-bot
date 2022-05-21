import os
import requests
import numpy

class MangaUpdates:
    def __init__(self):
        username = os.environ.get("MU_USER")
        password = os.environ.get("MU_PASS")
        loginurl = "https://api.mangaupdates.com/v1/account/login"
        login = requests.put(loginurl, json={"username": username, "password": password}).json()
        self.token = login["context"]["session_token"]
        self.headers = {"Authorization": f"Bearer {self.token}"}

    def convert_old_id(self, old_id):
        enc = numpy.base_repr(old_id, 36).lower()
        return enc
    
    def search_series(self, series_name):
        searchurl = f"https://api.mangaupdates.com/v1/series/search"
        search = requests.post(searchurl, data={"search": series_name, "perpage": 10}, headers=self.headers).json()
        return search
    
    def series_info(self, series_id):
        infourl = f"https://api.mangaupdates.com/v1/series/{series_id}"
        info = requests.get(infourl, headers=self.headers).json()
        return info