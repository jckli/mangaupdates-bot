package mubot

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/mangaupdates-bot/bridge"
	"github.com/jckli/mangaupdates-bot/utils"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Client       *bot.Client
	ApiClient    *utils.Client
	BridgeServer *bridge.Server
	Logger       *slog.Logger
	InternalPort string
	Version      string
	Config       Config
	GuildCount   atomic.Int64
	MemberCount  atomic.Int64
	StartTime    time.Time
}

func (b *Bot) ResetTarget(targetID string, targetType string) {
	if b.BridgeServer != nil {
		b.BridgeServer.ResetTarget(targetID, targetType)
	}
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := slog.Default()
	if os.Getenv("DEBUG_MODE") == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	logger.Info("Starting bot version: " + version)

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:3000"
	}

	port := os.Getenv("INTERNAL_PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logger.Error("Failed to get API_KEY from env")
	}

	apiClient := utils.NewClient(apiUrl, apiKey)

	return &Bot{
		Client:       nil,
		ApiClient:    apiClient,
		Logger:       logger,
		InternalPort: port,
		Version:      version,
		Config: Config{
			Token:       os.Getenv("TOKEN"),
			DevMode:     os.Getenv("DEV_MODE") == "true",
			DevServerID: snowflake.ID(devServerID),
		},
		StartTime: time.Now(),
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) *bot.Client {
	var err error
	b.Client, err = disgo.New(
		b.Config.Token,
		bot.WithShardManagerConfigOpts(
			sharding.WithLogger(b.Logger),
			sharding.WithAutoScaling(true),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(
					gateway.IntentGuilds,
				),
				gateway.WithPresenceOpts(
					gateway.WithPlayingActivity("/help"),
					gateway.WithOnlineStatus(discord.OnlineStatusOnline),
				),
				gateway.WithAutoReconnect(true),
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		bot.WithEventListeners(listeners...),
	)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error while building DisGo client: %s", err))
		os.Exit(1)
	}

	return b.Client
}

func (b *Bot) ReadyEvent(e *events.Ready) {
	b.Logger.Info("Bot shard connected and ready.")
	shardID := e.ShardID()
	shardCount := 0
	for range b.Client.ShardManager.Shards() {
		shardCount++
	}

	b.Logger.Info(fmt.Sprintf("Shard %d/%d is connected!", shardID+1, shardCount))
}

func (b *Bot) OnGuildLeave(e *events.GuildLeave) {
	b.Logger.Info("Left guild, cleaning up data", "guild_id", e.GuildID)
	utils.SendLogMessage(b.Client.Rest,
		fmt.Sprintf(
			"**Left Guild**\n**Server ID:** `%s`\n**Members:** `%d`\n*Deleting data...*",
			e.GuildID,
			e.Guild.MemberCount,
		),
	)

	err := b.ApiClient.DeleteServer(e.GuildID.String())
	if err != nil {
		b.Logger.Error("Failed to auto-delete server data",
			"guild_id", e.GuildID,
			"error", err,
		)
		utils.SendLogMessage(b.Client.Rest, fmt.Sprintf("**Cleanup Failed**\nServer ID: `%s`\nError: `%s`", e.GuildID, err.Error()))
	} else {
		b.Logger.Info("Successfully deleted server data", "guild_id", e.GuildID)
		utils.SendLogMessage(b.Client.Rest, fmt.Sprintf("**Cleanup Complete**\nServer ID: `%s` data deleted.", e.GuildID))
	}
}

func (b *Bot) OnGuildUpdate(e *events.GuildUpdate) {
	if e.Guild.MemberCount == 0 && e.OldGuild.MemberCount > 0 {
		guild := e.Guild
		guild.MemberCount = e.OldGuild.MemberCount
		b.Client.Caches.AddGuild(guild)
	}
}

func (b *Bot) UpdateStats() {
	var gCount, mCount int64
	
	shardCount := 1
	if b.Client != nil && b.Client.ShardManager != nil {
		c := 0
		for range b.Client.ShardManager.Shards() {
			c++
		}
		if c > 0 {
			shardCount = c
		}
	}

	shardStats := make(map[int]bridge.ShardStat)

	for guild := range b.Client.Caches.Guilds() {
		gCount++
		mCount += int64(guild.MemberCount)
		
		sID := sharding.ShardIDByGuild(guild.ID, shardCount)
		stat := shardStats[sID]
		stat.Servers++
		stat.Users += guild.MemberCount
		shardStats[sID] = stat
	}

	b.GuildCount.Store(gCount)
	b.MemberCount.Store(mCount)
	
	if b.BridgeServer != nil {
		for sID, stat := range shardStats {
			b.BridgeServer.UpdateShardStat(sID, stat)
		}
	}
}

func (b *Bot) StartStatsWorker() {
	go func() {
		b.UpdateStats()
		ticker := time.NewTicker(2 * time.Minute)
		for range ticker.C {
			b.UpdateStats()
		}
	}()
}
