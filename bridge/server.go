package bridge

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jckli/mangaupdates-bot/utils"
	"github.com/valyala/fasthttp"
)

var ColorPrimary = 0x3083e3

const (
	QueueSize   = 5000
	WorkerCount = 5
)

type Server struct {
	Client     bot.Client
	Logger     *slog.Logger
	Port       string
	updateChan chan BroadcastPayload
}

type BroadcastPayload struct {
	TargetID   string `json:"target_id"`
	TargetType string `json:"target_type"`
	Title      string `json:"title"`
	Chapter    string `json:"chapter"`
	Link       string `json:"link"`
	ImageURL   string `json:"image_url"`
	Groups     string `json:"groups"`
}

func (s *Server) logToOps(msg string) {
	utils.SendLogMessage(s.Client.Rest(), msg)
}

func New(client bot.Client, logger *slog.Logger, port string) *Server {
	return &Server{
		Client:     client,
		Logger:     logger,
		Port:       port,
		updateChan: make(chan BroadcastPayload, QueueSize),
	}
}

func (s *Server) Start() {
	for range WorkerCount {
		go s.processUpdates()
	}

	go func() {
		s.Logger.Info("Starting webhook bridge server on port " + s.Port)
		s.logToOps("**Bridge Server Started**")
		if err := fasthttp.ListenAndServe(":"+s.Port, s.handleRequest); err != nil {
			s.Logger.Error("Bridge server failed: " + err.Error())
		}
	}()
}

func (s *Server) handleRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/internal/broadcast" {
		ctx.Error("Not Found", fasthttp.StatusNotFound)
		return
	}

	if !ctx.IsPost() {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	var payload BroadcastPayload
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		errMsg := fmt.Sprintf("**API Error**: Bad JSON received\n```%v```", err)
		s.logToOps(errMsg)
		s.Logger.Error("Failed to parse broadcast payload", "error", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	s.updateChan <- payload
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *Server) processUpdates() {
	for payload := range s.updateChan {
		s.sendToDiscord(payload)
	}
}

func (s *Server) sendToDiscord(payload BroadcastPayload) {
	targetID, err := snowflake.Parse(payload.TargetID)
	if err != nil {
		s.Logger.Error("Invalid Snowflake ID", "id", payload.TargetID)
		s.logToOps(fmt.Sprintf("**Error**: Invalid Snowflake ID `%s` for series `%s`", payload.TargetID, payload.Title))
		return
	}

	botIcon := ""
	if self, ok := s.Client.Caches().SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	var channelID snowflake.ID = targetID
	if payload.TargetType == "user" {
		ch, err := s.Client.Rest().CreateDMChannel(targetID)
		if err != nil {
			s.Logger.Error("Failed to create DM", "user_id", targetID, "error", err)
			s.logToOps(fmt.Sprintf("**DM Error**: Could not open DM with User `%s`\nError: `%v`", payload.TargetID, err))
			return
		}
		channelID = ch.ID()
	}

	embed := discord.NewEmbedBuilder().
		SetAuthor("MangaUpdates", "", botIcon).
		SetDescriptionf("Chapter `%s` has been released for `%s`!\n\n_Note: Sources are now linked directly in the scanlator names below._", payload.Chapter, payload.Title).
		SetTitlef("New %s Chapter!", payload.Title).
		SetURL(payload.Link).
		SetColor(ColorPrimary).
		AddField("Chapter", payload.Chapter, true).
		AddField("Scanlators", payload.Groups, true).
		SetImage(payload.ImageURL).
		SetTimestamp(time.Now()).
		Build()

	_, err = s.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	})

	if err != nil {
		s.Logger.Error("Failed to send to Discord", "channel_id", channelID, "error", err)
		s.logToOps(fmt.Sprintf("**Delivery Failed**: `%s` Ch.%s -> Channel `%s`\nError: `%v`", payload.Title, payload.Chapter, channelID, err))
	} else {
		s.logToOps(fmt.Sprintf("**Sent**: `%s` Ch.%s -> %s (`%s`)", payload.Title, payload.Chapter, payload.TargetType, channelID))
	}
}
