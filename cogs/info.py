import discord
from discord.ext import commands
from discord.commands import slash_command
import time
import os
from datetime import datetime
from datetime import timedelta

startTime = time.time()
ghuser = os.environ.get("GITHUB_USER")

class Link(discord.ui.View):
    def __init__(self, label, link):
        super().__init__()
        self.add_item(discord.ui.Button(label=label, url=link))

class InfoButtons(discord.ui.View):
    def __init__(self, support_server, github):
        super().__init__()
        self.add_item(discord.ui.Button(label="Support Server", url=support_server))
        self.add_item(discord.ui.Button(label="GitHub", url=github))

class Information(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @slash_command(name="help", description="Displays all commands", guild_ids=[721216108668911636])
    async def help(self, ctx):
        embed = discord.Embed(title="MangaUpdates Commands", color=0x3083e3,
            description="""
                **mangaupdates**: Displays basic information about MangaUpdates.
                **help**: Shows this message.
                **ping**: Pong! Displays the ping.
                **invite**: Displays bot invite link.
                **alert**: Displays bot announcements.
            """)
        embed.add_field(name="__Manga__", 
            value="""
                **setup**: Sets up your server/user for manga updates.
                **search `manga`**: Searches for information about a manga series.
                **manga list**: Displays your list of tracked manga.
                **manga add `manga`**: Adds a manga to your list to be tracked.
                **manga remove `manga`**: Removes a manga from your list.
            """, inline=False)
        await ctx.respond(embed=embed, ephemeral=True)

    @slash_command(name="mangaupdates", description="Displays basic information about MangaUpdates", guild_ids=[721216108668911636])
    async def mangaupdates(self, ctx):
        activeServers = self.bot.guilds
        botUsers = 0
        for i in activeServers:
            botUsers += i.member_count
        currentTime = time.time()
        differenceUptime = int(round(currentTime - startTime))
        uptime = str(timedelta(seconds = differenceUptime))
        botinfo = discord.Embed(
            title="MangaUpdates",
            color=0x3083e3,
            timestamp=datetime.now(),
            description=f"Thanks for using MangaUpdates bot! Any questions can be brought up in the support server. This bot is also open-source! All code can be found on GitHub (Please leave a star ⭐ if you enjoy the bot).\n\n**Server Count:** {len(self.bot.guilds)}\n**Bot Users:** {botUsers}\n**Bot Uptime:** {uptime}"
        )
        botinfo.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        await ctx.respond(embed=botinfo, view=InfoButtons("https://discord.gg/UcYspqftTF", f"https://github.com/{ghuser}/mangaupdates-bot"))

    @slash_command(name="ping", description="Pong!", guild_ids=[721216108668911636])
    async def ping(self, ctx):
        await ctx.respond(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @slash_command(name="alert", description="Displays bot alerts/announcements.", guild_ids=[721216108668911636])
    async def alert(self, ctx):
        link = f"https://github.com/{ghuser}/mangaupdates-bot"
        description = """
        Ayo! Thanks for keeping MangaUpdates Bot. I have been working on this version for a while now, and I hope you enjoy it.
        
        I have changed the manga updates whole system to use mangaupdates.com new API, as well as changed the commands system to use Discord's new slash commands.
        
        Anyways, sorry for any inconveniences when the bot wasn't working. Cheers!
        """
        embed = discord.Embed(title="Alert - Bot revamp", color=0x3083e3, description=description)
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        await ctx.respond(embed=embed, view=Link("GitHub", link))
    
    @slash_command(name="invite", description="Invite MangaUpdates to your own server", guild_ids=[721216108668911636])
    async def invite(self, ctx):
        link = "https://jackli.dev/mangaupdates"
        embed = discord.Embed(title="Invite Link", color=0x3083e3, description="Invite me to your own servers!")
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        await ctx.respond(embed=embed, view=Link("Invite", link))

def setup(bot):
    bot.add_cog(Information(bot))