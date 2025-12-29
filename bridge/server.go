package bridge

import (
	"encoding/json"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Client bot.Client
	Logger *slog.Logger
	Port   string
}

type BroadcastPayload struct {
	TargetID string        `json:"target_id"`
	Embed    discord.Embed `json:"embed"`
}

func New(client bot.Client, logger *slog.Logger, port string) *Server {
	return &Server{
		Client: client,
		Logger: logger,
		Port:   port,
	}
}

func (s *Server) Start() {
	go func() {
		s.Logger.Info("Starting webhook bridge server on port " + s.Port)
		if err := fasthttp.ListenAndServe(":"+s.Port, s.handleRequest); err != nil {
			s.Logger.Error("Bridge server failed: " + err.Error())
		}
	}()
}

func (s *Server) handleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/internal/broadcast":
		s.handleBroadcast(ctx)
	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}
}

func (s *Server) handleBroadcast(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	var payload BroadcastPayload
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		s.Logger.Error("Failed to parse broadcast payload", "error", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	targetID, err := snowflake.Parse(payload.TargetID)
	if err != nil {
		s.Logger.Error("Invalid Snowflake ID", "id", payload.TargetID)
		ctx.Error("Invalid ID", fasthttp.StatusBadRequest)
		return
	}

	_, err = s.Client.Rest().CreateMessage(targetID, discord.MessageCreate{
		Embeds: []discord.Embed{payload.Embed},
	})

	if err != nil {
		s.Logger.Error("Failed to send to Discord", "error", err)
		ctx.Error("Discord API Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
