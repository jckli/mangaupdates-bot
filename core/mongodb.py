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
col = database["servers"]

def add_server(serverName, serverID):
    col.insert_one({"serverid": serverID, "serverName": serverName, "manga": []})

