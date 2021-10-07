import discord
from discord.ext import commands
from discord.ext.commands import BotMissingPermissions, NoPrivateMessage

import asyncio
import pymanga
import validators
import time

from core import mongodb, update

class Manga(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    """
    async def cog_command_error(self, ctx, error):
        if isinstance(error, BotMissingPermissions):
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
            await ctx.send(embed=permissionError, delete_after=10.0)
        elif isinstance(error, NoPrivateMessage):
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="This command cannot be ran in DMs. Please re-run this command in a server.")
            await ctx.send(embed=permissionError, delete_after=10.0)
        else:
            pass
    """
        
    @commands.command(name="addmanga")
    async def addmanga(self, ctx, *, arg=None):
        userid = ctx.message.author.id
        userExist = mongodb.checkUserExist(userid)
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
                await ctx.message.delete()
                mode = arg
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
                if mode == "user":
                    if userExist == True:
                        addMangaEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="What manga do you want to add? (Can also use mangaupdates.com link)")
                        sentEmbedAddManga = await ctx.send(embed=addMangaEmbed)
                        try:
                            manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                        except asyncio.TimeoutError:
                            await sentEmbedAddManga.delete()
                            await ctx.send(embed=timeoutError, delete_after=5.0)
                        else:
                            await sentEmbedAddManga.delete()
                            await manga.delete()
                            if validators.url(manga.content) == True:
                                mangaTitle = update.getTitle(manga.content)
                                link = manga.content
                            elif validators.url(manga.content) != True:
                                searchRaw = pymanga.api.search(manga.content)
                                description = "Type the number of the manga you want to add.\n"
                                searchNames = []
                                if searchRaw["series"] == []:
                                    resultError = discord.Embed(title="Error", color=0xff4f4f, description="No mangas were found.")
                                    await ctx.send(embed=resultError, delete_after=5.0)
                                    return
                                elif searchRaw["series"] != []:
                                    i = 1
                                    for result in searchRaw["series"]:
                                        name = result["name"]
                                        year = result["year"]
                                        rating = result["rating"]
                                        description += f"{i}. {name} ({year}, Rating: {rating})\n"
                                        searchNames.append(name)
                                        i += 1
                                    searchEmbed = discord.Embed(title="Search Results", color=0x3083e3, description=description)
                                    sentEmbedSearch = await ctx.send(embed=searchEmbed)
                                    try:
                                        search = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                                    except asyncio.TimeoutError:
                                        await sentEmbedSearch.delete()
                                        await ctx.send(embed=timeoutError, delete_after=5.0)
                                    else:
                                        await sentEmbedSearch.delete()
                                        await search.delete()
                                        if search.content.isnumeric() is True and int(search.content) in range(1, 11):
                                            mangaTitle = searchNames[int(search.content)-1]
                                            link = update.getLink(mangaTitle)
                                        else:
                                            countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `10`")
                                            await ctx.send(embed=countError, delete_after=5.0)
                                            return
                                else:
                                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                                    await ctx.send(embed=completeError, delete_after=5.0)
                                    return
                            else:
                                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                                await ctx.send(embed=completeError, delete_after=5.0)
                                return
                            mangaInDb = mongodb.checkMangaAlreadyWithinDb(userid, link, "user")
                            if mangaInDb == True:
                                mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to your list.")
                                await ctx.send(embed=mangaExist, delete_after=10.0)
                                return
                            elif mangaInDb == False:
                                mongodb.addManga(userid, mangaTitle, link, "user")
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
                    serverid = ctx.message.guild.id
                    serverExist = mongodb.checkServerExist(serverid)
                    if ctx.author.guild_permissions.administrator == True:
                        if serverExist == True:
                            addMangaEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="What manga do you want to add? (Can also use mangaupdates.com link)")
                            sentEmbedAddManga = await ctx.send(embed=addMangaEmbed)
                            try:
                                manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                            except asyncio.TimeoutError:
                                await sentEmbedAddManga.delete()
                                await ctx.send(embed=timeoutError, delete_after=5.0)
                            else:
                                await sentEmbedAddManga.delete()
                                await manga.delete()
                                if validators.url(manga.content) == True:
                                    mangaTitle = update.getTitle(manga.content)
                                    link = manga.content
                                elif validators.url(manga.content) != True:
                                    searchRaw = pymanga.api.search(manga.content)
                                    description = "Type the number of the manga you want to add.\n"
                                    searchNames = []
                                    if searchRaw["series"] == []:
                                        resultError = discord.Embed(title="Error", color=0xff4f4f, description="No mangas were found.")
                                        await ctx.send(embed=resultError, delete_after=5.0)
                                        return
                                    elif searchRaw["series"] != []:
                                        i = 1
                                        for result in searchRaw["series"]:
                                            name = result["name"]
                                            year = result["year"]
                                            rating = result["rating"]
                                            description += f"{i}. {name} ({year}, Rating: {rating})\n"
                                            searchNames.append(name)
                                            i += 1
                                        searchEmbed = discord.Embed(title="Search Results", color=0x3083e3, description=description)
                                        sentEmbedSearch = await ctx.send(embed=searchEmbed)
                                        try:
                                            search = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                                        except asyncio.TimeoutError:
                                            await sentEmbedSearch.delete()
                                            await ctx.send(embed=timeoutError, delete_after=5.0)
                                        else:
                                            await sentEmbedSearch.delete()
                                            await search.delete()
                                            if search.content.isnumeric() is True and int(search.content) in range(1, 11):
                                                mangaTitle = searchNames[int(search.content)-1]
                                                link = update.getLink(mangaTitle)
                                            else:
                                                countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `10`")
                                                await ctx.send(embed=countError, delete_after=5.0)
                                                return
                                    else:
                                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                                        await ctx.send(embed=completeError, delete_after=5.0)
                                        return
                                else:
                                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                                    await ctx.send(embed=completeError, delete_after=5.0)
                                    return
                                mangaInDb = mongodb.checkMangaAlreadyWithinDb(serverid, link, "server")
                                if mangaInDb == True:
                                    mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to the server's list.")
                                    await ctx.send(embed=mangaExist, delete_after=10.0)
                                elif mangaInDb == False:
                                    mongodb.addManga(serverid, mangaTitle, link, "server")
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
        else:
            if userExist == True:
                addMangaEmbed = discord.Embed(title="Add Manga", color=0x3083e3, description="What manga do you want to add? (Can also use mangaupdates.com link)")
                await ctx.send(embed=addMangaEmbed)
                try:
                    manga = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                except asyncio.TimeoutError:
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    if validators.url(manga.content) == True:
                        mangaTitle = update.getTitle(manga.content)
                        link = manga.content
                    elif validators.url(manga.content) != True:
                        searchRaw = pymanga.api.search(manga.content)
                        description = "Type the number of the manga you want to add.\n"
                        searchNames = []
                        if searchRaw["series"] == []:
                            resultError = discord.Embed(title="Error", color=0xff4f4f, description="No mangas were found.")
                            await ctx.send(embed=resultError, delete_after=5.0)
                            return
                        elif searchRaw["series"] != []:
                            i = 1
                            for result in searchRaw["series"]:
                                name = result["name"]
                                year = result["year"]
                                rating = result["rating"]
                                description += f"{i}. {name} ({year}, Rating: {rating})\n"
                                searchNames.append(name)
                                i += 1
                            searchEmbed = discord.Embed(title="Search Results", color=0x3083e3, description=description)
                            await ctx.send(embed=searchEmbed)
                            try:
                                search = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                            except asyncio.TimeoutError:
                                await ctx.send(embed=timeoutError, delete_after=5.0)
                            else:
                                if search.content.isnumeric() is True and int(search.content) in range(1, 11):
                                    mangaTitle = searchNames[int(search.content)-1]
                                    link = update.getLink(mangaTitle)
                                else:
                                    countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `10`")
                                    await ctx.send(embed=countError, delete_after=5.0)
                                    return
                        else:
                            completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                            await ctx.send(embed=completeError, delete_after=5.0)
                            return
                    else:
                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                        await ctx.send(embed=completeError, delete_after=5.0)
                        return
                    mangaInDb = mongodb.checkMangaAlreadyWithinDb(userid, link, "user")
                    if mangaInDb == True:
                        mangaExist = discord.Embed(title="Add Manga", color=0x3083e3, description="This manga is already added to your list.")
                        await ctx.send(embed=mangaExist)
                        return
                    elif mangaInDb == False:
                        mongodb.addManga(userid, mangaTitle, link, "user")
                        mangaAdded = discord.Embed(title="Add Manga", color=0x3083e3, description="Manga succesfully added.")
                        await ctx.send(embed=mangaAdded)
                        return
            elif userExist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first.")
                await ctx.send(embed=setupError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)

    @commands.command(name="removemanga")
    async def removemanga(self, ctx, *, arg=None):
        userid = ctx.message.author.id
        userExist = mongodb.checkUserExist(userid)
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
                await ctx.message.delete()
                mode = arg
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
                if mode == "user":
                    if userExist == True:
                        mangaList = mongodb.getMangaList(userid, "user")
                        i = 1
                        description = "Type the number of the manga you want to remove.\n"
                        for manga in mangaList:
                            description += f"{i}. {manga}\n"
                            i += 1
                        removeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description=description)
                        sentEmbedRemove = await ctx.send(embed=removeEmbed)
                        try:
                            remove = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                        except asyncio.TimeoutError:
                            await sentEmbedRemove.delete()
                            await ctx.send(embed=timeoutError, delete_after=5.0)
                        else:
                            await sentEmbedRemove.delete()
                            await remove.delete()
                            if remove.content.isnumeric() is True and int(remove.content) in range(1, i):
                                mangaTitle = mangaList[int(remove.content)-1]
                                mongodb.removeManga(userid, mangaTitle, "user")
                                mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
                                await ctx.send(embed=mangaRemoved, delete_after=10.0)
                            else:
                                countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `{}`".format(i-1))
                                await ctx.send(embed=countError, delete_after=5.0)
                                return
                    elif userExist == False:
                        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and add some manga using the `+addmanga` command.")
                        await ctx.send(embed=setupError, delete_after=5.0)
                    else:
                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                        await ctx.send(embed=completeError, delete_after=5.0)
                elif mode == "server":
                    serverid = ctx.message.guild.id
                    serverExist = mongodb.checkServerExist(serverid)
                    if ctx.author.guild_permissions.administrator == True:
                        if serverExist == True:
                            mangaList = mongodb.getMangaList(serverid, "server")
                            i = 1
                            description = "Type the number of the manga you want to remove.\n"
                            for manga in mangaList:
                                description += f"{i}. {manga}\n"
                                i += 1
                            removeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description=description)
                            sentEmbedRemove = await ctx.send(embed=removeEmbed)
                            try:
                                remove = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                            except asyncio.TimeoutError:
                                await sentEmbedRemove.delete()
                                await ctx.send(embed=timeoutError, delete_after=5.0)
                            else:
                                await sentEmbedRemove.delete()
                                await remove.delete()
                                if remove.content.isnumeric() is True and int(remove.content) in range(1, i):
                                    mangaTitle = mangaList[int(remove.content)-1]
                                    mongodb.removeManga(serverid, mangaTitle, "server")
                                    mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
                                    await ctx.send(embed=mangaRemoved, delete_after=10.0)
                                else:
                                    countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `{}`".format(i-1))
                                    await ctx.send(embed=countError, delete_after=5.0)
                                    return
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
        else:
            if userExist == True:
                mangaList = mongodb.getMangaList(userid, "user")
                i = 1
                description = "Type the number of the manga you want to remove.\n"
                for manga in mangaList:
                    description += f"{i}. {manga}\n"
                    i += 1
                removeEmbed = discord.Embed(title="Remove Manga", color=0x3083e3, description=description)
                await ctx.send(embed=removeEmbed)
                try:
                    remove = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                except asyncio.TimeoutError:
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    if remove.content.isnumeric() is True and int(remove.content) in range(1, i):
                        mangaTitle = mangaList[int(remove.content)-1]
                        mongodb.removeManga(userid, mangaTitle, "user")
                        mangaRemoved = discord.Embed(title="Remove Manga", color=0x3083e3, description="Manga succesfully removed.")
                        await ctx.send(embed=mangaRemoved)
                    else:
                        countError = discord.Embed(title="Error", color=0xff4f4f, description="You didn't select a number from `1` to `{}`".format(i-1))
                        await ctx.send(embed=countError)
                        return
            elif userExist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and add some manga using the `+addmanga` command.")
                await ctx.send(embed=setupError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)

    
    @commands.command(name="mangalist")
    async def mangalist(self, ctx, *, arg=None):
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
                mode = arg
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
                                noMangaError.set_thumbnail(url = iconUrl)
                                await ctx.send(embed=noMangaError)
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
                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                        await ctx.send(embed=completeError, delete_after=5.0)
                else:
                    modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
                    await ctx.send(embed=modeError, delete_after=5.0)
        else:
            userid = ctx.message.author.id
            name = ctx.message.author
            iconUrl = ctx.message.author.avatar_url
            exist = mongodb.checkUserExist(userid)
            if exist == True:
                try:
                    mangaList = mongodb.getMangaList(userid, "user")
                    if mangaList == None:
                        noMangaError = discord.Embed(title=f"{name}'s Manga List", color=0x3083e3, description="You have added no manga to your list.")
                        noMangaError.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                        noMangaError.set_thumbnail(url = iconUrl)
                        await ctx.send(embed=noMangaError)
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
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)


    @commands.command(name="clearmanga")
    async def clearmanga(self, ctx, *, arg=None):
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
                await ctx.message.delete()
                mode = arg
                modeEntry = False
                if (mode == None) or (mode != "server" and mode != "user"):
                    modeEntry = True
                    modeEmbed = discord.Embed(title="Manga List", color=0x3083e3, description="Do you want to clear your manga list or this server's manga list (Type user or server)")
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
                elif mode == "server":
                    givenid = ctx.message.guild.id
                    name = ctx.message.guild.name
                if mode == "user" or mode == "server":
                    if mode == "user":
                        exist = mongodb.checkUserExist(givenid)
                    elif mode == "server":
                        exist = mongodb.checkServerExist(givenid)
                    if exist == True:
                        try:
                            mangaList = mongodb.getMangaList(givenid, mode)
                            if mangaList == None:
                                noMangaError = discord.Embed(title=f"Cannot Clear {name}'s Manga list", color=0x3083e3, description="You have added no manga to your list.")
                                noMangaError.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                                await ctx.send(embed=noMangaError, delete_after=5.0)
                            else:
                                for manga in mangaList:
                                    mongodb.removeManga(givenid, manga, mode)
                                mangaListEmbed = discord.Embed(title="Cleared", color=0x3083e3, description=f"Succesfully cleared {name}'s manga list.")
                                mangaListEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                                await ctx.send(embed=mangaListEmbed, delete_after=5.0)
                        except:
                            completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                            await ctx.send(embed=completeError, delete_after=5.0)
                    elif exist == False:
                        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and have some manga before clearing the list.")
                        await ctx.send(embed=setupError, delete_after=5.0)
                    else:
                        completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                        await ctx.send(embed=completeError, delete_after=5.0)
                else:
                    modeError = discord.Embed(title="Error", color=0xff4f4f, description="You did not type in either `user` or `server`.")
                    await ctx.send(embed=modeError, delete_after=5.0)
        else:
            userid = ctx.message.author.id
            name = ctx.message.author
            exist = mongodb.checkUserExist(userid)
            if exist == True:
                try:
                    mangaList = mongodb.getMangaList(userid, "user")
                    if mangaList == None:
                        noMangaError = discord.Embed(title=f"Cannot Clear {name}'s Manga list", color=0x3083e3, description="You have added no manga to your list.")
                        noMangaError.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                        await ctx.send(embed=noMangaError)
                    else:
                        for manga in mangaList:
                            mongodb.removeManga(userid, manga, "user")
                        mangaListEmbed = discord.Embed(title="Cleared", color=0x3083e3, description=f"Succesfully cleared {name}'s manga list.")
                        mangaListEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar_url)
                        await ctx.send(embed=mangaListEmbed)
                except:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            elif exist == False:
                setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the `+setup` command first and have some manga before clearing the list.")
                await ctx.send(embed=setupError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)

def setup(bot):
    bot.add_cog(Manga(bot))