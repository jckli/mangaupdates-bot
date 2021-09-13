import discord
from discord.ext import commands

import time
import datetime
from datetime import datetime
from datetime import timedelta

from core import mongodb

startTime = time.time()

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="ping")
    async def ping(self, ctx):
        await ctx.send(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @commands.command(description="Shows the bot uptime.")
    async def botinfo(self, ctx):
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
        botinfo.set_footer(
            text=f'Requested by {ctx.message.author.name}',
            icon_url=ctx.author.avatar_url
        )
        await ctx.send(embed=botinfo)

def setup(bot):
    bot.add_cog(Information(bot))