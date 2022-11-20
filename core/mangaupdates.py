import os
import numpy
import aiohttp

class RequestsMU:
    def __init__(self):
        self.username = os.environ.get("MU_USER")
        self.password = os.environ.get("MU_PASS")
        self.loginurl = "https://api.mangaupdates.com/v1/account/login"
        self.session = aiohttp.ClientSession()
    async def login(self):
        async with self.session.put(self.loginurl, json={"username": self.username, "password": self.password}) as resp:
            login = await resp.json()
            self.token = login["context"]["session_token"]
            self.headers = {"Authorization": f"Bearer {self.token}"}
            return self.headers
    async def get(self, url):
        await self.login()
        async with self.session.get(url, headers=self.headers) as resp:
            return await resp.json()
    async def put(self, url, data):
        await self.login()
        async with self.session.put(url, headers=self.headers, json=data) as resp:
            return await resp.json()
    async def post(self, url, data):
        await self.login()
        async with self.session.post(url, headers=self.headers, json=data) as resp:
            return await resp.json()
    async def close(self):
        await self.session.close()

class MangaUpdates:
    def __init__(self):
        self.rq = RequestsMU()

    async def convert_old_id(self, old_id):
        enc = numpy.base_repr(old_id, 36).lower()
        return enc

    async def convert_new_id(self, new_id):
        return int(new_id, 36)
    
    async def search_series(self, series_name):
        searchurl = f"https://api.mangaupdates.com/v1/series/search"
        search = await self.rq.post(searchurl, data={"search": series_name, "perpage": 10})
        return search
    
    async def series_info(self, series_id):
        infourl = f"https://api.mangaupdates.com/v1/series/{series_id}"
        info = await self.rq.get(infourl)
        return info

    async def search_groups(self, group_name):
        searchurl = f"https://api.mangaupdates.com/v1/groups/search"
        search = await self.rq.post(searchurl, data={"search": group_name, "perpage": 10})
        return search

    async def group_info(self, group_id):
        infourl = f"https://api.mangaupdates.com/v1/groups/{group_id}"
        info = await self.rq.get(infourl)
        return info

    async def series_groups(self, series_id):
        groupsurl = f"https://api.mangaupdates.com/v1/series/{series_id}/groups"
        groups = await self.rq.get(groupsurl)
        return groups
    async def close(self):
        await self.rq.close()