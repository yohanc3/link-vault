__import__('pysqlite3')
import sys
sys.modules['sqlite3'] = sys.modules.pop('pysqlite3')

import os
import requests
import chromadb
from datetime import datetime
from dotenv import load_dotenv

import discord
from discord.ext import commands

load_dotenv()
token = os.getenv("TOKEN")

intents = discord.Intents.default()
intents.message_content = True

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

  # await ctx.send(f"🎉 Saved into database!")


bot.run(os.environ['DISCORD_CLIENT_TOKEN'])
