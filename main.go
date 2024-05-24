package main

import (
	"fmt"
	"os"
	"github.com/yohanc3/link-vault/bot"
	"github.com/yohanc3/link-vault/storage"

	"github.com/joho/godotenv"
)

func main(){

	godotenv.Load()

	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	POSTGRES_DB_URI := os.Getenv("POSTGRES_DB_URI")

	fmt.Println("Discord bot token: ", DISCORD_BOT_TOKEN)
	
	storage := storage.NewPostgresDb(POSTGRES_DB_URI)

	bot := bot.NewBot(DISCORD_BOT_TOKEN, storage)

	bot.Run()

}