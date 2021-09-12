import pymongo
from pymongo import MongoClient
from pymongo import collection

import json
import os

with open("config.json", "r") as f:
    config = json.load(f)

dbPassword = config["dbPassword"]
mongo = MongoClient(f"mongodb+srv://baka:{dbPassword}@akane.dsytm.mongodb.net/Kana?retryWrites=true&w=majority")

database = mongo["Kana"]
srv = database["servers"]
usr = database["users"]

def addServer(serverName, serverID, channelID):
    srv.insert_one({"serverid": serverID, "serverName": serverName, "channelid": channelID, "manga": []})

def addUser(userName, userID):
    usr.insert_one({"userid": userID, "userName": userName, "manga": []})

def removeServer(serverID):
    srv.delete_one({"serverid": serverID})

def removeUser(userID):
    usr.delete_one({"userid": userID})

def updateChannel(idName, channelID):
    srv.update_one({"serverid": idName}, {"$set": {"channelid": channelID}})

def checkServerExist(idName):
    result = srv.find_one({"serverid": idName})
    if result != None:
        return True
    else:
        return False

def addManga(idName, title, mode):
    if mode == "user":
        usr.update_one({"userid": idName}, {"$push": {"manga": {"title": title}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$push": {"manga": {"title": title}}})
    
def removeManga(idName, title, mode):
    if mode == "user":
        usr.update_one({"userid": idName}, {"$pull": {"manga": {"title": title}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$pull": {"manga": {"title": title}}})

def mangaWanted(title, mode):
    if mode == "user":
        result = srv.find_one({"manga": {"$elemMatch": {"title": title}}})
    elif mode == "server":
        result = srv.find_one({"manga": {"$elemMatch": {"title": title}}})

    if result != None:
        return True
    else:
        return False

def test(title):
    sheeesh = []
    for document in srv.find():
        print(document[2])

#        if result != None:
#            sheeesh.append(result.serverid)
#    return result