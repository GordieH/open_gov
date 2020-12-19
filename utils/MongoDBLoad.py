
import pandas as pd
from pymongo import MongoClient
import numpy as np
import requests
import json
import csv
from bson import json_util

## TEMP

#load in configurations
with open('../data/config.json') as config_file:
    config = json.load(config_file)

# save house members list
URL = 'https://api.propublica.org/congress/v1/116/house/members.json'
headers = {'x-api-key': config['proPublica']}

response = requests.get(url = URL, headers = headers)
df = json_util.loads(response.text)
members_list = df['results'][0]['members']
len(members_list)
client = MongoClient("127.0.0.1")
db = client.database
db.representatives.insert_many(members_list)

db.representatives.find_one({'id':'D000197'})


#senate
# save house members list
URL = 'https://api.propublica.org/congress/v1/116/senate/members.json'
headers = {'x-api-key': config['proPublica']}

response = requests.get(url = URL, headers = headers)
df = json_util.loads(response.text)
members_list = df['results'][0]['members']
len(members_list)
client = MongoClient("127.0.0.1")
db = client.database
db.representatives.count_documents({})
db.representatives.insert_many(members_list)

db.representatives.find_one({'id':'W000817'})
for element in db.representatives.find({'state':'CO'}):
    print(element)
    print('\n')


## TEMP


def saveData(URL, headers, outputFile):
    outputFilePath = '../data/' + outputFile
    response = requests.get(url = URL, headers = headers)
    df = json.loads(response.text)
    members_list = df['results'][0]['members']
    outputFileObj = open(outputFilePath, 'w')
    csvWriter = csv.writer(outputFileObj)
    count=0
    for member in members_list:
        if count == 0:
            header = member.keys()
            csvWriter.writerow(header)
            count += 1
        csvWriter.writerow(member.values()) 
    outputFileObj.close()

    # df['results'][0]['members'][0].values()
    # with open(outputFilePath, 'w') as json_file:
    #     json.dump(df, json_file)

#load in configurations
with open('../data/config.json') as config_file:
    config = json.load(config_file)



# save house members list
URL = 'https://api.propublica.org/congress/v1/116/house/members.json'
headers = {'x-api-key': config['proPublica']}
saveData(URL=URL, headers=headers, outputFile='house_members.csv')

# save senate members list
URL = 'https://api.propublica.org/congress/v1/116/senate/members.json'
saveData(URL=URL, headers=headers, outputFile='senate_members.csv')