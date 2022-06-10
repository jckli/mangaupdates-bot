import discord
from discord.ext import commands
from discord.commands import Option, slash_command, SlashCommandGroup
import validators
import os
from core.mongodb import Mongo
from core.mangaupdates import MangaUpdates
from core.manga_util import SearchData

mongo = Mongo()
mangaupdates = MangaUpdates()
ghuser = os.environ.get("GITHUB_USER")

timeoutError = discord.Embed(title="Error", description="You didn't respond in time! Please rerun the command.", color=0xff4f4f)

class Mode(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None
        self.interaction = None

    @discord.ui.button(label=f'User (DMs)', style=discord.ButtonStyle.grey)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "user"
        self.interaction = interaction
        self.stop()

    @discord.ui.button(label=f'Server', style=discord.ButtonStyle.grey)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = "server"
        self.interaction = interaction
        self.stop()

# search command views
class SelectMangaView(discord.ui.View):
    def __init__(self, manga_list):
        super().__init__(timeout=15.0)
        self.select_manga = SelectManga(manga_list=manga_list)
        self.add_item(self.select_manga)

class SelectManga(discord.ui.Select):
    def __init__(self, manga_list):
        manga_desc = []
        for manga in manga_list:
            manga_desc.append(discord.SelectOption(label=manga["name"], description=manga["description"]))
        super().__init__(
            placeholder="Choose a manga series...",
            min_values=1,
            max_values=1,
            options=manga_desc
        )
        self.manga_list = manga_list
        self.finish = None

    async def callback(self, interaction: discord.Interaction):
        value = int(self.values[0][:2].replace(".", "")) - 1
        mangaid = self.manga_list[value]["id"]
        series_info = await mangaupdates.series_info(mangaid)
        data = SearchData(series_info)
        result = discord.Embed(title=f"{data.title} ({data.status})", url=data.url, color=0x3083e3, description=data.description)
        if data.image != None:
            result.set_image(url=data.image)
        result.set_author(name="MangaUpdates", icon_url=interaction.client.user.avatar.url)
        result.set_footer(text=f"Requested by {interaction.user.name}", icon_url=interaction.user.display_avatar)
        result.add_field(name="Year", value=data.year, inline=True)
        result.add_field(name="Type", value=data.mangatype, inline=True)
        result.add_field(name="Latest Chapter", value=data.latest_chapter, inline=True)
        result.add_field(name="Author(s)", value=data.authors, inline=True)
        result.add_field(name="Artist(s)", value=data.artists, inline=True)
        result.add_field(name="Rating", value=data.rating, inline=True)
        await interaction.response.edit_message(embed=result, view=None)
        self.finish = True

class MangaGeneral(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @slash_command(name="search", description="Searches for a manga series", guild_ids=[721216108668911636])
    async def search(self, ctx, manga: Option(str, description="The name of the manga series (can use mangaupdates links)", required=True)):
        if validators.url(manga) is True:
            link = manga.partition("https://www.mangaupdates.com/series/")[2]
            mangaid = link.partition("/")[0]
            mangaid = await mangaupdates.convert_new_id(mangaid)
            series_info = await mangaupdates.series_info(mangaid)
            data = SearchData(series_info)
            result = discord.Embed(title=f"{data.title} ({data.status})", url=data.url, color=0x3083e3, description=data.description)
            if data.image != None:
                result.set_image(url=data.image)
            result.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
            result.set_footer(text=f"Requested by {ctx.author.name}", icon_url=ctx.author.display_avatar)
            result.add_field(name="Year", value=data.year, inline=True)
            result.add_field(name="Type", value=data.mangatype, inline=True)
            result.add_field(name="Latest Chapter", value=data.latest_chapter, inline=True)
            result.add_field(name="Author(s)", value=data.authors, inline=True)
            result.add_field(name="Artist(s)", value=data.artists, inline=True)
            result.add_field(name="Rating", value=data.rating, inline=True)
            await ctx.respond(embed=result)
        elif validators.url(manga) is not True:
            await ctx.defer()
            search_results = []
            description = "Type the number of the manga you want to see information for.\n"
            search_data = await mangaupdates.search_series(manga)
            if search_data["results"] == []:
                resultError = discord.Embed(title="Error", color=0xff4f4f, description="No mangas were found.")
                await ctx.respond(embed=resultError)
                return
            elif search_data["results"] != []:
                count = 1
                for manga_data in search_data["results"]:
                    manga = manga_data["record"]
                    name = manga["title"]
                    year = manga["year"]
                    rating = manga["bayesian_rating"]
                    description += f"{count}. {name} ({year}, Rating: {rating})\n"
                    search_results.append({"id": manga["series_id"], "name": f"{count}. {name}", "description": f"{year}, Rating: {rating}", "info": manga})
                    count += 1
                searchEmbed = discord.Embed(title="Search Results", color=0x3083e3, description=description)
                manga_drop = SelectMangaView(manga_list=search_results)
                search = await ctx.respond(embed=searchEmbed, view=manga_drop)
                await manga_drop.wait()
                if manga_drop.select_manga.finish is None:
                    await search.edit(embed=timeoutError, view=None)
                    return
                else:
                    return
            else:
                completeError = discord.Embed(title="Error", color=0xff4f4f, description=f"Something went wrong. Create an issue here for support: https://github.com/{ghuser}/mangaupdates-bot")
                await ctx.respond(embed=completeError)
                return
        else:
            completeError = discord.Embed(title="Error", color=0xff4f4f, description=f"Something went wrong. Create an issue here for support: https://github.com/{ghuser}/mangaupdates-bot")
            await ctx.respond(embed=completeError)
            return

    setup = SlashCommandGroup("setup", "Setup commands")

    @setup.command(name="server", description="Sets up your server for manga updates", guild_only=True, guild_ids=[721216108668911636])
    async def server(self, ctx, channel: Option(discord.TextChannel, required=True)):
        serverExist = await mongo.check_server_exist(ctx.guild.id)
        if serverExist is True:
            alrfinishSS = discord.Embed(title="Setup", color=0x3083e3, description="This server is already setup.")
            await ctx.respond(embed=alrfinishSS, view=None)
            return
        if ctx.author.guild_permissions.administrator is False:
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to setup this server's account. You need `Administrator` permission to use this.")
            await ctx.respond(embed=permissionError, view=None)
            return
        channelid = channel.id
        await mongo.add_server(ctx.guild.name, ctx.guild.id, channelid)
        embedServerF = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up and can add manga now.")
        await ctx.respond(embed=embedServerF, view=None)

    @setup.command(name="user", description="Sets up your user for manga updates", guild_ids=[721216108668911636])
    async def user(self, ctx):
        userExist = await mongo.check_user_exist(ctx.author.id)
        if userExist is True:
            alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup.")
            await ctx.respond(embed=alrfinishUS, view=None)
            return
        username = f"{ctx.author.name}#{ctx.author.discriminator}"
        await mongo.add_user(username, ctx.author.id)
        embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up and can add manga now.")
        await ctx.respond(embed=embedUser, view=None)

    @slash_command(name="setchannel", description="Sets the server's that manga chapter updates are sent to", guild_only=True, guild_ids=[721216108668911636])
    async def setchannel(self, ctx, channel: Option(discord.TextChannel, required=True)):
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        serverExist = await mongo.check_server_exist(ctx.guild.id)
        if serverExist is False:
            await ctx.respond(embed=setupError, view=None)
            return
        if ctx.author.guild_permissions.administrator is False:
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to change this server's channel for manga updates. You need `Administrator` permission to use this.")
            await ctx.respond(embed=permissionError, view=None)
            return
        channelid = channel.id
        curchannel = await mongo.get_channel(ctx.guild.id)
        if channelid == curchannel:
            sameError = discord.Embed(title="Error", color=0xff4f4f, description="This channel is already set as the channel for manga updates.")
            await ctx.respond(embed=sameError, view=None)
            return
        await mongo.set_channel(ctx.guild.id, channelid)
        embedChannel = discord.Embed(title="Set Channel", color=0x3083e3, description=f"The server's channel has been successfully changed to `#{ctx.guild.get_channel(channelid)}`.")
        await ctx.respond(embed=embedChannel, view=None)

def setup(bot):
    bot.add_cog(MangaGeneral(bot))