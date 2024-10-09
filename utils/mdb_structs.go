package utils

type MDbManga struct {
	Title     string `bson:"title"`
	Id        int64  `bson:"id"`
	GroupName string `bson:"groupName"`
	GroupId   int64  `bson:"groupid"`
}

type MDbServer struct {
	Id         int32      `bson:"_id"`
	ServerId   int64      `bson:"serverid"`
	ServerName string     `bson:"serverName"`
	ChannelId  int64      `bson:"channelid"`
	Manga      []MDbManga `bson:"manga"`
}

type MDbUser struct {
	Id       int32      `bson:"_id"`
	UserId   int64      `bson:"userid"`
	Username string     `bson:"username"`
	Manga    []MDbManga `bson:"manga"`
}

type MDbCounter struct {
	ID  string `bson:"_id"`
	Seq int32  `bson:"seq"`
}
