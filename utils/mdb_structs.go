package utils

type MDbManga struct {
	Title     string `bson:"title"`
	Id        string `bson:"id"`
	GroupName string `bson:"groupName"`
	GroupId   string `bson:"groupId"`
}

type MDbServer struct {
	Id         string     `bson:"_id"`
	ServerId   string     `bson:"serverid"`
	ServerName string     `bson:"serverName"`
	ChannelId  string     `bson:"channelid"`
	Manga      []MDbManga `bson:"manga"`
}

type MDbUser struct {
	Id       string     `bson:"_id"`
	UserId   string     `bson:"userid"`
	Username string     `bson:"username"`
	Manga    []MDbManga `bson:"manga"`
}
