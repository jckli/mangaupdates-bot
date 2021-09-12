import discord
from discord.ext import commands

import datetime
import threading
import asyncio

from core import update
from core import mongodb

class Manga(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    """
    async def notify(title, chapter, group, link):
        serverWant = mongodb.manga_wanted(title, "server")
        userWant = mongodb.manga_wanted(title, "user")
        if serverWant == True:
            
        if userWant == True:
            
        else:
            print("New manga not wanted lol")
    """

    def checkForUpdates(bot):
        old = update.getLatest()
        while True:
            new = update.getLatest()
            if new != old:
                newMangas = [manga for manga in new if manga not in old]
                for manga in newMangas:
                    asyncio.run_coroutine_threadsafe(notify(manga["title"], manga["chapter"], manga["group"], manga["link"]), bot.loop)
                old = new

    @commands.command()
    async def setup(self, ctx):
        await ctx.message.delete()
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        setupEmbedS = discord.Embed(title="Setup", color=0x3083e3, description="Do you want manga updates sent to your DMs or a server? (Type user or server)")
        sentEmbed = await ctx.send(embed=setupEmbedS)
        try:
            mode = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id)
        except asyncio.TimeoutError:
            await sentEmbed.delete()
            await ctx.send(embed=timeoutError, delete_after=5.0)
        else:
            await sentEmbed.delete()
            await mode.delete()
            if mode.content == "user":
                embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command +addmanga to add manga.")
                await ctx.send(embed=embedUser, delete_after=10.0)
            if mode.content == "server":
                embedServer = discord.Embed(title="Setup", color=0x3083e3, description="What channel should I use?")
                sentEmbedServer = await ctx.send(embed=embedServer)
                try:
                    channel = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id)
                except asyncio.TimeoutError:
                    await sentEmbedServer.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    await sentEmbedServer.delete()
                    await channel.delete()
                    serverid = ctx.message.guild.id
                    if mongodb.checkServerExist(serverid) == True:
                        alrfinishSS = discord.Embed(title="Setup", color=0x3083e3, description="This server is already setup. Run the command +addmanga to add manga.")
                        await ctx.send(embed=alrfinishSS, delete_after=10.0)
                    elif mongodb.checkServerExist(serverid) == False:
                        channelid = channel.channel_mentions[0].id
                        serverInfo = self.bot.get_guild(serverid)
                        mongodb.addServer(serverInfo.name, serverid, channelid)
                        finishServerSetup = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command +addmanga to add manga.")
                        await ctx.send(embed=finishServerSetup, delete_after=10.0)

def setup(bot):
    bot.add_cog(Manga(bot))