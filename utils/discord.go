package utils

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

const LogChannelID = "990005048408936529"

func SendLogMessage(client rest.Rest, content string) {
	go func() {
		destID, err := snowflake.Parse(LogChannelID)
		if err != nil {
			return
		}
		_, _ = client.CreateMessage(destID, discord.MessageCreate{Content: content})
	}()
}
