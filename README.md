<p align="center">
  <a href="https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=bot&permissions=268856384">
    <img src="https://github.com/ohashizu/mangaupdates-bot/blob/master/icon.png" alt="Logo" width="80" height="80">
  </a>
  <h3 align="center">MangaUpdates Bot</h3>
  <p align="center">A Discord bot that can be used to keep track of your favorite mangas.</p>
</p>

## About

MangaUpdates is a simple but powerful bot that sends every new manga, manhwa, or doujin chapter update to either your direct messages or a server channel. You simply search for your manga using mangaupdates.com's search and select your favorite manga for MangaUpdates to track!

This bot utilizes a MongoDB database to store all manga lists as well as mangaupdates.com to query mangas and their RSS feed to track when new updates of mangas are released.

> **Note**: This bot is still in development. There may still be some bugs, as well as many features that are coming in the future.

## Invite

Invite the bot [here](https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=bot&permissions=268856384).

## Features/Why this bot?

- Send manga updates to you through discord DMs or in a server channel.
- Import MyAnimeList manga list to be tracked (In development).
- Specific scan group selection
- Supports most mangas, manhwas, doujins (utilizing mangaupdates.com)
- Easy setup
- Search information for your favorite manga
- Consistently updated, with many more features planned

## Why?

Personally, I read quite a bit of manga, manhwa, and doujins. However, MyAnimeList doesn't want to track any less popular mangas as well as basically all manhwas and doujins.

I wanted to create something that would track all my mangas/manhwas/doujins on a platform that I personally use often on both desktop and mobile. Thus, I created this bot as Discord is cross-platform and I use it quite a lot to chat and has a very easy to use interface.

## Issues

If you have any issues, please don't be afraid to raise an issue on [GitHub](https://github.com/ohashizu/mangaupdates-bot) or join the [Support Server](https://discord.gg/UcYspqftTF).

## Commands

You can see all commands within the bot with the `+help` command.

+help - Returns all commands.

+ping - Pong! Displays the ping.

+invite - Displays bot invite link.

+source - Displays bot's GitHub repository.

+setup `user`/`server` - Setup your user/server for manga updates.

+addmanga `user`/`server` - Adds manga to your list to be tracked.

+removemanga `user`/`server` - Removes manga from your list that were tracked.

+mangalist `user`/`server` - Lists all manga that are being tracked.

+clearmanga `user`/`server` - Removes all manga from your current manga list.

+setchannel - Changes the server's channel that manga chapter updates are sent to.

+deleteaccount `user`/`server` - Deletes your account and your manga list.

+setgroup `user`/`server` - Sets a manga's scan group. Only that scan group's chapter updates for that manga will be sent.

+search `manga` - Searches for information about a manga.