package main

import (
	"fmt"
	"os"

	"github.com/yohanc3/link-vault/bot"
	. "github.com/yohanc3/link-vault/logger"
	"github.com/yohanc3/link-vault/storage"

	"github.com/joho/godotenv"
)

func main(){

	godotenv.Load()
	

	var DISCORD_BOT_TOKEN string

	MODE := os.Getenv("MODE")

	if MODE == "PROD"{
		DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
	} else if MODE == "DEV" {
		DISCORD_BOT_TOKEN = os.Getenv("DISCORD_TEST_BOT_TOKEN")
	} else {
		GeneralLogger.Panic().Str("MODE", MODE).Msg(".env var MODE is not a valid value. Only 'DEV' or 'PROD' are valid")
	}

	
	POSTGRES_DB_URI := os.Getenv("POSTGRES_DB_URI")

	fmt.Println("Discord bot token: ", DISCORD_BOT_TOKEN)
	
	storage := storage.NewPostgresDb(POSTGRES_DB_URI)

	bot := bot.NewBot(DISCORD_BOT_TOKEN, storage)

	bot.Run()

}