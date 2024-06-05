package bot

import (
	"fmt"
	"os"
	"os/signal"

	. "github.com/yohanc3/link-vault/logger"
	"github.com/yohanc3/link-vault/storage"

	"github.com/bwmarrin/discordgo"
)


type Bot struct {
	token	string
	storage	storage.Storage
}

var BotToken string

func NewBot(botToken string, storage storage.Storage) *Bot{
	return &Bot{
		token: botToken,
		storage: storage,
	}
}

func (b *Bot) Run(){

	//create a session
	discord, err := discordgo.New("Bot " + b.token)
	if err != nil {
		GeneralLogger.Fatal().Msg("Discordgo bot could not run. Error: " + err.Error())
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsDirectMessages | discordgo.IntentGuildMessages

	//add event handler
	discord.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate){
		b.NewMessage(session, message, b.storage)
	})
	discord.AddHandler(b.GuildMemberAdd)

	err = discord.Open()
	if err != nil {
			GeneralLogger.Fatal().Str("error", err.Error()).Msg("Could not open discord session.")
			return
	}

	defer discord.Close()

	err = discord.UpdateGameStatus(0, "+help | linkvault.me")

	if err != nil {
		GeneralLogger.Panic().Str("error", err.Error()).Msg("Error when setting bot's complex status")
	}

	//Keep running the bot until forced termination (ctrl + c)
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1 )
	signal.Notify(c, os.Interrupt)
	<-c

}