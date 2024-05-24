package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	. "wisdom-hoard/config"

	. "wisdom-hoard/error"
	. "wisdom-hoard/logger"
	"wisdom-hoard/storage"
	"wisdom-hoard/util"

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
	
	//add event handler
	discord.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate){
		b.NewMessage(session, message, b.storage)
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

func (b *Bot) NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	//Prevent bot responding to its own message
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, BOT_PREFIX+"find"):
		b.handleFindCommand(discord, message, storage)

	case strings.Contains(message.Content, BOT_PREFIX+"save"):
		b.handleSaveCommand(discord, message, storage)

	case strings.Contains(message.Content, BOT_PREFIX+"tags"):
		b.handleGetTagsCommand(discord, message, storage)

	case strings.Contains(message.Content, BOT_PREFIX+"help"):
		handleHelpCommand(discord, message)

	case strings.HasPrefix(message.Content, BOT_PREFIX):
		discord.ChannelMessageSend(message.ChannelID, "Oops, that's not a valid command! Try "+BOT_PREFIX+"help to find all valid commands!")
	}

}

func (b *Bot) handleFindCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

		username := message.Author.Username

		tags, err := util.ParseFindCommand(message.Content)

		if err != nil {
			fmt.Println("error is: ", err)
			b.sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		linksArr, err := storage.GetLinks(username, tags)
		if err != nil {
			b.sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		if len(linksArr) == 0 {
			b.sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		var stringifiedLinks string = strings.Join(linksArr, " ")

		discord.ChannelMessageSend(message.ChannelID, stringifiedLinks)

}

func (b *Bot ) handleSaveCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	username := message.Author.Username

	url, tags, err := util.ParseSaveCommand(message.Content)

	if err != nil {
		fmt.Println("error is ")
		b.sendErrorMessage(discord, message.ChannelID, err)
		return
	}

	//Merged tags are received only if the link has been saved before
	mergedTags, err := storage.InsertLinkAndTags(username, url, tags)

	if err != nil {
		b.sendErrorMessage(discord, message.ChannelID, err)
		return
	}

	if mergedTags != nil {
		discord.ChannelMessageSend(message.ChannelID, "Your link has been successfully updated! The new tags for this link are: " + util.FormatArrayToString(mergedTags))
		return
	}

	discord.ChannelMessageSend(message.ChannelID, "Yay! Your save was a success!")

	fmt.Println("Success!!")	
}

func (b *Bot ) handleGetTagsCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	if message.Content != BOT_PREFIX+"tags"{
		discord.ChannelMessageSend(message.ChannelID, "Wrong command call! \nExample: \n > "+BOT_PREFIX+"tags")
	}

	tags, err := storage.GetUserTags(message.Author.Username)

	if err != nil {
		b.sendErrorMessage(discord, message.ChannelID, err)
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
				Name:   BOT_PREFIX+"save",
				Value:  "Saves a url. You have to pass a valid url and tags that describe the type of content. \n Example: \n > "+BOT_PREFIX+"save https://example.com/ cinema sports literature\n\n",
				Inline: false,
			},
			{
				Name:   BOT_PREFIX+"find",
				Value:  "Given a list of categories, it retrieves all links that contain at least one of the given categories. \n Example: \n > "+BOT_PREFIX+"find cinema sports literature\n\n",
				Inline: false,
			},
			{
				Name:   BOT_PREFIX+"tags",
				Value:  "Returns all active tags you have previously used. \n Example: \n > "+BOT_PREFIX+"tags\n\n",
				Inline: false,
			},
		},
	}
	discord.ChannelMessageSendEmbed(message.ChannelID, embed)
}

func (b *Bot ) sendErrorMessage(discord *discordgo.Session, channelId string, err error){
	var message string
	if BotError, ok := err.(*Error); ok {
		message = BotError.UserMessage
	} else {
			message = GenericErrorMessage
	}
	discord.ChannelMessageSend(channelId, message)
}
