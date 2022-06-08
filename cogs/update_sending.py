import discord
from discord.ext import commands, tasks
from discord.commands import slash_command
import asyncio
import nest_asyncio
from datetime import datetime
from core.mongodb import Mongo
from core.rss import RSSParser
from core.mangaupdates import MangaUpdates

mongo = Mongo()
mangaupdates = MangaUpdates()
rss = RSSParser()

class UpdateSending(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.old = None
        self.check_for_updates.start()

    def cog_unload(self):
        self.check_for_updates.cancel()

    # todo: setup failsafe by saving last data to somewhere to check again if bot dies and needs to restart (redis maybe?)
    @tasks.loop(seconds=15)
    async def check_for_updates(self):
        new = await rss.parse_feed()
        print("Checking for new updates! " + (str(datetime.now().strftime("%H:%M:%S"))))
        if new != self.old:
            print("New update found!")
            new_mangas = [manga for manga in new if manga not in self.old]
            for manga in new_mangas:
                asyncio.run_coroutine_threadsafe(self.notify(manga["title"], manga["chapter"], manga["scan_group"], manga["link"]), self.bot.loop)
        self.old = new
    
    @check_for_updates.before_loop
    async def before_printer(self):
        self.old = await rss.parse_feed()
        await self.bot.wait_until_ready()

    async def notify(title, chapter, scan_group, link):
        if link is not None:
            link = link.partition("https://www.mangaupdates.com/series/")[2]
            mangaid = link.partition("/")[0]
            mangaid = await mangaupdates.convert_new_id(mangaid)
            data = mangaupdates.series_info(mangaid)
            associatedTitles = [manga["title"] for manga in data["associated"]]
            allTitles = [data["title"]] + associatedTitles
            image = data["image"]["url"]["original"]
        elif link is None:
            allTitles = [title]
            image = None

        if scan_group is not None:
            if "&" in group:
                group = group.split("&")
                group = [x.strip(' ') for x in group]
            else:
                group = [scan_group]
            sgs = []
            for g in group:
                scan_groups_search = await mangaupdates.search_groups(scan_group)
                scan_group = scan_groups_search["results"][0]
                sgs.append(scan_group["record"])

        serverNeed = False
        userNeed = False
        userList = []
        serverList = []

def setup(bot):
    bot.add_cog(UpdateSending(bot))