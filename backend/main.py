import os
from typing import List

import requests
from dotenv import load_dotenv
from fastapi import FastAPI, File, UploadFile

load_dotenv()
vectara_key = os.getenv("VECTARA_KEY")
# sabih's
customer_id = 2180615512

app = FastAPI()

@app.post("/file")
async def upload_file(file: UploadFile = File(...)):
    contents = await file.read()
    data = contents.decode("utf-8")

    url = "https://api.vectara.io/v1/upload"

    params = {
        'c': '2180615512', 
        'o': 5,
        'd': True
    }

    files = {
        "file": data
    }
    headers = {
        "Content-Type": "multipart/form-data",
        "Accept": "application/json",
        f"x-api-key": "{vectara_key}",
    }

    response = requests.request("POST", url, headers=headers, files=files, params=params)

    return {"data": data, "response": response.text}


@app.post("/related_posts")
async def related_posts():
    pass


@app.post("/question")
async def question():
    pass
