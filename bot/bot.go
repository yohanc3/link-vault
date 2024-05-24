package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"wisdom-hoard/config"
	"wisdom-hoard/storage"
	"wisdom-hoard/util"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token	string
	storage	storage.Storage
}

var BotToken string

const PREFIX string = config.BOT_PREFIX

func checkNill(e error){
	if e != nil{
		log.Fatal("Error message: ", e)
	}
}

func NewBot(botToken string, storage storage.Storage) *Bot{
	return &Bot{
		token: botToken,
		storage: storage,
	}
}

func (b *Bot) Run(){

	//create a session
	discord, err := discordgo.New("Bot " + b.token)
	checkNill(err)

	//add event handler
	discord.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate){
		NewMessage(session, message, b.storage)
	})

	//open session
	discord.Open()
	defer discord.Close()

	//Keep running the bot until forced termination (ctrl + c)
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1 )
	signal.Notify(c, os.Interrupt)
	<-c

}

func NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	//Prevent bot responding to its own message
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, PREFIX+"find"):
		handleFindCommand(discord, message, storage)

	case strings.Contains(message.Content, PREFIX+"save"):
		handleSaveCommand(discord, message, storage)

	case strings.Contains(message.Content, PREFIX+"tags"):
		handleGetTagsCommand(discord, message, storage)

	case strings.Contains(message.Content, PREFIX+"help"):
		handleHelpCommand(discord, message)

	case strings.HasPrefix(message.Content, PREFIX):
		discord.ChannelMessageSend(message.ChannelID, "Oops, that's not a valid command! Try "+PREFIX+"help to find all valid commands!")
	}

}

func handleFindCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

		username := message.Author.Username

		tags, err, notCriticalEror := util.ParseFindCommand(message.Content)

		if err != nil {
			fmt.Println("error is: ", err)
			if notCriticalEror {
				discord.ChannelMessageSend(message.ChannelID, err.Error())
				return
			}
			panic(err)
		}

		linksArr, err := storage.GetLinks(username, tags)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, err.Error())
			return
		}

		if len(linksArr) == 0 {
			discord.ChannelMessageSend(message.ChannelID, "You don't have any urls saved with this tag:(")
			return
		}

		var stringifiedLinks string = strings.Join(linksArr, " ")

		discord.ChannelMessageSend(message.ChannelID, stringifiedLinks)

}

func handleSaveCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	username := message.Author.Username

	url, tags, err, notCriticalError := util.ParseSaveCommand(message.Content)

	if err != nil {
		fmt.Println("error is ")
		if notCriticalError {
			discord.ChannelMessageSend(message.ChannelID, err.Error())
			return
		}
		panic(err)
	}

	//Merged tags are received only if the link has been saved before
	mergedTags, error := storage.InsertLinkAndTags(username, url, tags)

	if error != nil {
		discord.ChannelMessageSend(message.ChannelID, error.Error())
		return
	}

	if mergedTags != nil {
		discord.ChannelMessageSend(message.ChannelID, "Your link has been successfully updated! The new tags for this link are: " + util.FormatArrayToString(mergedTags))
		return
	}

	discord.ChannelMessageSend(message.ChannelID, "Yay! Your save was a success!")

	fmt.Println("Success!!")	
}

func handleGetTagsCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	if message.Content != PREFIX+"tags"{
		discord.ChannelMessageSend(message.ChannelID, "Wrong command call! \nExample: \n > "+PREFIX+"tags")
	}

	tags, err := storage.GetUserTags(message.Author.Username)

	if err != nil {
		discord.ChannelMessageSend(message.ChannelID, err.Error())
	}

	formattedTags := util.FormatArrayToString(tags)

	discord.ChannelMessageSend(message.ChannelID, "Your previously used tags: \n " + "**```\n" + formattedTags + "\n```**")
}

func handleHelpCommand(discord *discordgo.Session, message *discordgo.MessageCreate){
	embed := &discordgo.MessageEmbed{
		Title:       "All commands",
		Color:       0xfacd14,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   PREFIX+"save",
				Value:  "Saves a url. You have to pass a valid url and tags that describe the type of content. \n Example: \n > "+PREFIX+"save https://example.com/ cinema sports literature\n\n",
				Inline: false,
			},
			{
				Name:   PREFIX+"find",
				Value:  "Given a list of categories, it retrieves all links that contain at least one of the given categories. \n Example: \n > "+PREFIX+"find cinema sports literature\n\n",
				Inline: false,
			},
			{
				Name:   PREFIX+"tags",
				Value:  "Returns all active tags you have previously used. \n Example: \n > "+PREFIX+"tags\n\n",
				Inline: false,
			},
		},
	}
	discord.ChannelMessageSendEmbed(message.ChannelID, embed)
}