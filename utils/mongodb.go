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
