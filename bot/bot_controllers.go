package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	. "github.com/yohanc3/link-vault/config"
	. "github.com/yohanc3/link-vault/logger"
	"github.com/yohanc3/link-vault/storage"
)

func (b *Bot) NewMessage(discord *discordgo.Session, message *discordgo.MessageCreate, storage storage.Storage){

	//Prevent bot responding to its own message
	if message.Author.ID == discord.State.User.ID {
		return
	}

	var hasPrefix bool = strings.HasPrefix(message.Content, "+")

	if hasPrefix{

		switch {
			case strings.Contains(message.Content, BOT_PREFIX+"find"):
				b.handleFindCommand(discord, message, storage)
	
			case strings.Contains(message.Content, BOT_PREFIX+"save"):
				b.handleSaveCommand(discord, message, storage)
	
			case strings.Contains(message.Content, BOT_PREFIX+"tags"):
				b.handleGetTagsCommand(discord, message, storage)
	
			case strings.Contains(message.Content, BOT_PREFIX+"delete"):
				b.handleDeleteCommand(discord, message, storage)
	
			case strings.Contains(message.Content, BOT_PREFIX+"help"):
				b.handleHelpCommand(discord, message.ChannelID)
	
			default:
				var messageContent string = "Oops, that's not a valid command! Try "+BOT_PREFIX+"help to find all valid commands!"
				sendMessage(discord, message.ChannelID, messageContent)
		}
	
	}
}

func (b *Bot) GuildMemberAdd(discord *discordgo.Session, m *discordgo.GuildMemberAdd){
	channels, err := discord.GuildChannels(m.GuildID)

	if err != nil {
		GeneralLogger.Error().Str("error", err.Error()).Str("user", m.Member.User.Username).Msg("errow when dming intro message")
		return
	}

	var mainChannel *discordgo.Channel

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
				mainChannel = channel
				break
		}
	}

	if mainChannel == nil {
			fmt.Println("No text channel found in the guild")
			return
	}

	sendMessage(discord, mainChannel.ID, "Welcome to the server, "+m.Member.Mention()+"!"+"\nI can help you save and retrieve links. For more information do ` +help `")

}
