import discord
from discord.ext import commands

import datetime
import threading
import asyncio

from core import update

class Manga(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    async def notify(title, chapter, group, link):
        if title == "Manga Not Found":
            print("not a manga wanted lol")
        else:
            for 

    def checkForUpdates(bot):
        old = update.getLatest()
        while True:
            new = update.getLatest()
            if new != old:
                newMangas = [manga for manga in new if manga not in old]
                for manga in newMangas:
                    asyncio.run_coroutine_threadsafe(notify(manga["title"], manga["chapter"], manga["group"], manga["link"]), bot.loop)
                old = new

def setup(bot):
    bot.add_cog(Manga(bot))