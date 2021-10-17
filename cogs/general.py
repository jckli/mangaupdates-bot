import discord
from discord.ext import commands

import asyncio

from core import mongodb, update

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
    async def setup(self, ctx):
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
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
                            embedServer = discord.Embed(title="Setup", color=0x3083e3, description="What channel would you like me to send manga updates to?")
                            sentEmbedServer = await ctx.send(embed=embedServer)
                            try:
                                channel = await self.bot.wait_for('message', check=lambda x: x.author.id == ctx.author.id, timeout=15)
                            except asyncio.TimeoutError:
                                await sentEmbedServer.delete()
                                await ctx.send(embed=timeoutError, delete_after=5.0)
                            else:
                                await sentEmbedServer.delete()
                                await channel.delete()
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
    async def deleteaccount(self, ctx):
        confirmView = Confirm()
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
                modeView = Mode()
                deleteEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Do you want to delete your user account or the server's account?")
                sentEmbed = await ctx.send(embed=deleteEmbed, view=modeView)
                await modeView.wait()
                if modeView.value is None:
                    await sentEmbed.delete()
                    await ctx.send(embed=timeoutError, delete_after=5.0)
                else:
                    await sentEmbed.delete()
                    if modeView.value == "user":
                        check = mongodb.checkUserExist(ctx.message.author.id)
                    else:
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
                            if modeView.value == "user":
                                userid = ctx.message.author.id
                                mongodb.removeUser(userid)
                            elif modeView.value == "server":
                                serverid = ctx.message.guild.id
                                mongodb.removeServer(serverid)
                            completeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Your account has been deleted.")
                            await ctx.send(embed=completeEmbed, delete_after=10.0)
                        else:
                            await sentEmbedConfirm.delete()
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
        timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)
        if isinstance(ctx.channel, discord.DMChannel) == False:
            guild = ctx.message.guild
            if guild.me.guild_permissions.manage_messages == False:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description="I don't have permission to run this command. Please grant me the: `Manage Messages` permission.")
                await ctx.send(embed=permissionError, delete_after=10.0)
                return
            else:
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
                        await sentCEmbed.delete()
                        await channel.delete()
                        try:
                            channelid = channel.channel_mentions[0].id
                        except:
                            channelError = discord.Embed(title="Error", color=0xff4f4f, description="You did not input a valid channel. Please re-run the command to try again.")
                            await ctx.send(embed=channelError, delete_after=5.0)
                            return
                        if channelid == mongodb.findChannel(serverid):
                            sameError = discord.Embed(title="Error", color=0xff4f4f, description="This channel is already set as the channel for manga updates.")
                            await ctx.send(embed=sameError, delete_after=5.0)
                        else:
                            mongodb.updateChannel(serverid, channelid)
                            finishChange = discord.Embed(title="Change Channel", color=0x3083e3, description=f"The server's channel has been sucessfully changed to `#{guild.get_channel(channelid)}`.")
                            await ctx.send(embed=finishChange, delete_after=10.0)
                elif check == False:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="You are not setup yet. Run the command `+setup` to setup.")
                    await ctx.send(embed=completeError, delete_after=5.0)
                else:
                    completeError = discord.Embed(title="Error", color=0xff4f4f, description="Something went wrong. Create an issue here for support: https://github.com/ohashizu/mangaupdates-bot")
                    await ctx.send(embed=completeError, delete_after=5.0)
        else:
            dmError = discord.Embed(title="Error", color=0xff4f4f, description="This command cannot be run in DMs.")
            await ctx.send(embed=dmError)

def setup(bot):
    bot.add_cog(General(bot))