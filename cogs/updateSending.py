import discord
from discord.ext import commands

import threading
import asyncio
import nest_asyncio
import time
from datetime import datetime

from core import update
from core import mongodb

class UpdateSending(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

        # todo: setup failsafe by saving last data to somewhere to check again if bot dies and needs to restart
        def checkForUpdates():
            old = update.getLatest()
            while True:
                new = update.getLatest()
                print("Checking for new updates! " + (str(datetime.now().strftime("%H:%M:%S"))))
                if new != old:
                    print("New update found!")
                    newMangas = [manga for manga in new if manga not in old]
                    for manga in newMangas:
                        asyncio.run_coroutine_threadsafe(notify(manga["title"], manga["chapter"], manga["scanGroup"], manga["link"]), bot.loop)
                time.sleep(60)
                old = new
            
        async def notify(title, chapter, group, link):
            if link != None:
                data = update.getAllData(link)
                allTitles = data["associatedNames"]
                image = data["image"]
            elif link == None:
                allTitles = [title]
                image = None

            serverNeed = False
            userNeed = False
            userList = []
            serverList = []
            for title in allTitles:
                serverWant = mongodb.mangaWanted(title, group, "server")
                userWant = mongodb.mangaWanted(title, group, "user")
                if userWant != None:
                    userNeed = True
                    for user in userWant:
                        userList.append({
                            "id": user,
                            "title": title
                        })
                if serverWant != None:
                    serverNeed = True
                    i = 0
                    for channel in serverWant.channelList:
                        serverList.append({
                            "serverid": serverWant.serverList[i],
                            "channelid": channel,
                            "title": title
                        })
                        i += 1
                        
            if userNeed == True:
                for user in userList:
                    userObject = await self.bot.fetch_user(user["id"])
                    mangaTitle = user["title"]
                    embed = discord.Embed(title=f"New {mangaTitle} chapter released!", description=f"There is a new `{mangaTitle}` chapter.", color=0x3083e3)
                    embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                    embed.add_field(name="Chapter", value=chapter, inline=True)
                    embed.add_field(name="Group", value=group, inline=True)
                    embed.add_field(name="Link", value=link, inline=False)
                    if image != None:
                        embed.set_image(url=image)
                    await userObject.send(embed=embed)
            if serverNeed == True:
                for server in serverList:
                    channelObject = self.bot.get_channel(server["channelid"])
                    mangaTitle = server["title"]
                    embed = discord.Embed(title=f"New {mangaTitle} chapter released!", description=f"There is a new `{mangaTitle}` chapter.", color=0x3083e3)
                    embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                    embed.add_field(name="Chapter", value=chapter, inline=True)
                    embed.add_field(name="Group", value=group, inline=True)
                    embed.add_field(name="Link", value=link, inline=False)
                    if image != None:
                        embed.set_image(url=image)
                    await channelObject.send(embed=embed)
            else:
                print(f"New manga not wanted. ({title})")

        nest_asyncio.apply()
        checkThread = threading.Thread(target=checkForUpdates)
        checkThread.start()

def setup(bot):
    bot.add_cog(UpdateSending(bot))