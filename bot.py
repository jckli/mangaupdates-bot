import discord
from discord.ext import commands

import json
import os

from core import mongodb

# Load config
with open("config.json", "r") as f:
    config = json.load(f)

# Initialize
bot = commands.Bot(
    command_prefix = config["prefix"]
    )
bot.remove_command("help")

@bot.event
async def on_ready():
    print(f"Bot is online")
    await bot.change_presence(activity=discord.Game(name="+help"))
    for file in os.listdir("./cogs"):
        if file.endswith(".py"):
            name = file[:-3]
            bot.load_extension(f"cogs.{name}")

'''
@bot.event
async def on_guild_remove(guild):
    mongodb.removeServer(guild.id)
'''

try:
    bot.run(config["token"])
except Exception as err:
    print(f"Error: {err}")