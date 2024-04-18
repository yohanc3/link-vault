import requests
import json
import re
class Vectara:

    def __init__(self, api_key, customer_id, corpus_id) -> None:
        self.api_key =api_key
        self.customer_id = customer_id
        self.corpus_id = corpus_id
    
    def file_upload(self, file, link = ""):
        session = requests.Session()
        files={
            "file": (file, open(file, "rb"), "application/octet-stream"),
            "doc_metadata": json.dumps({"link": link})
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
            "numResults": 5,
            "contextConfig": {
                "sentences_before": 3,
                "sentences_after": 3,
                "start_tag": "<b>",
                "end_tag": "</b>"
            },
            "corpusKey": [
                {
                "corpus_id": self.corpus_id
                }
            ],
            "summary": [
                {
                "max_summarized_results": 1,
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
    
    # API key, customer ID, and corpus ID
    vectara = Vectara("zut_yym2wk3zzNjsw3-xyZV04OV-ibmTlynUbrJQkw", 3408508610, 3)
    
    # print(vectara.file_upload("cities.txt", "https://facebook.com"))
    
    convo = vectara.ask_question(input("Enter your message: "))
    while convo is not None:
        print(convo["answer"])
        print("Documents and Links:", convo["documents"])
        convo = vectara.ask_question(input("Enter your question: "), convo["conversation_id"])
    
    
