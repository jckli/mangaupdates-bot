import discord
from discord.ext import commands

import json
import asyncio

from core import mongodb, update

with open("config.json", "r") as f:
    config = json.load(f)

class Mode(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None

    @discord.ui.button(label=f'User (DMs)', style=discord.ButtonStyle.grey)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "user"
        self.stop()

    @discord.ui.button(label=f'Server', style=discord.ButtonStyle.grey)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "server"
        self.stop()

class Confirm(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None

    sep = '\u2001'

    @discord.ui.button(label=f'{sep*6}Confirm{sep*6}', style=discord.ButtonStyle.green)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = True
        self.stop()

    @discord.ui.button(label=f'{sep*6}Cancel{sep*6}', style=discord.ButtonStyle.red)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = False
        self.stop()

class General(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="setup")
    async def setup(self, ctx, *, arg=None):
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            modeView = Mode()
            mode = arg
            modeEntry = False
            if (mode == None) or (mode != "server" and mode != "user"):
                modeEntry = True
                modeEmbed = discord.Embed(title="Setup", color=0x3083e3, description="Do you want manga updates sent to your DMs or a server?")
                sentEmbedMode = await ctx.send(embed=modeEmbed, view=modeView)
                await modeView.wait()
                if modeView.value is None:
                    await sentEmbedMode.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                    return
                else:
                    mode = modeView.value
            if modeEntry == True:
                await sentEmbedMode.delete()
            if mode == "user":
                userid = ctx.message.author.id
                if mongodb.checkUserExist(userid) == True:
                    alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup. Run the command `+addmanga` to add manga.")
                    await ctx.send(embed=alrfinishUS)
                elif mongodb.checkUserExist(userid) == False:
                    userInfo = await self.bot.fetch_user(userid)
                    username = f"{userInfo.name}#{userInfo.discriminator}"
                    mongodb.addUser(username, userid)
                    embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command `+addmanga` to add manga.")
                    await ctx.send(embed=embedUser)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
            elif mode == "server":
                if ctx.author.guild_permissions.administrator == False:
                    permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to setup this server's account. You need `Administrator` permission to use this.")
                    await ctx.send(embed=permissionError, delete_after=5.0)
                    return
                serverid = ctx.message.guild.id
                if mongodb.checkServerExist(serverid) == True:
                    alrfinishSS = discord.Embed(title="Setup", color=0x3083e3, description="This server is already setup. Run the command `+addmanga` to add manga.")
                    await ctx.send(embed=alrfinishSS, delete_after=10.0)
                elif mongodb.checkServerExist(serverid) == False:
                    embedServer = discord.Embed(title="Setup", color=0x3083e3, description="What channel would you like me to send manga updates to?")
                    sentEmbedServer = await ctx.send(embed=embedServer)
                    try:
                        channel = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                    except asyncio.TimeoutError:
                        await sentEmbedServer.delete()
                        await ctx.send(embed=timeoutError, delete_after=5.0)
                    else:
                        try:
                            channelid = channel.channel_mentions[0].id
                        except:
                            channelError = discord.Embed(title="Error", color=0xff4f4f, description="You did not input a valid channel. Please re-run the command to try again.")
                            await ctx.send(embed=channelError, delete_after=5.0)
                            return
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
        else:
            userid = ctx.message.author.id
            if mongodb.checkUserExist(userid) == True:
                alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup. Run the command `+addmanga` to add manga.")
                await ctx.send(embed=alrfinishUS)
            elif mongodb.checkUserExist(userid) == False:
                userInfo = await self.bot.fetch_user(userid)
                username = f"{userInfo.name}#{userInfo.discriminator}"
                mongodb.addUser(username, userid)
                embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up. Run the command `+addmanga` to add manga.")
                await ctx.send(embed=embedUser)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
    
    @commands.command(name="deleteaccount")
    async def deleteaccount(self, ctx, *, arg=None):
        confirmView = Confirm()
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            modeView = Mode()
            mode = arg
            modeEntry = False
            if (mode == None) or (mode != "server" and mode != "user"):
                modeEntry = True
                modeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Do you want to delete your user account or the server's account?")
                sentEmbedMode = await ctx.send(embed=modeEmbed, view=modeView)
                await modeView.wait()
                if modeView.value is None:
                    await sentEmbedMode.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                    return
                else:
                    mode = modeView.value
            if modeEntry == True:
                await sentEmbedMode.delete()
            if mode == "user":
                check = mongodb.checkUserExist(ctx.message.author.id)
            else:
                if ctx.author.guild_permissions.administrator == False:
                    permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to delete this server's account. You need `Administrator` permission to use this.")
                    await ctx.send(embed=permissionError, delete_after=5.0)
                    return
                check = mongodb.checkServerExist(ctx.message.guild.id)
            if check == True:
                confirmEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Are you sure you want to delete the account? (Your manga list will be gone forever!)")
                sentEmbedConfirm = await ctx.send(embed=confirmEmbed, view=confirmView)
                await confirmView.wait()
                if confirmView.value is None:
                    await sentEmbedConfirm.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                elif confirmView.value == True:
                    await sentEmbedConfirm.delete()
                    if mode == "user":
                        userid = ctx.message.author.id
                        mongodb.removeUser(userid)
                    elif mode == "server":
                        serverid = ctx.message.guild.id
                        mongodb.removeServer(serverid)
                    completeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Your account has been deleted.")
                    await ctx.send(embed=completeEmbed)
                else:
                    await sentEmbedConfirm.delete()
                    cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
                    await ctx.send(embed=cancelEmbed)
                    return
            elif check == False:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="You are not setup yet. Run the command `+setup` to setup.")
                await ctx.send(embed=completeError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
        else:
            check = mongodb.checkUserExist(ctx.message.author.id)
            if check == True:
                confirmEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Are you sure you want to delete the account? (The account's manga list will be gone forever!)")
                sentEmbedConfirm = await ctx.send(embed=confirmEmbed, view=confirmView)
                await confirmView.wait()
                if confirmView.value is None:
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                elif confirmView.value == True:
                    mongodb.removeUser(ctx.message.author.id)
                    completeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Your account has been deleted.")
                    await ctx.send(embed=completeEmbed)
                else:
                    return
            elif check == False:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="You are not setup yet. Run the command `+setup` to setup.")
                await ctx.send(embed=completeError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
            
    @commands.command(name="setchannel")
    async def setchannel(self, ctx):
        guild = ctx.message.guild
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            if ctx.author.guild_permissions.administrator == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to set this server's update channel. You need `Administrator` permission to use this.")
                await ctx.send(embed=permissionError, delete_after=5.0)
                return
            serverid = ctx.message.guild.id
            check = mongodb.checkServerExist(serverid)
            if check == True:
                channelEmbed = discord.Embed(title="Change Channel", color=0x3083e3, description="What channel would you like me to send manga updates to?")
                sentCEmbed = await ctx.send(embed=channelEmbed)
                try:
                    channel = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                except asyncio.TimeoutError:
                    await sentCEmbed.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    try:
                        channelid = channel.channel_mentions[0].id
                    except:
                        channelError = discord.Embed(title="Error", color=0xff4f4f, description="You did not input a valid channel. Please re-run the command to try again.")
                        await ctx.send(embed=channelError, delete_after=5.0)
                        return
                    if channelid == mongodb.findChannel(serverid):
                        sameError = discord.Embed(title="Error", color=0xff4f4f, description="This channel is already set as the channel for manga updates.")
                        await ctx.send(embed=sameError)
                    else:
                        mongodb.updateChannel(serverid, channelid)
                        finishChange = discord.Embed(title="Change Channel", color=0x3083e3, description=f"The server's channel has been sucessfully changed to `#{guild.get_channel(channelid)}`.")
                        await ctx.send(embed=finishChange)
            elif check == False:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="You are not setup yet. Run the command `+setup` to setup.")
                await ctx.send(embed=completeError, delete_after=5.0)
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                await ctx.send(embed=completeError, delete_after=5.0)
        else:
            dmError = discord.Embed(title="Error", color=0xff4f4f, description="This command cannot be run in DMs.")
            await ctx.send(embed=dmError)

    @commands.command(name="serverdbcp")
    async def serverdbcp(self, ctx):
        if ctx.message.author.id == config["ownerid"]:
            allServers = mongodb.getAllIds("server")
            currentServers = []
            description = ""
            for i in self.bot.guilds:
                currentServers.append(i.id)
            for server in allServers:
                if server not in currentServers:
                    description += f"• {server}\n"
            if description == "":
                description = "All servers in the DB still contain the bot."
            serversEmbed = discord.Embed(title=f"Non-Existing Servers in DB", color=0x3083e3, description=description)
            serversEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
            await ctx.send(embed=serversEmbed)
        else:
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You found a hidden command! Too bad only the bot owner can use this.")
            await ctx.send(embed=permissionError, delete_after=5.0)

def setup(bot):
    bot.add_cog(General(bot))