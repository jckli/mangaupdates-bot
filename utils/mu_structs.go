package utils

type MangaSearchResult struct {
	ID       int64   `json:"id"`
	Base36ID string  `json:"base36_id"`
	Title    string  `json:"title"`
	Year     string  `json:"year"`
	Rating   float64 `json:"rating"`
}

type GroupSearchResult struct {
	ID   int64  `json:"group_id"`
	Name string `json:"name"`
}

type MangaMetadata struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Year          string   `json:"year"`
	Rating        *float64 `json:"rating"`
	Type          string   `json:"kind"`
	Status        string   `json:"status"`
	LatestChapter string   `json:"latest_chapter"`
	FinalChapter  string   `json:"final_chapter"`
	Authors       []string `json:"authors"`
	Artists       []string `json:"artists"`
	CoverURL      string   `json:"cover_url"`
}

type GroupDetails struct {
	GroupID int64       `json:"group_id"`
	Name    string      `json:"name"`
	URL     string      `json:"url"`
	Social  GroupSocial `json:"social"`
	Active  bool        `json:"active"`
}

type GroupSocial struct {
	Site    string `json:"site"`
	Discord string `json:"discord"`
	Twitter string `json:"twitter"`
}

type TrackedManga struct {
	Title     string `json:"title"`
	ID        int64  `json:"id"`
	GroupName string `json:"groupName,omitempty"`
	GroupID   int64  `json:"groupid,omitempty"`
}

type SetupServerRequest struct {
	ServerName string `json:"server_name"`
	ChannelID  string `json:"channel_id"`
}

type SetupUserRequest struct {
	Username string `json:"username"`
}

type AddMangaRequest struct {
	ID int64 `json:"id"`
}

type SetGroupRequest struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
}

type SetChannelRequest struct {
	ChannelID string `json:"channel_id"`
}

type SetRoleRequest struct {
	RoleID   int64  `json:"role_id"`
	RoleType string `json:"type"`
}

type ServerConfig struct {
	Roles struct {
		Admin int64 `json:"admin,omitempty"`
	} `json:"roles"`
	ChannelID int64 `json:"channelid,omitempty"`
}

type UserConfig struct {
}
