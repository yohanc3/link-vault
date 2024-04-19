import requests
import json
import re
import tempfile
import uuid
import os
from dotenv import load_dotenv

class Vectara:

    def __init__(self, api_key, customer_id, corpus_id) -> None:
        self.api_key =api_key
        self.customer_id = customer_id
        self.corpus_id = corpus_id
    
    def file_upload(self, file_text = "", link = "", user_id = "", topic = ""):
        session = requests.Session()
        with tempfile.NamedTemporaryFile(delete=False, mode='w+') as tmp_file:
            tmp_file.write(file_text)
            tmp_file_name = tmp_file.name
            file = tmp_file_name
        files={
            "file": (str(uuid.uuid4()), open(file, "rb"), "application/octet-stream"),
            "doc_metadata": json.dumps({
                "link": link,
                "user_id": user_id,
                "topic": topic
            })
        }
        post_headers = {
        "x-api-key": self.api_key
        }
        response = session.post(
            url = f"https://api.vectara.io/v1/upload?c={self.customer_id}&o={self.corpus_id}",
            headers=post_headers,
            files=files)

        return_value = response.text
        session.close()
        return return_value
    
    def ask_question(self, question, conversation_id = None):
        url = "https://api.vectara.io/v1/query"

        payload = json.dumps({
        "query": [
            {
            "query": question,
            "start": 0,
            "numResults": 10,
            "contextConfig": {
                "chars"
                "sentences_before": 3,
                "sentences_after": 3,
                "start_tag": "<b>",
                "end_tag": "</b>"
            },
            "corpusKey": [
                {
                "corpus_id": self.corpus_id,
                "semantics": 0,
                "lexical_interpolation_config": {
                    "lambda": 0.025,
                },
                "dim": []
                }
            ],
            "summary": [
                {
                "max_summarized_results": 5,
                "response_lang": "en",
                "chat":  {
                "store": True,
                } if conversation_id is None else {
                     "conversationId": conversation_id,
                        "store": True
                }
                }
            ]
            }
        ]
        })
        headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'x-api-key': self.api_key
        }

        response = requests.request("POST", url, headers=headers, data=payload)
        
        status_code = response.status_code
        if status_code != 200:
            print(f"Error: {status_code}, {response.text}")
            return None

        data = json.loads(response.text)

        # Extract the summary text
        summary_text = data['responseSet'][0]['summary'][0]['text']

        # Extract the conversation ID
        conversation_id = data['responseSet'][0]['summary'][0]['chat']['conversationId']

        # Extract documents and their links
        documents = []
        for doc in data['responseSet'][0]['document']:
            doc_id = doc['id']
            link = next(item for item in doc['metadata'] if item['name'] == 'link')['value']
            documents.append({'id': doc_id, 'link': link})
        
        # Find all references in the summary text and print corresponding response texts
        references = re.findall(r'\[\d+\]', summary_text)
        response_texts = []
        for ref in references:
            index = int(ref.strip('[]')) - 1  # Convert '[1]' to 0, '[2]' to 1, etc.
            if index < len(data['responseSet'][0]['response']):
                response_text = data['responseSet'][0]['response'][index]['text']
                response_texts.append(f"{ref} {response_text}")

        # Printing the extracted data
        # print("Summary Text:", summary_text)
        # print("Conversation ID:", conversation_id)
        # print("Documents and Links:", documents)
        # print("References and Response Texts:", response_texts)
        return {"answer": summary_text, "conversation_id": conversation_id, "documents": documents, "references" :response_texts}


if __name__ == "__main__":
    load_dotenv()
    vectara_key = os.getenv("VECTARA_KEY")
    customer_id = os.getenv("CUSTOMER_ID")
    corpus_id = os.getenv("CORPUS_ID")

    # vectara api key, customer ID, and corpus ID
    vectara = Vectara(vectara_key, customer_id, int(corpus_id))
    
    # Either feed it file=filename or file_text="text" 
    print(vectara.file_upload(file_text='''Networking helps me in the following ways.
                            - people are more likely to reach out to me for help since they already know me and feel comfortable doing it
                            - I feel more comfortable reaching out to people for help because I’ve spoken to them already
                            - it helps me recognize patterns of problems teams/people are dealing with, and I can come up with common solutions that have impact across teams
                            - in general it opens up my mind to new perspectives
                            How I keep on top of it:
                            - I have a target for myself that I will reach out to two people every week to just say hello and introduce myself if we’re meeting for the first tins
                            How I find the people to network with
                            - I always make note of people names when they’re mentioned during conversations. The more conversations I have, the more names I add to my list!
                              ''', 
                              link="https://linkedin.com"))
    
    # convo = vectara.ask_question(input("Enter your message: "))
    # while convo is not None:
    #     print(convo["answer"])
    #     print("Documents and Links:", convo["documents"])
    #     convo = vectara.ask_question(input("Enter your question: "), convo["conversation_id"])