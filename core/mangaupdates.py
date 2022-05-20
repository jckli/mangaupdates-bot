import os
import requests

class MangaUpdates:
    def __init__(self):
        username = os.environ.get("MU_USER")
        password = os.environ.get("MU_PASS")
        loginurl = "https://api.mangaupdates.com/v1/account/login"
        login = requests.put(loginurl, json={"username": username, "password": password}).json()
        self.token = login["context"]["session_token"]
    
    def get_token(self):
        return self.token
