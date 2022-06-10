import discord
from discord.ext import commands
from dotenv import load_dotenv
import os
from core.mongodb import Mongo

load_dotenv()
mongo = Mongo()

bot = commands.Bot(
    intents=discord.Intents(guilds=True)
)
bot.remove_command("help")

for file in os.listdir("./cogs"):
    if file.endswith(".py"):
        name = file[:-3]
        bot.load_extension(f"cogs.{name}")

@bot.event
async def on_ready():
    print(f"Bot is online.")
    await bot.change_presence(activity=discord.Game(name="/help"))

@bot.event
async def on_guild_remove(guild):
    mongo.remove_server(guild.id)

try:
    bot.run(os.environ.get("TOKEN"))
except Exception as err:
    print(f"Error: {err}")