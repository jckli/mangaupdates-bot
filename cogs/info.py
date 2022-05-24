import discord
from discord.ext import commands
from discord.commands import Option, slash_command

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
        embed.add_field(name="setup `user/server`", value="Setup your user/server for manga updates.", inline=False)
        embed.add_field(name="addmanga `user/server`", value="Adds manga to your list to be tracked.", inline=False)
        embed.add_field(name="removemanga `user/server`", value="Removes manga from your list that were tracked.", inline=False)
        embed.add_field(name="mangalist `user/server`", value="Lists all manga that are being tracked.", inline=False)
        embed.add_field(name="clearmanga `user/server`", value="Removes all manga from your current manga list.", inline=False)
        embed.add_field(name="setchannel", value="Changes the server's channel that manga chapter updates are sent to.", inline=False)
        embed.add_field(name="deleteaccount `user/server`", value="Deletes your account and your manga list.", inline=False)
        embed.add_field(name="setgroup `user/server`", value="Sets a manga's scan group. Only that scan group's chapter updates for that manga will be sent.", inline=False)
        embed.add_field(name="search `manga`", value="Searches for information about a manga.", inline=False)
        await ctx.respond(embed=embed, ephemeral=True)

def setup(bot):
    bot.add_cog(Information(bot))