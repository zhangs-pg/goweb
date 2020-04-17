import requests
import json

def testregist():
    data = {"name":"wangwu", "password":"123"}
    url = "http://127.0.0.1:8080/user/regist"
    r = requests.post(url, json.dumps(data))
    print(r.text)

def login():
    headers={}
    # headers['Content-Type']='application/json'
    headers["token"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NywiTmFtZSI6Indhbmd3dSJ9.0KGy5bmDZBpYlEqlflXZ57swa6PcbCoZF6MDJ9-QjR4"
    data = {"name":"wangwu", "password":"123"}
    url = "http://127.0.0.1:8080/user/login"
    r = requests.post(url, json.dumps(data), headers=headers)
    print(r.text)
login()

def he():
    headers={}
    # headers['Content-Type']='application/json'
    headers["token"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NywiTmFtZSI6Indhbmd3dSJ9.0KGy5bmDZBpYlEqlflXZ57swa6PcbCoZF6MDJ9-QjR4"

    url = "http://127.0.0.1:8080/test/hello"
    r = requests.get(url,headers=headers)
    print(r.text)
# he()