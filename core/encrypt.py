from Crypto.Cipher import AES
from base64 import b64encode, b64decode
import json
import os

with open("config.json", "r") as f:
    config = json.load(f)

key = config["encryptionKey"]
salt = config["encryptionSalt"]

class Crypt:
    def __init__(self):
        self.salt = salt.encode('utf8')
        self.method = 'utf-8'

    def encrypt(self, data):
        aes_obj = AES.new(key.encode('utf-8'), AES.MODE_CFB, self.salt)
        enc = aes_obj.encrypt(data.encode('utf8'))
        encstr = b64encode(enc).decode(self.method)
        return encstr

    def decrypt(self, data):
        aes_obj = AES.new(key.encode('utf8'), AES.MODE_CFB, self.salt)
        str_tmp = b64decode(data.encode(self.method))
        str_dec = aes_obj.decrypt(str_tmp)
        decstr = str_dec.decode(self.method)
        return decstr
