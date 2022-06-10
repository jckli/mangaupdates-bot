import discord
from discord.ext import commands, tasks
import asyncio
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

    async def notify(self, title, chapter, scan_group, link):
        if link:
            link = link.partition("https://www.mangaupdates.com/series/")[2]
            mangaid = link.partition("/")[0]
            mangaid = await mangaupdates.convert_new_id(mangaid)
            data = await mangaupdates.series_info(mangaid)
            image = data["image"]["url"]["original"]
        else:
            image = None

        if scan_group:
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

        if link:
            serverWant = mongo.manga_wanted_server(sgs, manga_id=mangaid)
            userWant = mongo.manga_wanted_user(sgs, manga_id=mangaid)
        else:
            serverWant = mongo.manga_wanted_server(sgs, manga_title=title)
            userWant = mongo.manga_wanted_user(sgs, manga_title=title)
        
        if userWant or serverWant:
            print(f"Manga Wanted ({title})")

        if sgs[0]["social"]["site"]:
            scanLink = sgs[0]["social"]["site"]
        elif sgs[0]["social"]["discord"]:
            scanLink = sgs[0]["social"]["discord"]
        else:
            scanLink = sgs[0]["url"]
        
        if userWant:
            for user in userWant:
                userObject = await self.bot.fetch_user(user["id"])
                userEmbed = discord.Embed(title=f"New {user['title']} chapter released!", url=link, description=f"There is a new `{user['title']}` chapter.", color=0x3083e3)
                userEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                userEmbed.add_field(name="Chapter", value=chapter, inline=True)
                userEmbed.add_field(name="Group", value=group, inline=True)
                userEmbed.add_field(name="Scanlator Link", value=scanLink, inline=False)
                if image != None:
                    userEmbed.set_image(url=image)
                await userObject.send(embed=userEmbed)
        else:
            print(f"New manga not wanted. (User: {title})")
        if serverWant:
            for server in serverWant:
                channelObject = self.bot.get_channel(server["channelid"])
                channelEmbed = discord.Embed(title=f"New {server['title']} chapter released!", url=link, description=f"There is a new `{server['title']}` chapter.", color=0x3083e3)
                channelEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                channelEmbed.add_field(name="Chapter", value=chapter, inline=True)
                channelEmbed.add_field(name="Group", value=group, inline=True)
                channelEmbed.add_field(name="Scanlator Link", value=scanLink, inline=False)
                if image != None:
                    channelEmbed.set_image(url=image)
                await channelObject.send(embed=channelEmbed)
        else:
            print(f"New manga not wanted. (Server: {title})")

def setup(bot):
    bot.add_cog(UpdateSending(bot))