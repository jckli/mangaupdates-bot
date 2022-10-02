import discord
from discord.ext import commands
from discord.commands import Option, slash_command, SlashCommandGroup
import validators
import os
from core.mongodb import Mongo
from core.mangaupdates import MangaUpdates


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

class Confirm(discord.ui.View):
    def __init__(self):
        super().__init__(timeout=15.0)
        self.value = None
        self.interaction = None

    @discord.ui.button(label="Confirm", style=discord.ButtonStyle.green)
    async def confirm(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = True
        self.interaction = interaction
        self.stop()

    @discord.ui.button(label="Cancel", style=discord.ButtonStyle.red)
    async def cancel(self, button: discord.ui.Button, interaction: discord.Interaction):
        self.value = False
        self.interaction = interaction
        self.stop()

# Helpers
async def is_admin(ctx):
    has_add_permission = ctx.author.guild_permissions.administrator
    if not has_add_permission:
        permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to delete this server's account. You need `Administrator` permission to use this.")
        await ctx.respond(embed=permissionError, view=None)
    return has_add_permission

async def is_server_exists(ctx):
    server_exists = await mongo.check_server_exist(ctx.guild.id)
    if not server_exists:
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        await ctx.respond(embed=setupError, view=None)
    return server_exists

async def validate_admin_and_server(ctx):
    return await is_admin(ctx) and await is_server_exists(ctx)

class MangaGeneral(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    server = SlashCommandGroup("server", description="Server commands", guild_only=True)

    @server.command(name="setup", description="Sets up your server for updates")
    async def setup(self, ctx, channel: Option(discord.TextChannel, required=True)):
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

    @server.command(name="setchannel", description="Sets the server's channel that updates are sent to")
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

    @server.command(name="delete", description="Deletes your server's account and list")
    async def delete(self, ctx):
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        serverExist = await mongo.check_server_exist(ctx.guild.id)
        if serverExist is False:
            await ctx.respond(embed=setupError, view=None)
            return
        if ctx.author.guild_permissions.administrator is False:
            permissionError = discord.Embed(title="Error", color=0xff4f4f, description="You don't have permission to delete this server's account. You need `Administrator` permission to use this.")
            await ctx.respond(embed=permissionError, view=None)
            return
        embed = discord.Embed(title="Delete Account", color=0x3083e3, description="Are you sure you want to delete your account? (Your manga list will be gone forever!)")
        confirm = Confirm()
        await ctx.respond(embed=embed, view=confirm)
        await confirm.wait()
        if confirm.value is None:
            await confirm.interaction.response.edit_message(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            await mongo.remove_server(ctx.guild.id)
            completeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Successfully deleted your account.")
            await confirm.interaction.response.edit_message(embed=completeEmbed, view=None)

    @server.command(name="addadminrole", description="Sets a role that can modify the manga list")
    async def add_admin_role(self, ctx, role: Option(discord.Role, required=True)):
        hasPermission = await validate_admin_and_server(ctx)
        if not hasPermission:
            return
        curRole = await mongo.get_admin_role_server(ctx.guild.id)
        if curRole == role.id:
            sameError = discord.Embed(title="Error", color=0xff4f4f, description="This role is already set as the admin role.")
            await ctx.respond(embed=sameError, view=None)
            return
        await mongo.add_admin_role_server(ctx.guild.id, role.id)
        embedChannel = discord.Embed(title="Add Admin Role", color=0x3083e3, description=f"Successfully allowed role `{role}` to modify to the manga list.")
        await ctx.respond(embed=embedChannel, view=None)

    @server.command(name="removeadminrole", description="Removes the role that can modify the manga list")
    async def remove_admin_role(self, ctx):
        hasPermission = await validate_admin_and_server(ctx)
        if not hasPermission:
            return
        embed = discord.Embed(title="Remove Admin Role", color=0x3083e3, description="Are you sure you want to remove the admin role?")
        confirm = Confirm()
        await ctx.respond(embed=embed, view=confirm)
        await confirm.wait()
        if confirm.value is None:
            await confirm.interaction.response.edit_message(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            await mongo.remove_admin_role_server(ctx.guild.id)
            embedChannel = discord.Embed(title="Remove Admin Role", color=0x3083e3, description=f"Successfully removed the admin role.")
            await confirm.interaction.response.edit_message(embed=embedChannel, view=None)

    @server.command(name="test", description="Tests sending updates to your server")
    async def testsending(self, ctx):
        channel = await mongo.get_channel(ctx.guild.id)
        if channel is None:
            noChannelError = discord.Embed(title="Error", color=0xff4f4f, description="You have no channel set up. Please set one up with the `setup` command.")
            await ctx.respond(embed=noChannelError, view=None)
            return
        else:
            hasPermission = False
            adminRole = await mongo.get_admin_role_server(ctx.guild.id)
            authorRoles = [r.id for r in ctx.author.roles]
            hasPermission = ctx.author.guild_permissions.administrator or (adminRole in authorRoles)
            if not hasPermission:
                permissionError = discord.Embed(title="Error", color=0xff4f4f, description=("You don't have permission to test alerts. Set a role to modify manga with `/server addadminrole` or have `Administrator` permission."))
                await ctx.respond(embed=permissionError, view=None)
                return
            else:
                try:
                    channelObject = self.bot.get_channel(channel)
                    if channelObject is None:
                        noChannelError = discord.Embed(title="Error", color=0xff4f4f, description="Channel doesn't exist anymore, please reset the channel with the `server setchannel` command.")
                        await ctx.respond(embed=noChannelError, view=None)
                        return
                    mangaTestEmbed = discord.Embed(title=f"Testing Update", url="https://hayasaka.moe/", description=f"This is a test alert for the MangaUpdates Bot.", color=0x3083e3)
                    if self.bot.user and self.bot.user.avatar:
                        mangaTestEmbed.set_author(name="MangaUpdates", icon_url=self.bot.user.avatar.url)
                    mangaTestEmbed.add_field(name="Chapter", value="c.1", inline=True)
                    mangaTestEmbed.add_field(name="Group", value="The MangaUpdates Bot Team", inline=True)
                    mangaTestEmbed.add_field(name="Scanlator Link", value="https://hayasaka.moe/", inline=False)
                    await channelObject.send(embed=mangaTestEmbed, view=None)
                    mainResponseEmbed = discord.Embed(title="Test Manga Update Message", color=0x3083e3, description="Congrats! Manga can be sent to your server.")
                    await ctx.respond(embed=mainResponseEmbed, view=None)
                except:
                    mainResponseEmbed = discord.Embed(title="Test Manga Update Message", color=0xff4f4f, description="Oops! Something went wrong. Give the bot permissions to send messages in the channel you have set and then try again.")
                    await ctx.respond(embed=mainResponseEmbed, view=None)
                    return

    user = SlashCommandGroup(name="user", description="User commands")

    @user.command(name="setup", description="Sets up your user for updates")
    async def setup(self, ctx):
        userExist = await mongo.check_user_exist(ctx.author.id)
        if userExist is True:
            alrfinishUS = discord.Embed(title="Setup", color=0x3083e3, description="You are already setup.")
            await ctx.respond(embed=alrfinishUS, view=None)
            return
        username = f"{ctx.author.name}#{ctx.author.discriminator}"
        await mongo.add_user(username, ctx.author.id)
        embedUser = discord.Embed(title="Setup", color=0x3083e3, description="Great! You're all set up and can add manga now.")
        await ctx.respond(embed=embedUser, view=None)
    
    @user.command(name="delete", description="Deletes your user account and list")
    async def delete(self, ctx):
        setupError = discord.Embed(title="Error", color=0xff4f4f, description="Sorry! Please run the setup command first.")
        userExist = await mongo.check_user_exist(ctx.author.id)
        if userExist is False:
            await ctx.respond(embed=setupError, view=None)
            return
        embed = discord.Embed(title="Delete Account", color=0x3083e3, description="Are you sure you want to delete your account? (Your manga list will be gone forever!)")
        confirm = Confirm()
        await ctx.respond(embed=embed, view=confirm)
        await confirm.wait()
        if confirm.value is None:
            await confirm.interaction.response.edit_message(embed=timeoutError, view=None)
            return
        elif confirm.value is False:
            cancelEmbed = discord.Embed(title=f"Canceled", color=0x3083e3, description="Successfully canceled.")
            await confirm.interaction.response.edit_message(embed=cancelEmbed, view=None)
            return
        else:
            await mongo.remove_user(ctx.author.id)
            completeEmbed = discord.Embed(title="Delete Account", color=0x3083e3, description="Successfully deleted your account.")
            await confirm.interaction.response.edit_message(embed=completeEmbed, view=None)

def setup(bot):
    bot.add_cog(MangaGeneral(bot))