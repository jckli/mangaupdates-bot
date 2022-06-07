from pymongo import MongoClient
import certifi
import os

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

    async def get_manga_list_server(self, server_id):
        manga = []
        result = self.srv.find_one({"serverid": server_id}, {"manga": 1})
        for i in result["manga"]:
            manga.append(i["title"])
        if manga != []:
            return manga
        else:
            return None
    
    async def get_manga_list_user(self, user_id):
        manga = []
        result = self.usr.find_one({"userid": user_id}, {"manga": 1})
        for i in result["manga"]:
            manga.append(i["title"])
        if manga != []:
            return manga
        else:
            return None