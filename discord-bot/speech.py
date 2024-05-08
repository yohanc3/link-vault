from openai import OpenAI

client = OpenAI(api_key="")

# response = client.audio.speech.create(
#     model="tts-1",
#     voice="shimmer",
#     input='''TESTING TESTING TESTING [pause] [pause] [pause] [pause]
    
#     Hello and welcome to wisdom hoard! [pause] Your knowledge base redefined! [pause] [pause] [pause] [pause]

#     Here is the dream team. It has Abdool, ghiridar, and sabee. [pause] [pause] [pause] [pause]

#     Wisdom Hoard is an universal bookmarking tool. It is your personal knowledge vault, powered by a simple Discord bot. [pause] [pause]

#     To use it, simply send the bot a link to the post you want to save. You can also ask the bot for summary of a topic based on all of your saved posts so far. [pause] [pause] [pause] [pause] [pause] [pause]
    
#     Features include: [pause] bookmarking links from any social media site, [pause] intelligent search across all your posts, [pause] and LLM powered insights combining posts saved for a certain topic. [pause] [pause] [pause] [pause] [pause] [pause]
    
#     The tech stack The tech stack is made off of 3 tools: [pause] python, the discord bot API, and vectara. [pause] [pause]
    
#     Vectara provides rag as a service which is very convenient. [pause] [pause] [pause] [pause]

#     Now onto to the most exciting part, the demo. [pause]

#     This wisdom hoard bot has 3 commands: [pause] save, chat, and sources. [pause] [pause]

#     The save command saves a link. Let me show you by saving 2 tweets. I may have to try multiple times as the scraping service I depend on timeouts frequently. [pause] [pause] [pause] [pause]
#     '''
# )

# response.stream_to_file("body.mp3")

def text_to_speech(speech, name):
    response = client.audio.speech.create(
        model="tts-1",
        voice="shimmer",
        input=speech
    )

    response.stream_to_file(f"{name}.mp3")

phrases = {
    "testing": "Testing testing testing [pause] [pause] [pause] [pause]",
    # "chat": "The chat command lets you ask a question that is answered by a LLM. Related posts from your knowledge will be fed into it. Watch me use it now. I will begin by asking a question that I do not have a post saved about for. [pause] [pause] [pause] [pause]",
    # "save": "Now I will add a post giving advice to new grads. [pause] [pause] [pause] [pause]",
    # "ask": "Then let's ask the same question again. [pause] [pause] [pause] [pause]",
    # "neat": "Neat, right? [pause] [pause] [pause] [pause]",
    "sources": "Lastly, the sources command shows the links for each post you have saved in your knowledge.  Watch me use it now. [pause] [pause] [pause] [pause]",
    # "finish":"Awesome, all 3 commands are working! That's all for the demo and presentation! Thank you for watching!",
}

for name, speech in phrases.items():
    text_to_speech(speech, name)