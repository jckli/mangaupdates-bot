package bridge

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
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
	Client              *bot.Client
	Logger              *slog.Logger
	Port                string
	updateChan          chan BroadcastPayload
	GlobalFeedChannelID string
	StartTime           time.Time
	shardStats          sync.Map // map[int]ShardStat
}

type ShardStat struct {
	Servers int
	Users   int
}

func (s *Server) UpdateShardStat(shardID int, stat ShardStat) {
	s.shardStats.Store(shardID, stat)
}

type BroadcastPayload struct {
	TargetID   string   `json:"target_id"`
	TargetType string   `json:"target_type"`
	Title      string   `json:"title"`
	Chapter    string   `json:"chapter"`
	Link       string   `json:"link"`
	ImageURL   string   `json:"image_url"`
	Groups     string   `json:"groups"`
	RoleIDs    []string `json:"role_ids,omitempty"`
}

func (s *Server) logToOps(msg string) {
	utils.SendLogMessage(s.Client.Rest, msg)
}

func New(client *bot.Client, logger *slog.Logger, port string) *Server {
	return &Server{
		Client:              client,
		Logger:              logger,
		Port:                port,
		updateChan:          make(chan BroadcastPayload, QueueSize),
		GlobalFeedChannelID: os.Getenv("GLOBAL_FEED_CHANNEL_ID"),
		StartTime:           time.Now(),
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
	if string(ctx.Path()) == "/internal/status" && ctx.IsGet() {
		s.handleStatus(ctx)
		return
	}

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
	if self, ok := s.Client.Caches.SelfUser(); ok {
		botIcon = self.EffectiveAvatarURL()
	}

	var channelID snowflake.ID = targetID
	isDM := payload.TargetType == "user"
	if isDM {
		ch, err := s.Client.Rest.CreateDMChannel(targetID)
		if err != nil {
			s.Logger.Error("Failed to create DM", "user_id", targetID, "error", err)
			s.logToOps(fmt.Sprintf("**DM Error**: Could not open DM with User `%s`\nError: `%v`", payload.TargetID, err))
			return
		}
		channelID = ch.ID()
	}

	embed := discord.NewEmbed().
		WithAuthor("MangaUpdates", "", botIcon).
		WithDescriptionf("Chapter `%s` has been released for `%s`!", payload.Chapter, payload.Title).
		WithTitlef("New %s Chapter!", payload.Title).
		WithURL(payload.Link).
		WithColor(ColorPrimary).
		AddField("Chapter", payload.Chapter, true).
		AddField("Scanlators", payload.Groups, true).
		WithImage(payload.ImageURL).
		WithTimestamp(time.Now())

	var content string
	if len(payload.RoleIDs) > 0 {
		var pings []string
		for _, rid := range payload.RoleIDs {
			pings = append(pings, fmt.Sprintf("<@&%s>", rid))
		}
		content = strings.Join(pings, " ")
	}

	message, err := s.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content: content,
		Embeds:  []discord.Embed{embed},
	})

	if err != nil {
		s.Logger.Error("Failed to send to Discord", "channel_id", channelID, "error", err)
		s.logToOps(fmt.Sprintf("**Delivery Failed**: `%s` Ch.%s -> %s (`%s`)\nError: `%v`",
			payload.Title,
			payload.Chapter,
			payload.TargetType,
			payload.TargetID,
			err,
		))
	} else {
		if !isDM && s.GlobalFeedChannelID != "" && payload.TargetID == s.GlobalFeedChannelID {
			go func() {
				_, err := s.Client.Rest.CrosspostMessage(channelID, message.ID)
				if err != nil {
					s.Logger.Error("Failed to auto-publish global feed", "error", err)
				}
			}()
		}
		s.logToOps(fmt.Sprintf("**Sent**: `%s` Ch.%s -> %s (`%s`)", payload.Title, payload.Chapter, payload.TargetType, payload.TargetID))
	}
}

type ShardInfo struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Latency int    `json:"latency"`
	Servers int    `json:"servers"`
	Users   int    `json:"users"`
}

type StatusData struct {
	Uptime        string      `json:"uptime"`
	MemoryUsageMB float64     `json:"memory_usage_mb"`
	TotalServers  int         `json:"total_servers"`
	Shards        []ShardInfo `json:"shards"`
}

func (s *Server) handleStatus(ctx *fasthttp.RequestCtx) {
	uptime := time.Since(s.StartTime).Round(time.Second)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memMB := float64(memStats.Alloc) / 1024 / 1024

	var totalServers int
	if s.Client != nil && s.Client.Caches != nil {
		totalServers = s.Client.Caches.GuildsLen()
	}

	shards := []ShardInfo{}
	if s.Client != nil && s.Client.ShardManager != nil {
		for shard := range s.Client.ShardManager.Shards() {
			statusStr := ""
			switch shard.Status() {
			case 0:
				statusStr = "Unconnected"
			case 1:
				statusStr = "Connecting"
			case 2:
				statusStr = "WaitingForHello"
			case 3:
				statusStr = "Identifying"
			case 4:
				statusStr = "Resuming"
			case 5:
				statusStr = "Ready"
			case 6:
				statusStr = "Disconnected"
			case 7:
				statusStr = "Fatal"
			default:
				statusStr = fmt.Sprintf("Unknown(%d)", shard.Status())
			}

			sCount, uCount := 0, 0
			if val, ok := s.shardStats.Load(shard.ShardID()); ok {
				if stat, ok := val.(ShardStat); ok {
					sCount = stat.Servers
					uCount = stat.Users
				}
			}
			
			shards = append(shards, ShardInfo{
				ID:      shard.ShardID(),
				Status:  statusStr,
				Latency: int(shard.Latency().Milliseconds()),
				Servers: sCount,
				Users:   uCount,
			})
		}
	}

	data := StatusData{
		Uptime:        uptime.String(),
		MemoryUsageMB: float64(int(memMB*10)) / 10, // 1 decimal point
		TotalServers:  totalServers,
		Shards:        shards,
	}

	body, _ := json.Marshal(data)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
}
