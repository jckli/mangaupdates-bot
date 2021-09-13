import discord
from discord.ext import commands

import datetime
import threading
import asyncio
import nest_asyncio
import threading

from core import update
from core import mongodb

class Manga(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
    
        async def notify(title, chapter, group, link, image):
            serverWant = mongodb.mangaWanted(title, "server")
            userWant = mongodb.mangaWanted(title, "user")
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

        def checkForUpdates():
            old = update.getLatest()
            while True:
                new = update.getLatest()
                if new != old:
                    newMangas = [manga for manga in new if manga not in old]
                    for manga in newMangas:
                        asyncio.run_coroutine_threadsafe(notify(manga["title"], manga["chapter"], manga["group"], manga["link"], manga["image"]), bot.loop)
                    old = new

        nest_asyncio.apply()
        db_thread = threading.Thread(target=checkForUpdates)
        db_thread.start()

    @commands.command()
    async def setup(self, ctx):
        await ctx.message.delete()
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        setupEmbedS = discord.Embed(title="Setup", color=0x3083e3, description="Do you want manga updates sent to your DMs or a server? (Type user or server)")
        sentEmbed = await ctx.send(embed=setupEmbedS)
        try:
            mode = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
        except asyncio.TimeoutError:
            await sentEmbed.delete()
            await ctx.send(embed=timeoutError, delete_after=5.0)
        else:
            await sentEmbed.delete()
            await mode.delete()
            if mode.content == "user":
                userid = ctx.message.author.id
                if mongodb.checkUserExist(userid) == True:
                    alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup. Run the command +addmanga to add manga.")
                    await ctx.send(embed=alrfinishUS, delete_after=10.0)
                elif mongodb.checkUserExist(userid) == False:
                    userInfo = await self.bot.fetch_user(userid)
                    username = f"{userInfo.name}#{userInfo.discriminator}"
                    mongodb.addUser(username, userid)
                    embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command +addmanga to add manga.")
                    await ctx.send(embed=embedUser, delete_after=10.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong.")
                    await ctx.send(embed=completeError, delete_after=5.0)
            if mode.content == "server":
                serverid = ctx.message.guild.id
                if mongodb.checkServerExist(serverid) == True:
                    alrfinishSS = discord.Embed(title="Setup", color=0x3083e3, description="This server is already setup. Run the command +addmanga to add manga.")
                    await ctx.send(embed=alrfinishSS, delete_after=10.0)
                elif mongodb.checkServerExist(serverid) == False:
                    embedServer = discord.Embed(title="Setup", color=0x3083e3, description="What channel should I use?")
                    sentEmbedServer = await ctx.send(embed=embedServer)
                    try:
                        channel = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                    except asyncio.TimeoutError:
                        await sentEmbedServer.delete()
                        await ctx.send(embed=timeoutError, delete_after=5.0)
                    else:
                        await sentEmbedServer.delete()
                        await channel.delete()
                        channelid = channel.channel_mentions[0].id
                        serverInfo = self.bot.get_guild(serverid)
                        mongodb.addServer(serverInfo.name, serverid, channelid)
                        finishServerSetup = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command +addmanga to add manga.")
                        await ctx.send(embed=finishServerSetup, delete_after=10.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong.")
                    await ctx.send(embed=completeError, delete_after=5.0)

def setup(bot):
    bot.add_cog(Manga(bot))