import discord
from discord.ext import commands
from discord.commands import Option, slash_command
import time
from datetime import datetime
from datetime import timedelta

startTime = time.time()

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
        embed = discord.Embed(title="MangaUpdates Help", color=0x3083e3)
        embed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
        embed.add_field(name="help", value="Displays this message.", inline=False)
        embed.add_field(name="alert", value="Displays bot alerts/announcements.", inline=False)
        embed.add_field(name="ping", value="Pong! Displays the ping.", inline=False)
        embed.add_field(name="invite", value="Displays bot invite link.", inline=False)
        embed.add_field(name="source", value="Displays bot's GitHub repository.", inline=False)
        embed.add_field(name="setup", value="Sets up your user/server for manga updates.", inline=False)
        embed.add_field(name="addmanga", value="Adds manga to your list to be tracked.", inline=False)
        embed.add_field(name="removemanga `user/server`", value="Removes manga from your list that were tracked.", inline=False)
        embed.add_field(name="mangalist `user/server`", value="Lists all manga that are being tracked.", inline=False)
        embed.add_field(name="clearmanga `user/server`", value="Removes all manga from your current manga list.", inline=False)
        embed.add_field(name="setchannel", value="Changes the server's channel that manga chapter updates are sent to.", inline=False)
        embed.add_field(name="deleteaccount `user/server`", value="Deletes your account and your manga list.", inline=False)
        embed.add_field(name="setgroup `user/server`", value="Sets a manga's scan group. Only that scan group's chapter updates for that manga will be sent.", inline=False)
        embed.add_field(name="search `manga`", value="Searches for information about a manga.", inline=False)
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
        await ctx.respond(embed=botinfo, view=InfoButtons("https://discord.gg/Fr2BhuCkET", "https://github.com/jckli/mangaupdates-bot"))

    @slash_command(name="ping", description="Pong!", guild_ids=[721216108668911636])
    async def ping(self, ctx):
        await ctx.respond(f"🏓 Pong! My ping is {round(self.bot.latency * 1000)}ms")

    @slash_command(name="alert", description="Displays bot alerts/announcements.", guild_ids=[721216108668911636])
    async def alert(self, ctx):
        link = "https://github.com/jckli/mangaupdates-bot"
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