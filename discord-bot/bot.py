__import__('pysqlite3')
import sys
sys.modules['sqlite3'] = sys.modules.pop('pysqlite3')

import os
import uuid
import requests
import chromadb
from datetime import datetime


import discord
from discord.ext import commands

intents = discord.Intents.default()
intents.message_content = True

chroma_client = chromadb.Client()
collection = chroma_client.create_collection(name="bookmarks")

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
  except requests.exceptions.RequestException as e:
    await ctx.send(f"Error fetching content: {e}")
    return

  data = response.text
  collection.add(
    documents=[data],
    metadatas=[{"source": link, "saved_by": user_id, "timestamp": f"datetime.now()"}],
    ids=[f"{uuid.uuid4()}"]
)
  

  results = collection.query(
        query_texts=["Title"],
        n_results=2
    )
  with open("query_results.txt", "a+") as f:
    # Write the results to the file
    for result in results:
        f.write(str(result) + "\n")  # Convert each result to a string and add a newline

    print("Results written to query_results.txt")

  await ctx.send(f"🎉 Saved into database!")



bot.run(os.environ['DISCORD_CLIENT_TOKEN'])

