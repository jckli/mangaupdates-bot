import discord
from discord.ext import commands
from discord.commands import Option, slash_command
import validators
import os
from core.mangaupdates import MangaUpdates

mangaupdates = MangaUpdates()
ghuser = os.environ.get("GITHUB_USER")

timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)

# todo: make dropdown only work once, also make the timeout stop showing if it didnt time out (think it already does but dunno)
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
        title = series_info["title"]
        description = series_info["description"]
        if series_info["completed"] is True:
            status = "Completed"
        else:
            status = "Ongoing"
        image = series_info["image"]["url"]["original"]
        url = series_info["url"]
        mangatype = series_info["type"]
        year = series_info["year"]
        latest_chapter = series_info["latest_chapter"]
        rating = series_info["bayesian_rating"]
        authorsList = []
        artistsList = []
        for author in series_info["authors"]:
            if author["type"] == "Author":
                authorsList.append(author["name"])
            if author["type"] == "Artist":
                artistsList.append(author["name"])
        authors = ", ".join(authorsList)
        artists = ", ".join(artistsList)

        result = discord.Embed(title=title, url=url, color=0x3083e3, description=description)
        if image != None:
            result.set_image(url=image)
        result.set_author(name="MangaUpdates", icon_url=interaction.client.user.avatar.url)
        result.set_footer(text=f"Requested by {interaction.user.name}", icon_url=interaction.user.display_avatar)
        result.add_field(name="Year", value=year, inline=True)
        result.add_field(name="Type", value=mangatype, inline=True)
        result.add_field(name="Latest Chapter", value=latest_chapter, inline=True)
        result.add_field(name="Author(s)", value=authors, inline=True)
        result.add_field(name="Artist(s)", value=artists, inline=True)
        result.add_field(name="Rating", value=rating, inline=True)
        await interaction.response.send_message(embed=result)
        self.finish = True

class MangaGeneral(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @slash_command(name="search", description="Searches for a manga series", guild_ids=[721216108668911636])
    async def search(self, ctx, manga):
        if validators.url(manga) is True:
            link = manga.partition("https://www.mangaupdates.com/series/")[2]
            mangaid = link.partition("/")[0]
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
                await ctx.respond(embed=searchEmbed, view=manga_drop)
                await manga_drop.wait()
                manga_drop.disable_all_items()
                if manga_drop.select_manga.finish is None:
                    await ctx.respond(embed=timeoutError)
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

def setup(bot):
    bot.add_cog(MangaGeneral(bot))