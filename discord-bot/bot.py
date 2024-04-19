# __import__('pysqlite3')
import sys
# sys.modules['sqlite3'] = sys.modules.pop('pysqlite3')

import os
import requests
from vectara import Vectara
# import chromadb
from datetime import datetime
from dotenv import load_dotenv

import discord
from discord.ext import commands

intents = discord.Intents.default()
intents.message_content = True
prev_message = {}

bot = commands.Bot(command_prefix='!', intents=intents)

@bot.event
async def on_ready():
    print(f'We have logged in as {bot.user}')

@bot.command(name="save")
async def process_link(ctx, link=None):
  """
  This command processes a link (optional) provided by the user.

  Args:
      ctx: The context of the command invocation.
      link: The link provided by the user as a string (optional).
  """

  user = ctx.author
  user_id = user.id

  if not link:
    await ctx.send("Please provide a link to process. You can type the command again with a link after it (e.g., `!save https://www.example.com`).")
    return

  if not link.startswith(("http://", "https://")):
    await ctx.send(f"Invalid link provided. Please enter a valid URL starting with 'http://' or 'https://'.")
    return

  full_link = f"https://r.jina.ai/{link}"

  try:
    response = requests.get(full_link)
    response.raise_for_status()  # Raise an exception for non-2xx status codes
    data = response.text
  except requests.exceptions.RequestException as e:
    await ctx.send(f"Error fetching content: {e}")
    return
  
  # save data to vectara with these properties: the chunk is the content. metadata is topic, user_id, and link to post

  print(f"content: {data}")
  load_dotenv()
  vectara_key = os.getenv("VECTARA_KEY")
  customer_id = os.getenv("CUSTOMER_ID")
  corpus_id = os.getenv("CORPUS_ID")
  vec = Vectara(vectara_key, customer_id, int(corpus_id))
  returns = vec.file_upload(file_text=data, link=link)
  
  if returns is None:
    await ctx.send("Error processing your content :(")
  else:
    await ctx.send("Content processed successfully!")

@bot.command(name="chat")
async def process_link(ctx, chat):
  """
  This command sends a chat message from the user to the Vectara API.

  Args:
      ctx: The context of the command invocation.
      chat: The chat message provided by the user.
  """

  user = ctx.author
  user_id = user.id

  if not chat:
    await ctx.send("Please provide a complete sentence for the chat message.")
    return

  vectara_key = os.getenv("VECTARA_KEY")
  customer_id = os.getenv("CUSTOMER_ID")
  corpus_id = os.getenv("CORPUS_ID")
  vec = Vectara(vectara_key, customer_id, int(corpus_id))
  returns = vec.ask_question(chat)
  
  if returns is None:
    await ctx.send("Error processing your chat message :(")
  else:
    global prev_message
    prev_message[user_id] = returns
    await ctx.send(returns["answer"])
    
@bot.command(name="sources")
async def process_link(ctx):
  """
  This command prints the sources (links) for the previous response.

  Args:
      ctx: The context of the command invocation.
      chat: The chat message provided by the user.
  """

  user = ctx.author
  user_id = user.id

  if prev_message[user_id] is None:
    await ctx.send("There is no previous message to show sources for.")
    return


  for doc in prev_message[user.id]["documents"]:
    await ctx.send(doc["link"])

load_dotenv()
discord_token = os.getenv("TOKEN")

bot.run(discord_token)
