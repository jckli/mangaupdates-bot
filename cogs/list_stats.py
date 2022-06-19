import discord
from discord.ext import commands, tasks
import requests
import os

class ListStatistics(commands.Cog):
    def __init__(self, bot):
        self.bot = bot
        self.update_stats.start()

    @tasks.loop(hours=1)
    async def update_stats(self):
        servers = self.bot.guilds
        botUsers = 0
        for i in servers:
            botUsers += i.member_count
        serverCount = len(servers)
        # top.gg update
        topggToken = os.environ.get("TOPGG_TOKEN")
        topggurl = "https://top.gg/api/bots/880694914365685781/stats"
        header = {"Authorization": topggToken}

        body = {"server_count": f"{serverCount}"}
        try:
            requests.post(topggurl, data=body, headers=header)
            print("Successfully posted stats to top.gg.")
        except Exception as err:
            print(f"Failed to post to top.gg: {err}")
        
        # discordbotlist.com update
        dblToken = os.environ.get("DBL_TOKEN")
        dblurl = "https://discordbotlist.com/api/v1/bots/880694914365685781/stats"
        header = {"Authorization": dblToken}
        body = {"guilds": f"{serverCount}", "users": f"{botUsers}"}
        try:
            requests.post(dblurl, data=body, headers=header)
            print("Successfully posted stats to dbl.com.")
        except Exception as err:
            print(f"Failed to post to dbl.com: {err}")

    @update_stats.before_loop
    async def before_update_stats(self):
        await self.bot.wait_until_ready()

def setup(bot):
    bot.add_cog(ListStatistics(bot))