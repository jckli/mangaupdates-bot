package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/jckli/mangaupdates-bot/commands"
	"github.com/jckli/mangaupdates-bot/mubot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mu := mubot.New(os.Getenv("VERSION"))

	h := commands.CommandHandlers(mu)

	client := mu.Setup(
		h,
		bot.NewListenerFunc(mu.ReadyEvent),
	)
	mu.Client = client

	var err error
	if mu.Config.DevMode {
		mu.Logger.Info(
			fmt.Sprintf(
				"Running in dev mode. Syncing commands to server ID: %s",
				mu.Config.DevServerID,
			),
		)
	} else {
		mu.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
		_, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands.CommandList)
	}
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to sync commands: %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.OpenGateway(ctx)
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Error while connecting: %s", err))
	}
	defer client.Close(context.TODO())

	mu.Logger.Info("Bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	mu.Logger.Info("Shutting down...")
}
