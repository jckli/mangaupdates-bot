package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/bot"
	"github.com/jckli/mangaupdates-bot/bridge"
	"github.com/jckli/mangaupdates-bot/commands"
	"github.com/jckli/mangaupdates-bot/mubot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mu := mubot.New(os.Getenv("VERSION"))

	h := commands.CommandHandlers(mu)

	mu.Setup(
		h,
		bot.NewListenerFunc(mu.ReadyEvent),
	)

	var err error
	if mu.Config.DevMode {
		mu.Logger.Info(
			fmt.Sprintf(
				"Running in dev mode. Syncing commands to server ID: %s",
				mu.Config.DevServerID,
			),
		)
		_, err = mu.Client.Rest().
			SetGuildCommands(mu.Client.ApplicationID(), mu.Config.DevServerID, commands.CommandList)
		if err == nil {
			_, err = mu.Client.Rest().SetGlobalCommands(mu.Client.ApplicationID(), nil)
		}
	} else {
		mu.Logger.Info(
			"Running in global mode. Syncing commands globally.",
		)
		_, err = mu.Client.Rest().SetGlobalCommands(mu.Client.ApplicationID(), commands.CommandList)
		if err == nil {
			_, err = mu.Client.Rest().SetGuildCommands(mu.Client.ApplicationID(), mu.Config.DevServerID, nil)
		}
	}
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Failed to sync commands: %s", err))
	}

	bridgeServer := bridge.New(mu.Client, mu.Logger, mu.InternalPort)
	bridgeServer.Start()

	err = mu.Client.OpenShardManager(context.TODO())
	if err != nil {
		mu.Logger.Error(fmt.Sprintf("Error while connecting: %s", err))
	}
	defer mu.Client.Close(context.TODO())

	mu.Logger.Info("Bot is now running.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	mu.Logger.Info("Shutting down...")
}
