from pymongo import MongoClient

import json
import certifi

from core import update

with open("config.json", "r") as f:
    config = json.load(f)

ca = certifi.where()

dbPassword = config["dbPassword"]
mongo = MongoClient(f"mongodb+srv://baka:{dbPassword}@akane.dsytm.mongodb.net/Kana?retryWrites=true&w=majority", tlsCAFile=ca)

database = mongo["Kana"]
srv = database["servers"]
usr = database["users"]

def addServer(serverName, serverID, channelID):
    documentCount = srv.count_documents({})
    while documentCount >= 0:
        try:
            srv.insert_one({"_id": documentCount, "serverid": serverID, "serverName": serverName, "channelid": channelID, "manga": []})
            break
        except:
            documentCount -= 1

def addUser(userName, userID):
    documentCount = usr.count_documents({})
    while documentCount >= 0:
        try:
            usr.insert_one({"_id": documentCount, "userid": userID, "username": userName, "manga": []})
            break
        except:
            documentCount -= 1

def removeServer(serverID):
    srv.delete_one({"serverid": serverID})

def removeUser(userID):
    usr.delete_one({"userid": userID})

def updateChannel(idName, channelID):
    srv.update_one({"serverid": idName}, {"$set": {"channelid": channelID}})

def findChannel(idName):
    result = srv.find_one({"serverid": idName})
    return result["channelid"]

def checkServerExist(idName):
    result = srv.find_one({"serverid": idName})
    if result != None:
        return True
    else:
        return False

def checkUserExist(idName):
    result = usr.find_one({"userid": idName})
    if result != None:
        return True
    else:
        return False

def addManga(idName, title, link, mode):
    mangaid = link.partition("https://www.mangaupdates.com/series.html?id=")[2]
    if mode == "user":
        usr.update_one({"userid": idName}, {"$push": {"manga": {"title": title, "id": mangaid}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$push": {"manga": {"title": title, "id": mangaid}}})
    
def removeManga(idName, title, mode):
    if mode == "user":
        usr.update_one({"userid": idName}, {"$pull": {"manga": {"title": title}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$pull": {"manga": {"title": title}}})

def mangaWanted(title, group, mode):
    if "&" in group:
        group = group.split("&")
        group = [x.strip(' ') for x in group]
    if mode == "user":
        idList = []
        result = usr.find({"manga.title": title}, {"_id": 0, "userid": 1, "manga": 1})
        for i in result:
            for m in i["manga"]:
                if m["title"] == title:
                    if "groupName" in m:
                        if m["groupName"] in group:
                            idList.append(i["userid"])
                    elif "groupName" not in m:
                        idList.append(i["userid"])
                    break
        if idList != []:
            return idList
        else:
            return None
    elif mode == "server":
        class list:
            serverList = []
            channelList = []
        result = srv.find({"manga.title": title}, {"_id": 0, "serverid": 1, "channelid": 1, "manga": 1})
        for i in result:
            for m in i["manga"]:
                if m["title"] == title:
                    if "groupName" in m:
                        if m["groupName"] in group:
                            list.serverList.append(i["serverid"])
                            list.channelList.append(i["channelid"])
                    elif "groupName" not in m:
                        list.serverList.append(i["serverid"])
                        list.channelList.append(i["channelid"])
                    break
        if list.serverList != []:
            return list
        else:
            return None

def checkMangaAlreadyWithinDb(id, link, mode):
    mangaid = link.partition("https://www.mangaupdates.com/series.html?id=")[2]
    if mode == "user":
        result = usr.find_one({"userid": id}, {"manga": 1})
        for i in result["manga"]:
            if i["id"] == mangaid:
                return True
        return False
    elif mode == "server":
        result = srv.find_one({"serverid": id}, {"manga": 1})
        for i in result["manga"]:
            if i["id"] == mangaid:
                return True
        return False

def getMangaList(id, mode):
    manga = []
    if mode == "user":
        result = usr.find_one({"userid": id}, {"manga": 1})
        for i in result["manga"]:
            manga.append(i["title"])
        if manga != []:
            return manga
        else:
            return None
    elif mode == "server":
        result = srv.find_one({"serverid": id}, {"manga": 1})
        for i in result["manga"]:
            manga.append(i["title"])
        if manga != []:
            return manga
        else:
            return None
        
def getMangaID(title, mode):
    if mode == "user":
        result = usr.find_one({"manga.title": title}, {"manga": 1})
    elif mode == "server":
        result = srv.find_one({"manga.title": title}, {"manga": 1})
    for i in result["manga"]:
        if i["title"] == title:
            return i["id"]

def setScanGroup(id, title, group, groupid, mode):
    if mode == "user":
        usr.update_one({"userid": id, "manga.title": title}, {"$set": {"manga.$.groupName": group, "manga.$.groupid": groupid}})
    elif mode == "server":
        srv.update_one({"serverid": id, "manga.title": title}, {"$set": {"manga.$.groupName": group, "manga.$.groupid": groupid}})
    
def removeScanGroup(id, title, mode):
    if mode == "user":
        usr.update_one({"userid": id, "manga.title": title}, {"$unset": {"manga.$.groupName": "", "manga.$.groupid": ""}})
    elif mode == "server":
        srv.update_one({"serverid": id, "manga.title": title}, {"$unset": {"manga.$.groupName": "", "manga.$.groupid": ""}})

def getAllIds(mode):
    if mode == "user":
        result = usr.find({}, {"userid": 1})
    elif mode == "server":
        result = srv.find({}, {"serverid": 1})
    idList = []
    for i in result:
        if mode == "user":
            idList.append(i["userid"])
        elif mode == "server":
            idList.append(i["serverid"])
    return idList   

def test():
    documentCount = usr.count_documents({})
    print(documentCount)

#        if result != None:
#            sheeesh.append(result.serverid)
#    return result