import discord
from discord.ext import commands

import asyncio

from core import mongodb

class Manga(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

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
                    alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup. Run the command `+addmanga` to add manga.")
                    await ctx.send(embed=alrfinishUS, delete_after=10.0)
                elif mongodb.checkUserExist(userid) == False:
                    userInfo = await self.bot.fetch_user(userid)
                    username = f"{userInfo.name}#{userInfo.discriminator}"
                    mongodb.addUser(username, userid)
                    embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command `+addmanga` to add manga.")
                    await ctx.send(embed=embedUser, delete_after=10.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            elif mode.content == "server":
                serverid = ctx.message.guild.id
                if mongodb.checkServerExist(serverid) == True:
                    alrfinishSS = discord.Embed(title="Setup", color=0x3083e3, description="This server is already setup. Run the command `+addmanga` to add manga.")
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
                        finishServerSetup = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command `+addmanga` to add manga.")
                        await ctx.send(embed=finishServerSetup, delete_after=10.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            else:
                modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
                await ctx.send(embed=modeError, delete_after=5.0)

# Search it (https://github.com/alisww/py-manga), link capability, command ailias of +addmanga server or +addmanga user
    @commands.command()
    async def addmanga(self, ctx, *, arg=None):
        await ctx.message.delete()
        userid = ctx.message.author.id
        serverid = ctx.message.guild.id
        mode = arg
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        modeEntry = False
        if (mode == None) or (mode != "server" and mode != "user"):
            modeEntry = True
            modeEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="Do you want this manga added to your list or this server's list? (Type user or server)")
            sentEmbedMode = await ctx.send(embed=modeEmbed)
            try:
                modeObject = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                mode = modeObject.content
            except asyncio.TimeoutError:
                await sentEmbedMode.delete()
                await ctx.send(embed=timeoutError, delete_after=5.0)
                return
        if modeEntry == True:
            await sentEmbedMode.delete()
            await modeObject.delete()
        serverExist = mongodb.checkServerExist(serverid)
        userExist = mongodb.checkUserExist(userid)
        if mode == "user":
            if userExist == True:
                addMangaEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="What manga do you want to add? (Please type exact title name with correct punctuation)")
                sentEmbedAddManga = await ctx.send(embed=addMangaEmbed)
                try:
                    manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                except asyncio.TimeoutError:
                    await sentEmbedAddManga.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    await sentEmbedAddManga.delete()
                    await manga.delete()
                    mangaInDb = mongodb.checkMangaAlreadyWithinDb(userid, manga.content, "user")
                    if mangaInDb == True:
                        mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to your list.")
                        await ctx.send(embed=mangaExist, delete_after=10.0)
                        return
                    elif mangaInDb == False:
                        mongodb.addManga(userid, manga.content, "user")
                        mangaAdded = discord.Embed(title="Add Manga", color=0x3083e3, description="Manga succesfully added.")
                        await ctx.send(embed=mangaAdded, delete_after=10.0)
                        return
            elif userExist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first.")
                await ctx.send(embed=setupError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
        elif mode == "server":
            if ctx.author.guild_permissions.administrator == True:
                if serverExist == True:
                    addMangaEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="What manga do you want to add? (Please type exact title name with correct punctuation)")
                    sentEmbedAddManga = await ctx.send(embed=addMangaEmbed)
                    try:
                        manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                    except asyncio.TimeoutError:
                        await sentEmbedAddManga.delete()
                        await ctx.send(embed=timeoutError, delete_after=5.0)
                    else:
                        await sentEmbedAddManga.delete()
                        await manga.delete()
                        mangaInDb = mongodb.checkMangaAlreadyWithinDb(serverid, manga.content, "server")
                        if mangaInDb == True:
                            mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to the server's list.")
                            await ctx.send(embed=mangaExist, delete_after=10.0)
                        elif mangaInDb == False:
                            mongodb.addManga(serverid, manga.content, "server")
                            mangaAdded = discord.Embed(title="Add Manga", color=0x3083e3, description="Manga succesfully added.")
                            await ctx.send(embed=mangaAdded, delete_after=10.0)
                elif serverExist == False:
                    setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first.")
                    await ctx.send(embed=setupError, delete_after=5.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            else:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to add manga. You need `Administrator` permission to use this.")
                await ctx.send(embed=permissionError, delete_after=5.0)
        else:
            modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
            await ctx.send(embed=modeError, delete_after=5.0)

    @commands.command()
    async def removemanga(self, ctx, *, arg=None):
        await ctx.message.delete()
        userid = ctx.message.author.id
        serverid = ctx.message.guild.id
        mode = arg
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        modeEntry = False
        if (mode == None) or (mode != "server" and mode != "user"):
            modeEntry = True
            modeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description="Do you want this manga removed from your list or this server's list? (Type user or server)")
            sentEmbedMode = await ctx.send(embed=modeEmbed)
            try:
                modeObject = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                mode = modeObject.content
            except asyncio.TimeoutError:
                await sentEmbedMode.delete()
                await ctx.send(embed=timeoutError, delete_after=5.0)
                return
        if modeEntry == True:
            await sentEmbedMode.delete()
            await modeObject.delete()
        serverExist = mongodb.checkServerExist(serverid)
        userExist = mongodb.checkUserExist(userid)
        if mode == "user":
            if userExist == True:
                remMangaEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description="What manga do you want to remove? (Please type exact title name with correct punctuation)")
                sentEmbedRemManga = await ctx.send(embed=remMangaEmbed)
                try:
                    manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                except asyncio.TimeoutError:
                    await sentEmbedRemManga.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    await sentEmbedRemManga.delete()
                    await manga.delete()
                    mangaInDb = mongodb.checkMangaAlreadyWithinDb(userid, manga.content, "user")
                    if mangaInDb == True:
                        mongodb.removeManga(userid, manga.content, "user")
                        mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
                        await ctx.send(embed=mangaRemoved, delete_after=10.0)
                    elif mangaInDb == False:
                        mangaNotExist = discord.Embed(title="Remove Manga", color=0x3083e3, description="This manga is not in your list.")
                        await ctx.send(embed=mangaNotExist, delete_after=10.0)
                    else:
                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                        await ctx.send(embed=completeError, delete_after=5.0)
            elif userExist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and add some manga using the `+addmanga` command.")
                await ctx.send(embed=setupError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
        elif mode == "server":
            if ctx.author.guild_permissions.administrator == True:
                if serverExist == True:
                    remMangaEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description="What manga do you want to remove? (Please type exact title name with correct punctuation)")
                    sentEmbedRemManga = await ctx.send(embed=remMangaEmbed)
                    try:
                        manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                    except asyncio.TimeoutError:
                        await sentEmbedRemManga.delete()
                        await ctx.send(embed=timeoutError, delete_after=5.0)
                    else:
                        await sentEmbedRemManga.delete()
                        await manga.delete()
                        mangaInDb = mongodb.checkMangaAlreadyWithinDb(serverid, manga.content, "server")
                        if mangaInDb == True:
                            mongodb.removeManga(serverid, manga.content, "server")
                            mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
                            await ctx.send(embed=mangaRemoved, delete_after=10.0)
                        elif mangaInDb == False:
                            mangaNotExist = discord.Embed(title="Remove Manga", color=0x3083e3, description="This manga is not in the server's list.")
                            await ctx.send(embed=mangaNotExist, delete_after=10.0)
                        else:
                            completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                            await ctx.send(embed=completeError, delete_after=5.0)
                elif serverExist == False:
                    setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and add some manga using the `+addmanga` command.")
                    await ctx.send(embed=setupError, delete_after=5.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            else:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to add manga. You need `Administrator` permission to use this.")
                await ctx.send(embed=permissionError, delete_after=5.0)
        else:
            modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
            await ctx.send(embed=modeError, delete_after=5.0)
    
    @commands.command()
    async def mangalist(self, ctx, *, arg=None):
        mode = arg
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        modeEntry = False
        if (mode == None) or (mode != "server" and mode != "user"):
            modeEntry = True
            modeEmbed = discord.Embed(title="Manga List", color=0x3083e3, description="Do you want to see your manga list or this server's manga list (Type user or server)")
            sentEmbedMode = await ctx.send(embed=modeEmbed)
            try:
                modeObject = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                mode = modeObject.content
            except asyncio.TimeoutError:
                await sentEmbedMode.delete()
                await ctx.send(embed=timeoutError, delete_after=5.0)
                return
        if modeEntry == True:
            await sentEmbedMode.delete()
            await modeObject.delete()
        if mode == "user":
            givenid = ctx.message.author.id
            name = ctx.message.author
            iconUrl = ctx.message.author.avatar_url
        elif mode == "server":
            givenid = ctx.message.guild.id
            name = ctx.message.guild.name
            iconUrl = ctx.message.guild.icon_url
        if mode == "user" or mode == "server":
            if mode == "user":
                exist = mongodb.checkUserExist(givenid)
            elif mode == "server":
                exist = mongodb.checkServerExist(givenid)
            if exist == True:
                try:
                    mangaList = mongodb.getMangaList(givenid, mode)
                    if mangaList == None:
                        noMangaError = discord.Embed(title=f"{name}'s Manga List", color=0x3083e3, description="You have added no manga to your list.")
                        noMangaError.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                        await ctx.send(embed=noMangaError, delete_after=5.0)
                    else:
                        description = ""
                        for manga in mangaList:
                            description += f"• {manga}\n"
                        mangaListEmbed = discord.Embed(title=f"{name}'s Manga List", color=0x3083e3, description=description)
                        mangaListEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                        mangaListEmbed.set_thumbnail(url = iconUrl)
                        await ctx.send(embed=mangaListEmbed)
                except:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            elif exist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and add some manga using the `+addmanga` command.")
                await ctx.send(embed=setupError, delete_after=5.0)
        else:
            modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
            await ctx.send(embed=modeError, delete_after=5.0)

def setup(bot):
    bot.add_cog(Manga(bot))