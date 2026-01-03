package mubot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/mangaupdates-bot/utils"
)

type Config struct {
	Token       string
	DevMode     bool
	DevServerID snowflake.ID
}

type Bot struct {
	Client       bot.Client
	ApiClient    *utils.Client
	Logger       *slog.Logger
	InternalPort string
	Version      string
	Config       Config
}

func New(version string) *Bot {
	devServerID, _ := strconv.Atoi(os.Getenv("DEV_SERVER_ID"))

	logger := slog.Default()
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
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) bot.Client {
	var err error
	b.Client, err = disgo.New(
		b.Config.Token,
		bot.WithLogger(b.Logger),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		bot.WithEventListeners(listeners...),
	)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error while building DisGo client: %s", err))
	}

	return b.Client
}

func (b *Bot) ReadyEvent(_ *events.Ready) {
	err := b.Client.SetPresence(
		context.TODO(),
		gateway.WithPlayingActivity("✨ Rewrite update | /alert | /help"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("Error while setting presence: %s", err))
	}

	b.Logger.Info("Bot presence set successfully.")
}
