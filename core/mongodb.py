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
    
    def add_server(self, server_name, server_id, channel_id):
        document_count = self.srv.count_documents({})
        while document_count >= 0:
            try:
                self.srv.insert_one({"_id": document_count, "serverid": server_id, "serverName": server_name, "channelid": channel_id, "manga": []})
                break
            except:
                document_count -= 1

    def add_user(self, user_name, user_id):
        document_count = self.usr.count_documents({})
        while document_count >= 0:
            try:
                self.usr.insert_one({"_id": document_count, "userid": user_id, "username": user_name, "manga": []})
                break
            except:
                document_count -= 1

    def remove_server(self, server_id):
        self.srv.delete_one({"serverid": server_id})

    def remove_user(self, user_id):
        self.usr.delete_one({"userid": user_id})