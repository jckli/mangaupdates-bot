import discord
from discord.ext import commands

import threading
import asyncio
import nest_asyncio
import threading
import time

from core import update
from core import mongodb

class UpdateSending(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

        # todo: setup failsafe by saving last data to somewhere to check again if bot dies and needs to restart
        def checkForUpdates():
            old = update.getLatest()
            value = True
            while value == True:
                new = update.getLatest()
                print("Checking for new updates!")
                if new != old:
                    print("New update found!")
                    newMangas = [manga for manga in new if manga not in old]
                    for manga in newMangas:
                        asyncio.run_coroutine_threadsafe(notify(manga["title"], manga["chapter"], manga["group"], manga["link"]), bot.loop)
                time.sleep(60)
                old = new
            
        async def notify(title, chapter, group, link):
            serverWant = mongodb.mangaWanted(title, "server")
            userWant = mongodb.mangaWanted(title, "user")
            image = update.getImage(link)
            embed = discord.Embed(title=f"New {title} chapter released!", description=f"There is a new {title} chapter released", color=0x3083e3)
            embed.add_field(name="Chapter", value=chapter, inline=True)
            embed.add_field(name="Group", value=group, inline=True)
            embed.add_field(name="Link", value=link, inline=False)
            embed.set_image(url=image)
            if userWant != None:
                for user in userWant:
                    userObject = await self.bot.fetch_user(user)
                    await userObject.send(embed=embed)
                    break
            if serverWant.serverList != None:
                for channel in serverWant.channelList:
                    channelObject = self.bot.get_channel(channel)
                    await channelObject.send(embed=embed)
                    break
            else:
                print("New manga not wanted lol")

        nest_asyncio.apply()
        checkThread = threading.Thread(target=checkForUpdates)
        checkThread.start()

def setup(bot):
    bot.add_cog(UpdateSending(bot))