package utils

type mDbManga struct {
	Title     string `bson:"title"`
	Id        string `bson:"id"`
	GroupName string `bson:"groupName"`
	GroupId   string `bson:"groupId"`
}

type mDbServer struct {
	Id         string     `bson:"_id"`
	ServerId   string     `bson:"serverid"`
	ServerName string     `bson:"serverName"`
	ChannelId  string     `bson:"channelid"`
	Manga      []mDbManga `bson:"manga"`
}

type mDbUser struct {
	Id       string     `bson:"_id"`
	UserId   string     `bson:"userid"`
	Username string     `bson:"username"`
	Manga    []mDbManga `bson:"manga"`
}
