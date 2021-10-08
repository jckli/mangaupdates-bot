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

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="help")
    async def help(self, ctx):
        embed = discord.Embed(title="MangaUpdates Help", color=0x3083e3)
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="+help", value="Displays this message", inline=False)
        embed.add_field(name="+ping", value="Pong! Displays the ping.", inline=False)
        embed.add_field(name="+invite", value="Displays bot invite link", inline=False)
        embed.add_field(name="+source", value="Displays bot's GitHub repository", inline=False)
        embed.add_field(name="+setup", value="Setup your user/server for manga updates.", inline=False)
        embed.add_field(name="+addmanga", value="Adds manga to your list to be tracked. (Optional: `user` or `server` after command for easy usage)", inline=False)
        embed.add_field(name="+removemanga", value="Removes manga from your list that were tracked. (Optional: `user` or `server` after command for easy usage)", inline=False)
        embed.add_field(name="+mangalist", value="Lists all manga that are being tracked. (Optional: `user` or `server` after command for easy usage)", inline=False)
        embed.add_field(name="+clearmanga", value="Removes all manga from your current manga list. (Optional: `user` or `server` after command for easy usage)", inline=False)
        await ctx.send(embed=embed)
        

    @commands.command(name="ping")
    async def ping(self, ctx):
        await ctx.send(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @commands.command(name="invite")
    async def invite(self, ctx):
        embed = discord.Embed(title="Invite Link", color=0x3083e3, description="Invite me to your own servers!")
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="Link", value="https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=bot&permissions=268856384")
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
        await ctx.send(embed=embed)

def setup(bot):
    bot.add_cog(Information(bot))