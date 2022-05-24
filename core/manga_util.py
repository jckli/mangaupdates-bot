import discord
import re
from core.utils import Util

util = Util()

timeoutError = discord.Embed(title="Error", description="You didn't respond in time!", color=0xff4f4f)

class ConfirmButtons(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None
        self.interaction = None

    @discord.ui.button(label=f"Confirm", style=discord.ButtonStyle.green)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = True
        self.interaction = interaction
        self.stop()

    @discord.ui.button(label=f"Cancel", style=discord.ButtonStyle.red)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = False
        self.interaction = interaction
        self.stop()

class SelectManga(discord.ui.Select):
    def __init__(self, manga_list, need_confirmation):
        manga_desc = []
        for manga in manga_list:
            manga_desc.append(discord.SelectOption(label=manga["name"], description=manga["description"]))
        super().__init__(
            placeholder="Choose a manga series...",
            min_values=1,
            max_values=1,
            options=manga_desc
        )
        self.confirm = need_confirmation
        self.manga_list = manga_list
        self.id = None
        self.value = None

    async def callback(self, interaction: discord.Interaction):
        if self.confirm is True:
            title = self.manga_list["info"]["title"]
            description = self.manga_list["info"]["description"]
            image = self.manga_list["info"]["image"]["url"]["original"]
            confirmEmbed = discord.Embed(title=f"Did you mean to add `{title}`?", description=description, color=0x3083e3)
            confirmEmbed.set_thumbnail(url=image)
            await interaction.response.edit_message(embed=confirmEmbed, view=ConfirmButtons())
            await ConfirmButtons().wait()
            if ConfirmButtons().value is None:
                await ConfirmButtons().interaction.send_message(embed=timeoutError)
            else:
                if ConfirmButtons().value is True:
                    cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
                    await ConfirmButtons().interaction.send_message(embed=cancelEmbed)
                else:
                    value = int(self.values[0][:2].replace(".", "")) - 1
                    mangaid = self.manga_list[value]["id"]
                    self.id = mangaid
                    self.i = ConfirmButtons().interaction
                    await interaction.response.edit_message()
        else:
            value = int(self.values[0][:2].replace(".", "")) - 1
            mangaid = self.manga_list[value]["id"]
            self.id = mangaid
            self.value = value
            await interaction.response.edit_message()

class SearchData:
    def __init__(self, series_info):
        self.title = series_info["title"]
        desc = series_info["description"]
        self.description = util.format_mu_description(desc)
        if series_info["completed"] is True:
            self.status = "Completed"
        else:
            self.status = "Ongoing"
        self.image = series_info["image"]["url"]["original"]
        self.url = series_info["url"]
        self.mangatype = series_info["type"]
        self.year = series_info["year"]
        self.latest_chapter = series_info["latest_chapter"]
        self.rating = series_info["bayesian_rating"]
        authorsList = []
        artistsList = []
        for author in series_info["authors"]:
            if author["type"] == "Author":
                authorsList.append(author["name"])
            if author["type"] == "Artist":
                artistsList.append(author["name"])
        self.authors = ", ".join(authorsList)
        self.artists = ", ".join(artistsList)