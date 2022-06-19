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

    async def convert_old_id(self, old_id):
        enc = numpy.base_repr(old_id, 36).lower()
        return enc

    async def convert_new_id(self, new_id):
        return int(new_id, 36)
    
    async def search_series(self, series_name):
        searchurl = f"https://api.mangaupdates.com/v1/series/search"
        search = requests.post(searchurl, data={"search": series_name, "perpage": 10}, headers=self.headers).json()
        return search
    
    async def series_info(self, series_id):
        infourl = f"https://api.mangaupdates.com/v1/series/{series_id}"
        info = requests.get(infourl, headers=self.headers).json()
        return info

    async def search_groups(self, group_name):
        searchurl = f"https://api.mangaupdates.com/v1/groups/search"
        search = requests.post(searchurl, data={"search": group_name, "perpage": 10}, headers=self.headers).json()
        return search

    async def group_info(self, group_id):
        infourl = f"https://api.mangaupdates.com/v1/groups/{group_id}"
        info = requests.get(infourl, headers=self.headers).json()
        return info

    async def series_groups(self, series_id):
        groupsurl = f"https://api.mangaupdates.com/v1/series/{series_id}/groups"
        groups = requests.get(groupsurl, headers=self.headers).json()
        return groups