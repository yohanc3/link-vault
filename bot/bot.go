package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	. "github.com/yohanc3/link-vault/config"

	. "github.com/yohanc3/link-vault/error"
	. "github.com/yohanc3/link-vault/logger"
	"github.com/yohanc3/link-vault/storage"
	"github.com/yohanc3/link-vault/util"

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
		var messageContent string = "Oops, that's not a valid command! Try "+BOT_PREFIX+"help to find all valid commands!"
		sendMessage(discord, message.ChannelID, messageContent)
	}

}

func (b *Bot) handleFindCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

		username := message.Author.Username

		tags, err := util.ParseFindCommand(message.Content)

		if err != nil {
			sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		linksArr, err := storage.GetLinks(username, tags)
		if err != nil {
			sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		if len(linksArr) == 0 {
			sendErrorMessage(discord, message.ChannelID, err)
			return
		}

		var embedFields []*discordgo.MessageEmbedField = []*discordgo.MessageEmbedField{}

		for i, v := range linksArr{
			var field discordgo.MessageEmbedField = discordgo.MessageEmbedField{
				Name: fmt.Sprint("→  ", i+1),
				Value: v,
				Inline: false,
				
			}
			embedFields = append(embedFields, &field)
		}

		stringifiedTags := strings.Join(tags, ", ")

		sendMessageEmbed(discord, message.ChannelID, "Results for tags " + "**```" + stringifiedTags + "```**", embedFields)

}

func (b *Bot ) handleSaveCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	username := message.Author.Username

	url, tags, err := util.ParseSaveCommand(message.Content)

	if err != nil {
		fmt.Println("error is ")
		sendErrorMessage(discord, message.ChannelID, err)
		return
	}

	//Merged tags are received only if the link has been saved before
	mergedTags, err := storage.InsertLinkAndTags(username, url, tags)

	if err != nil {
		sendErrorMessage(discord, message.ChannelID, err)
		return
	}

	if mergedTags != nil {
		var styledMergedTags string = strings.Join(mergedTags, " ")
		var messageContent string =  "Your link has been successfully updated! Now the tags for this link are " + "**```" + styledMergedTags + "```**"

		sendMessage(discord, message.ChannelID, messageContent)
		return
	}

	sendMessage(discord, message.ChannelID, "Yay! Your save was a success!")
}

func (b *Bot ) handleGetTagsCommand(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	if message.Content != BOT_PREFIX+"tags"{
		var messageContent string = "Wrong command call! \nExample: \n >" +BOT_PREFIX+ "tags" 
		sendMessage(discord, message.ChannelID, messageContent)
		return
	}

	tags, err := storage.GetUserTags(message.Author.Username)

	if err != nil {
		sendErrorMessage(discord, message.ChannelID, err)
	}

	formattedTags := strings.Join(tags, " ")
	var messageContent string = "Your previously used tags: \n " + "**```\n" + formattedTags + "\n```**"
	sendMessage(discord, message.ChannelID, messageContent)
}

func handleHelpCommand(discord *discordgo.Session, message *discordgo.MessageCreate){


	var embedFields []*discordgo.MessageEmbedField = []*discordgo.MessageEmbedField{
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
	}

	var embedTitle string = "All commands"
	sendMessageEmbed(discord, message.ChannelID, embedTitle, embedFields)
}

func sendMessage(discord *discordgo.Session, channelId string, message string){
	_, err := discord.ChannelMessageSend(channelId, message)

	if err != nil{
		GeneralLogger.Error().Str("error: ", err.Error()).Msg(DiscordMessageSendError.LogMessage)
	}

}

func sendMessageEmbed(discord *discordgo.Session, channelId string, title string, fields []*discordgo.MessageEmbedField){
	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0xfacd14,
		Fields: fields,
	}

	_, err := discord.ChannelMessageSendEmbed(channelId, embed)

	if err != nil {
		GeneralLogger.Error().Str("error", err.Error()).Msg(DiscordMessageSendError.LogMessage)
		return
	}
}

func sendErrorMessage(discord *discordgo.Session, channelId string, err error){
	var message string
	if BotError, ok := err.(*Error); ok {
		message = BotError.UserMessage
	} else {
			message = GenericErrorMessage
	}
	sendMessage(discord, channelId, message)
}
