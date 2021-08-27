import discord
from discord.ext import commands

import datetime
import threading

class Commands(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="ping")
    async def ping(self, ctx):
        await ctx.send(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    

def setup(bot):
    bot.add_cog(Commands(bot))