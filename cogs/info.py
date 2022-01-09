import discord
from discord.ext import commands

import time
import datetime
from datetime import datetime
from datetime import timedelta
import json

startTime = time.time()

# Load config
with open("config.json", "r") as f:
    config = json.load(f)

class InviteLink(discord.ui.View):
    def __init__(self, link):
        super().__init__()
        self.add_item(discord.ui.Button(label="Invite Link", url=link))

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="help")
    async def help(self, ctx):
        embed = discord.Embed(title="MangaUpdates Help", color=0x3083e3)
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="+help", value="Displays this message.", inline=False)
        embed.add_field(name="+alert", value="Displays bot alerts/announcements.", inline=False)
        embed.add_field(name="+ping", value="Pong! Displays the ping.", inline=False)
        embed.add_field(name="+invite", value="Displays bot invite link.", inline=False)
        embed.add_field(name="+source", value="Displays bot's GitHub repository.", inline=False)
        embed.add_field(name="+setup `user/server`", value="Setup your user/server for manga updates.", inline=False)
        embed.add_field(name="+addmanga `user/server`", value="Adds manga to your list to be tracked.", inline=False)
        embed.add_field(name="+removemanga `user/server`", value="Removes manga from your list that were tracked.", inline=False)
        embed.add_field(name="+mangalist `user/server`", value="Lists all manga that are being tracked.", inline=False)
        embed.add_field(name="+clearmanga `user/server`", value="Removes all manga from your current manga list.", inline=False)
        embed.add_field(name="+setchannel", value="Changes the server's channel that manga chapter updates are sent to.", inline=False)
        embed.add_field(name="+deleteaccount `user/server`", value="Deletes your account and your manga list.", inline=False)
        embed.add_field(name="+setgroup `user/server`", value="Sets a manga's scan group. Only that scan group's chapter updates for that manga will be sent.", inline=False)
        embed.add_field(name="+search `manga`", value="Searches for information about a manga.", inline=False)
        embed.set_footer(text="PLEASE READ (server owners): +alert")
        await ctx.send(embed=embed)

    @commands.command(name="ping")
    async def ping(self, ctx):
        await ctx.send(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @commands.command(name="invite")
    async def invite(self, ctx):
        embed = discord.Embed(title="Invite Link", color=0x3083e3, description="Invite me to your own servers!")
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="Link", value="https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=applications.commands%20bot&permissions=268856384")
        await ctx.send(embed=embed)

    @commands.command(description="Shows the bot uptime.")
    async def botinfo(self, ctx):
        if ctx.message.author.id == config["ownerid"]:
            # Get all users in all servers the bot is in.
            activeServers = self.bot.guilds
            botUsers = 0
            for i in activeServers:
                botUsers += i.member_count
            # Get the current uptime.
            currentTime = time.time()
            differenceUptime = int(round(currentTime - startTime))
            uptime = str(timedelta(seconds = differenceUptime))
            # Make the embed for the message.
            botinfo = discord.Embed(
                title="Bot info",
                color=0x3083e3,
                timestamp=datetime.now(),
                description=f"**Server Count:** {len(self.bot.guilds)}\n**Bot Users:** {botUsers}\n**Bot Uptime:** {uptime}"
            )
            botinfo.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
            await ctx.send(embed=botinfo)
        else:
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You found a hidden command! Too bad only the bot owner can use this.")
            await ctx.send(embed=permissionError, delete_after=5.0)

    @commands.command(name="source")
    async def source(self, ctx):
        embed = discord.Embed(title="Source Code", color=0x3083e3, description="MangaUpdate's source code can be found on GitHub. Any issues with the bot can be raised there.")
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="Link", value="https://github.com/ohashizu/mangaupdates-bot")
        embed.set_footer(text="PLEASE READ (server owners): +alert")
        await ctx.send(embed=embed)

    @commands.command(name="alert")
    async def alert(self, ctx):
        link = 'https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=applications.commands%20bot&permissions=268856384'
        description = f"*This alert is only if you invited the bot before December 25th.*\n\nYo everyone! Recently Discord changed their API to require message content as intents. They want every bot to move to slash commands. This means that this bot needs new permissions to use these slash commands (don't ask me why).\n\nPlease reinvite the bot with the link or else by April 30, 2022, you won't be able to use the bot. Thanks for understanding and using MangaUpdates Bot!"
        embed = discord.Embed(title="Alert - Please read", color=0x3083e3, description=description)
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.set_footer(text="If button doesn't work: https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=applications.commands%20bot&permissions=268856384")
        await ctx.send(embed=embed, view=InviteLink(link))

def setup(bot):
    bot.add_cog(Information(bot))