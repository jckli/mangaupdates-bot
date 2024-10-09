package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jckli/mangaupdates-bot/mubot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbName   = os.Getenv("MONGO_DB_NAME")
	serverMu sync.Mutex
	userMu   sync.Mutex
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

func dbGetNextID(ctx context.Context, collection *mongo.Collection, mu *sync.Mutex) (int32, error) {
	mu.Lock()
	defer mu.Unlock()

	expectedID := int32(0)

	// Create a cursor to iterate over _id's in ascending order
	findOptions := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetProjection(bson.D{{Key: "_id", Value: 1}})

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return 0, fmt.Errorf("Failed to find _id's: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			ID int32 `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("Failed to decode _id: %w", err)
		}

		if result.ID == expectedID {
			expectedID++
		} else if result.ID > expectedID {
			break
		}
	}

	if err := cursor.Err(); err != nil {
		return 0, fmt.Errorf("Cursor error: %w", err)
	}

	return expectedID, nil
}

func DbAddServer(b *mubot.Bot, serverName string, serverId, channelId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("servers")

	const maxRetries = 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		newID, err := dbGetNextID(ctx, collection, &serverMu)
		if err != nil {
			return fmt.Errorf("Error getting next server ID: %w", err)
		}

		doc := bson.M{
			"_id":        newID,
			"serverid":   serverId,
			"serverName": serverName,
			"channelid":  channelId,
			"manga":      []interface{}{},
		}

		_, err = collection.InsertOne(ctx, doc)
		if err != nil {
			var writeException mongo.WriteException
			if errors.As(err, &writeException) {
				duplicateKey := false
				for _, we := range writeException.WriteErrors {
					if we.Code == 11000 {
						duplicateKey = true
						break
					}
				}
				if duplicateKey {
					b.Logger.Error(fmt.Sprintf(
						"Duplicate _id %d detected for servers. Retrying (Attempt %d/%d)...",
						newID,
						attempt+1,
						maxRetries,
					))
					continue
				}
			}
			return fmt.Errorf("Failed to insert server: %w", err)
		}

		return nil
	}

	return fmt.Errorf("Failed to insert server after %d attempts due to duplicate _id", maxRetries)
}

func DbAddUser(b *mubot.Bot, username string, userId int64) error {
	collection := b.MongoClient.Database(dbName).Collection("users")

	const maxRetries = 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		newID, err := dbGetNextID(ctx, collection, &userMu)
		if err != nil {
			return fmt.Errorf("Error getting next user ID: %w", err)
		}

		doc := bson.M{
			"_id":      newID,
			"userid":   userId,
			"username": username,
			"manga":    []interface{}{},
		}

		_, err = collection.InsertOne(ctx, doc)
		if err != nil {
			var writeException mongo.WriteException
			if errors.As(err, &writeException) {
				duplicateKey := false
				for _, we := range writeException.WriteErrors {
					if we.Code == 11000 { // Duplicate Key Error Code
						duplicateKey = true
						break
					}
				}
				if duplicateKey {
					b.Logger.Error(fmt.Sprintf(
						"Duplicate _id %d detected for users. Retrying (Attempt %d/%d)...",
						newID,
						attempt+1,
						maxRetries,
					))
					continue
				}
			}
			return fmt.Errorf("Failed to insert user: %w", err)
		}

		return nil
	}

	return fmt.Errorf("Failed to insert user after %d attempts due to duplicate _id", maxRetries)
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
