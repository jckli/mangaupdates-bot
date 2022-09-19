from pymongo import MongoClient
import certifi
import os
import requests
from bs4 import BeautifulSoup as bs
import time
import re

class Mongo:
    def __init__(self):
        ca = certifi.where()
        username = os.environ.get("MONGO_USER")
        password = os.environ.get("MONGO_PASS")
        database_name = os.environ.get("MONGO_DB_NAME")
        mongo = MongoClient(f"mongodb+srv://{username}:{password}@akane.dsytm.mongodb.net/{database_name}?retryWrites=true&w=majority", tlsCAFile=ca)

        db = mongo[database_name]
        self.usr = db["users"]
        self.srv = db["servers"]
    
    async def add_server(self, server_name, server_id, channel_id):
        document_count = self.srv.count_documents({})
        while document_count >= 0:
            try:
                self.srv.insert_one({"_id": document_count, "serverid": server_id, "serverName": server_name, "channelid": channel_id, "manga": []})
                break
            except:
                document_count -= 1

    async def add_user(self, user_name, user_id):
        document_count = self.usr.count_documents({})
        while document_count >= 0:
            try:
                self.usr.insert_one({"_id": document_count, "userid": user_id, "username": user_name, "manga": []})
                break
            except:
                document_count -= 1

    async def remove_server(self, server_id):
        self.srv.delete_one({"serverid": server_id})

    async def remove_user(self, user_id):
        self.usr.delete_one({"userid": user_id})

    async def get_server(self, server_id):
        return self.srv.find_one({"serverid": server_id})

    async def get_user(self, user_id):
        return self.usr.find_one({"userid": user_id})

    async def set_channel(self, server_id, channel_id):
        self.srv.update_one({"serverid": server_id}, {"$set": {"channelid": channel_id}})

    async def get_channel(self, server_id):
        result = self.srv.find_one({"serverid": server_id}, {"channelid": 1})
        return result["channelid"]

    async def check_server_exist(self, server_id):
        result = await Mongo.get_server(self, server_id)
        if result is not None:
            return True
        else:
            return False
    
    async def check_user_exist(self, user_id):
        result = await Mongo.get_user(self, user_id)
        if result is not None:
            return True
        else:
            return False

    async def check_manga_exist_server(self, server_id, manga_id):
        result = self.srv.find_one({"serverid": server_id}, {"manga": 1})
        for i in result["manga"]:
            if i["id"] == manga_id:
                return True
        return False
    
    async def check_manga_exist_user(self, user_id, manga_id):
        result = self.usr.find_one({"userid": user_id}, {"manga": 1})
        for i in result["manga"]:
            if i["id"] == manga_id:
                return True
        return False
    
    async def add_manga_server(self, server_id, manga_id, manga_name):
        self.srv.update_one({"serverid": server_id}, {"$push": {"manga": {"title": manga_name, "id": manga_id}}})
    
    async def add_manga_user(self, user_id, manga_id, manga_name):
        self.usr.update_one({"userid": user_id}, {"$push": {"manga": {"title": manga_name, "id": manga_id}}})

    async def add_add_role_server(self, server_id, discord_role_id):
        self.srv.update_one({"serverid": server_id}, {"$set": {"role": discord_role_id}})

    async def get_manga_list_server(self, server_id):
        manga = []
        result = self.srv.find_one({"serverid": server_id}, {"manga": 1})
        for i in result["manga"]:
            manga.append({"id": i["id"], "title": i["title"]})
        if manga != []:
            return manga
        else:
            return None
    
    async def get_server_allow_add_role(self, server_id):
        result = self.srv.find_one({'serverid': server_id}, {'role': 1})
        return result.get('role', '')

    async def get_manga_list_user(self, user_id):
        manga = []
        result = self.usr.find_one({"userid": user_id}, {"manga": 1})
        for i in result["manga"]:
            manga.append({"id": i["id"], "title": i["title"]})
        if manga != []:
            return manga
        else:
            return None

    async def remove_manga_server(self, server_id, manga_id):
        self.srv.update_one({"serverid": server_id}, {"$pull": {"manga": {"id": manga_id}}})
    
    async def remove_manga_user(self, user_id, manga_id):
        self.usr.update_one({"userid": user_id}, {"$pull": {"manga": {"id": manga_id}}})

    async def remove_add_role_server(self, server_id):
        self.srv.update_one({"serverid": server_id}, {"$set": {"role": ''}})

    async def manga_wanted_server(self, group_list, manga_id=None, manga_title=None):
        serverList = []
        if manga_title is not None:
            result = self.srv.find({"manga.title": manga_title}, {"serverid": 1, "channelid": 1, "manga.$": 1})
        else:
            result = self.srv.find({"manga.id": manga_id}, {"serverid": 1, "channelid": 1, "manga.$": 1})
        for i in result:
            if "groupid" in i["manga"][0]:
                for group in group_list:
                    if i["manga"][0]["groupid"] == group["group_id"]:
                        serverList.append({"serverid": i["serverid"], "channelid": i["channelid"], "title": i["manga"][0]["title"]})
            elif "groupid" not in i["manga"][0]:
                serverList.append({"serverid": i["serverid"], "channelid": i["channelid"], "title": i["manga"][0]["title"]})
        if serverList != []:
            return serverList
        else:
            return None

    async def manga_wanted_user(self, group_list, manga_id=None, manga_title=None):
        userList = []
        if manga_title is not None:
            result = self.usr.find({"manga.title": manga_title}, {"userid": 1, "manga.$": 1})
        else:
            result = self.usr.find({"manga.id": manga_id}, {"userid": 1, "manga.$": 1})
        for i in result:
            if "groupid" in i["manga"][0]:
                for group in group_list:
                    if i["manga"][0]["groupid"] == group["group_id"]:
                        userList.append({"userid": i["userid"], "title": i["manga"][0]["title"]})
            elif "groupid" not in i["manga"][0]:
                userList.append({"userid": i["userid"], "title": i["manga"][0]["title"]})
        if userList != []:
            return userList
        else:
            return None

    async def set_scan_group_server(self, serverid, manga_id, group_id, group_name):
        self.srv.update_one({"serverid": serverid, "manga.id": manga_id}, {"$set": {"manga.$.groupName": group_name, "manga.$.groupid": group_id}})

    async def set_scan_group_user(self, userid, manga_id, group_id, group_name):
        self.usr.update_one({"userid": userid, "manga.id": manga_id}, {"$set": {"manga.$.groupName": group_name, "manga.$.groupid": group_id}})

    # hella scuffed, dont use lmao
    def update_all_ids(self, mode):
        if mode == "server":
            result = self.srv.find({}, {"_id": 1})
            for i in result:
                a = self.srv.find({"_id": i["_id"]}, {"manga": 1})
                for j in a:
                    print(j["_id"])
                    for k in j["manga"]:
                        print(k)
                        req = requests.get(f"https://www.mangaupdates.com/series.html?id={k['id']}")
                        soup = bs(req.text, "html.parser")
                        new = soup.find("link", {"rel": "canonical"})["href"]
                        link = new.partition("https://www.mangaupdates.com/series/")[2]
                        mangaid = link.partition("/")[0]
                        mangaid = int(mangaid, 36)
                        self.srv.update_one({"_id": i["_id"], "manga": {"$elemMatch": {"title": k["title"]}}}, {"$set": {"manga.$.id": mangaid}})
                        time.sleep(8)
                        if "groupid" in k:
                            req = requests.get(f"https://www.mangaupdates.com/groups.html?id={k['groupid']}")
                            soup = bs(req.text, "html.parser")
                            new = soup.find("link", {"rel": "canonical"})["href"]
                            link = new.partition("https://www.mangaupdates.com/group/")[2]
                            mangaid = link.partition("/")[0]
                            mangaid = int(mangaid, 36)
                            self.srv.update_one({"_id": i["_id"], "manga": {"$elemMatch": {"title": k["title"]}}}, {"$set": {"manga.$.groupid": mangaid}})
                            time.sleep(8)
        if mode == "user":
            result = self.usr.find({}, {"_id": 1}).skip(200)
            for i in result:
                a = self.usr.find({"_id": i["_id"]}, {"manga": 1})
                for j in a:
                    print(j["_id"])
                    for k in j["manga"]:
                        print(k)
                        req = requests.get(f"https://www.mangaupdates.com/series.html?id={k['id']}")
                        soup = bs(req.text, "html.parser")
                        new = soup.find("link", {"rel": "canonical"})["href"]
                        link = new.partition("https://www.mangaupdates.com/series/")[2]
                        mangaid = link.partition("/")[0]
                        mangaid = int(mangaid, 36)
                        self.usr.update_one({"_id": i["_id"], "manga": {"$elemMatch": {"title": k["title"]}}}, {"$set": {"manga.$.id": mangaid}})
                        time.sleep(8)
                        if "groupid" in k:
                            req = requests.get(f"https://www.mangaupdates.com/groups.html?id={k['groupid']}")
                            soup = bs(req.text, "html.parser")
                            new = soup.find("link", {"rel": "canonical"})["href"]
                            link = new.partition("https://www.mangaupdates.com/group/")[2]
                            mangaid = link.partition("/")[0]
                            mangaid = int(mangaid, 36)
                            self.usr.update_one({"_id": i["_id"], "manga": {"$elemMatch": {"title": k["title"]}}}, {"$set": {"manga.$.groupid": mangaid}})
                            time.sleep(8)
    
    def test(self):
        self.usr.update_many({"manga.title": "Berserk"}, {"$set": {"manga.$.id": 51239621230}})
        a = self.usr.find({"manga.id": ""}, {"manga": 1})
        for i in a:
            print(i)
