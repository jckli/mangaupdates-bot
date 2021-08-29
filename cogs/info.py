import discord
from discord.ext import commands

from core import mongodb

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.command(name="ping")
    async def ping(self, ctx):
        await ctx.send(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @commands.command(name="test")
    async def addserver(self, ctx):
        id = ctx.message.guild.id
        name = self.bot.fetch_guild(id)
        
        mongodb.add_server(id, name)

def setup(bot):
    bot.add_cog(Information(bot))