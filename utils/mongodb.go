package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jckli/mangaupdates-bot/mubot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbName = os.Getenv("MONGO_DB_NAME")
)

func DbConnect() (*mongo.Client, error) {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	mongoHost := os.Getenv("MONGO_HOST")
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", mongoUser, mongoPass, mongoHost)),
	)

	return client, err
}

func DbDisconnect(b *mubot.Bot) error {
	if err := b.MongoClient.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}

func DbAddServer(b *mubot.Bot, serverName string, serverId, channelId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"serverid":   serverId,
		"serverName": serverName,
		"channelid":  channelId,
		"manga":      []interface{}{},
	}

	_, err := collection.InsertOne(context.TODO(), doc)
	return err
}

func DbAddUser(b *mubot.Bot, username string, userId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("users")

	doc := bson.M{
		"_id":      primitive.NewObjectID(),
		"userid":   userId,
		"username": username,
		"manga":    []interface{}{},
	}

	_, err := collection.InsertOne(context.TODO(), doc)
	return err
}

func DbRemoveServer(b *mubot.Bot, serverId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"serverid": serverId})

	return err
}

func DbRemoveUser(b *mubot.Bot, userId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("users")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"userid": userId})

	return err
}

func DbGetServer(b *mubot.Bot, serverId int64) (bson.M, error) {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"serverid": serverId}).Decode(&result)

	return result, err
}

func DbGetUser(b *mubot.Bot, userId int64) (bson.M, error) {
	collection := b.MongoClient.Database(dbName).Collection("users")

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"userid": userId}).Decode(&result)

	return result, err
}

func DbSetChannel(b *mubot.Bot, serverId, channelId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"serverid": serverId},
		bson.M{"$set": bson.M{"channelid": channelId}},
	)

	return err
}

func DbGetChannel(b *mubot.Bot, serverId int64) (int64, error) {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"serverid": serverId}).Decode(&result)

	return result["channelid"].(int64), err
}

func DbServersWanted(
	b *mubot.Bot,
	groupList *[]MuSearchGroupsGroup,
	entry *MangaEntry,
) ([]MDbServer, error) {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	var filter bson.M
	if &entry.NewId != nil {
		filter = bson.M{"manga.id": entry.NewId}
	} else if entry.Title != "" {
		filter = bson.M{"manga.title": entry.Title}
	} else {
		return nil, fmt.Errorf("No entry to search for")
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []MDbServer
	for cursor.Next(context.Background()) {
		var result MDbServer
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		manga := result.Manga[0]
		if manga.GroupId != 0 {
			for _, group := range *groupList {
				if manga.GroupId == int64(group.Record.GroupID) {
					results = append(results, result)
					break
				}
			}
		} else {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func DbUsersWanted(
	b *mubot.Bot,
	groupList *[]MuSearchGroupsGroup,
	entry *MangaEntry,
) ([]MDbUser, error) {
	collection := b.MongoClient.Database(dbName).Collection("users")

	var filter bson.M
	if &entry.NewId != nil {
		filter = bson.M{"manga.id": entry.NewId}
	} else if entry.Title != "" {
		filter = bson.M{"manga.title": entry.Title}
	} else {
		return nil, fmt.Errorf("No entry to search for")
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []MDbUser
	for cursor.Next(context.Background()) {
		var result MDbUser
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		manga := result.Manga[0]
		if manga.GroupId != 0 {
			for _, group := range *groupList {
				if manga.GroupId == int64(group.Record.GroupID) {
					results = append(results, result)
					break
				}
			}
		} else {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func DbServerAddManga(b *mubot.Bot, serverId int64, manga MDbManga) error {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"serverid": serverId},
		bson.M{"$push": bson.M{"manga": manga}},
	)

	return err
}

func DbUserAddManga(b *mubot.Bot, userId int64, manga MDbManga) error {
	collection := b.MongoClient.Database(dbName).Collection("users")

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"userid": userId},
		bson.M{"$push": bson.M{"manga": manga}},
	)

	return err
}

func DbServerCheckMangaExists(b *mubot.Bot, serverId, mangaId int64) (bool, error) {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	count, err := collection.CountDocuments(
		context.TODO(),
		bson.M{"serverid": serverId, "manga.id": mangaId},
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DbUserCheckMangaExists(b *mubot.Bot, userId, mangaId int64) (bool, error) {
	collection := b.MongoClient.Database(dbName).Collection("users")

	count, err := collection.CountDocuments(
		context.TODO(),
		bson.M{"userid": userId, "manga.id": mangaId},
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DbServerCheckExists(b *mubot.Bot, serverId int64) (bool, error) {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	count, err := collection.CountDocuments(
		context.TODO(),
		bson.M{"serverid": serverId},
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DbUserCheckExists(b *mubot.Bot, userId int64) (bool, error) {
	collection := b.MongoClient.Database(dbName).Collection("users")

	count, err := collection.CountDocuments(
		context.TODO(),
		bson.M{"userid": userId},
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
