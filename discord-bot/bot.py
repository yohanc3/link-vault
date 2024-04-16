import nextcord
from nextcord.ext import commands
from dotenv import load_dotenv
import os

load_dotenv()
token = os.getenv("TOKEN")

bot = commands.Bot()

intents = nextcord.Intents.default()
intents.messages = True # Enable the messages intent

bot = commands.Bot(command_prefix='!', intents=intents)

@bot.event
async def on_message(message):
    # Check if the message is a DM
    if isinstance(message.channel, nextcord.DMChannel):
        # Reply to the DM
        await message.reply('This is a reply to your DM!')

bot.run(token)