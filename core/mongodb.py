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

def add_server(serverName, serverID):
    srv.insert_one({"serverid": serverID, "serverName": serverName, "manga": []})

def add_user(userName, userID):
    usr.insert_one({"userid": userID, "userName": userName, "manga": []})

def remove_server(serverID):
    srv.delete_one({"serverid": serverID})

def remove_user(userID):
    usr.delete_one({"userid": userID})

def add_manga(idName, title, mode):
    if mode == "user":
        usr.update_one({"userid": idName}, {"$push": {"manga": {"title": title}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$push": {"manga": {"title": title}}})
    
def remove_manga(idName, title, mode):
    if mode == "user":
        usr.update_one({"userid": idName}, {"$pull": {"manga": {"title": title}}})
    elif mode == "server":
        srv.update_one({"serverid": idName}, {"$pull": {"manga": {"title": title}}})

