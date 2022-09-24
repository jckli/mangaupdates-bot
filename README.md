</br>
<p align="center">
  <a href="https://discord.com/oauth2/authorize?client_id=880694914365685781&scope=applications.commands%20bot&permissions=268856384" style:"margin-bottom: 0;">
    <img src="https://github.com/jckli/mangaupdates-bot/blob/master/icon.png" alt="Logo" width="100" height="100">
  </a>
  <h3 align="center">MangaUpdates Bot</h3>
  <p align="center">A Discord bot that can be used to keep track of your favorite mangas.</p>
</p>

## About

MangaUpdates is a simple but powerful bot that sends every new manga, manhwa, or doujin chapter update to either your direct messages or a server channel. You simply search for your manga using mangaupdates.com's search and select your favorite manga for MangaUpdates to track!

This bot utilizes a MongoDB database to store all manga lists as well as mangaupdates.com to query mangas and their RSS feed to track when new updates of mangas are released.

## Links

Invite the bot [here](https://jackli.dev/mangaupdates).

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

If you have any issues, please don't be afraid to raise an issue on [GitHub](https://github.com/jckli/mangaupdates-bot) or join the [Support Server](https://jackli.dev/discord).

## Commands

**mangaupdates**: Displays basic information about MangaUpdates.

**help**: Shows this message.

**ping**: Pong! Displays the ping.

**invite**: Displays bot invite link.

**alert**: Displays bot announcements.

**setup**: Sets up your server/user for manga updates.

**delete**: Deletes your account and your manga list.

**setchannel**: Sets the server's that manga chapter updates are sent to.

**search `manga`**: Searches for information about a manga series.

**manga list**: Displays your list of tracked manga.

**manga add `manga`**: Adds a manga to your list to be tracked.

**manga remove**: Removes a manga from your list.

**manga setgroup**: Sets a manga's scan group. Only that scan group's chapter updates for that manga will be sent.

## Self-Hosting

I would prefer you not to self-host, as it is unnecessarily complicated. I would much rather a feature request on my support discord or here. However if you still wish to do so, here is how to do so.

This bot is not really written for someone else to host for, so if some things don't work, I am not going to help much. You have to change the code yourself.

### Environment Variables
- `TOKEN`: Discord tot token
- `MONGO_USER`: MongoDB username
- `MONGO_PASS`: MongoDB password
- `MONGO_DB_NAME`: MongoDB database name
- `MU_USER`: MangaUpdates username
- `MU_PASS`: MangaUpdates password
- `GITHUB_USER`: GitHub username (for error responses)
- `TOPGG_TOKEN`: Top.gg token
- `DBL_TOKEN`: Discordbotlist.com token